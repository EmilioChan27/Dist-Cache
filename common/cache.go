package common

import (
	"container/list"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"
)

type Article struct {
	Id        int    // primary key
	AuthorId  int    // foreign key into author
	Category  string // foreign key into category
	Title     string
	ImageUrl  string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Likes     int
	Size      int // size of obj in mb
}

type Author struct {
	Id            int // primary key
	Name          string
	Bio           string
	Email         string
	ImageUrl      string
	SpecialtyTags string // csv tags for their specialty categories
	NumArticles   int
}

type Category struct {
	Id   int // primary key
	Name string
}

type Stats struct {
	coldHits int
	misses   int
	hotHits  int
}
type Cache struct {
	sync.RWMutex
	hot       *LRU
	cold      *LRU
	dbTime    time.Time
	writes    list.List
	inBuffer  []*Article
	outBuffer []*Article
	stats     *Stats
}

type LRU struct {
	Capacity    int
	Size        int
	ArticleList *list.List
	IdToArticle map[int]*list.Element
}

func NewCache(hotCapacity int, coldCapacity int, bufferSize int) *Cache {
	return &Cache{hot: NewLRU(hotCapacity), cold: NewLRU(coldCapacity), inBuffer: make([]*Article, bufferSize), outBuffer: make([]*Article, bufferSize), writes: *list.New().Init(), dbTime: time.Now(), stats: &Stats{}}
}

func NewLRU(capacity int) *LRU {
	return &LRU{Capacity: capacity, Size: 0, ArticleList: list.New(), IdToArticle: make(map[int]*list.Element)}
}

func (lru *LRU) Move(object *list.Element) {
	lru.ArticleList.MoveToFront(object)
}

// adds a new article to the LRU and returns the list element that was removed to add it (or nil if none)
func (lru *LRU) Add(a *Article) *list.Element {
	var outgoingElem *list.Element
	if lru.Capacity <= lru.Size {
		// fmt.Println("in add, the capacity is too much")
		outgoingElem = lru.ArticleList.Back()
		// fmt.Printf("Am removing article %d\n", outgoingElem.Value.(*Article).Id)
		lru.Remove(outgoingElem.Value.(*Article))
	}
	elem := lru.ArticleList.PushFront(a)
	lru.IdToArticle[a.Id] = elem
	lru.Size++
	return outgoingElem
}

func (lru *LRU) Remove(a *Article) {
	if elem, found := lru.IdToArticle[a.Id]; found {
		delete(lru.IdToArticle, a.Id)
		lru.ArticleList.Remove(elem)
		lru.Size--
	}
}
func (c *Cache) coldToHot(a *Article) {
	c.cold.Remove(a)
	oldHotElem := c.hot.Add(a)
	if oldHotElem != nil {
		c.cold.Add(oldHotElem.Value.(*Article))
	}
}

func (c *Cache) GetArticleById(id int) *Article {
	hot := c.hot
	if elem, found := hot.IdToArticle[id]; found {
		hot.Move(elem)
		return elem.Value.(*Article)
	}
	cold := c.cold
	if elem, found := cold.IdToArticle[id]; found {
		c.coldToHot(elem.Value.(*Article))
		return elem.Value.(*Article)
	}
	return nil
}

func (c *Cache) ModifyArticle(a *Article) bool {
	hot := c.hot
	if elem, found := hot.IdToArticle[a.Id]; found {
		elem.Value = a
		hot.Move(elem)
		return true
	}
	cold := c.cold
	if elem, found := cold.IdToArticle[a.Id]; found {
		elem.Value = a
		c.coldToHot(elem.Value.(*Article))
		return true
	}
	return true
}

// returns all articles whose categories match a certain pattern BUT MIGHT NOT RETURN AS MANY AS THEY WANT
func (c *Cache) GetArticlesByCategory(category string, limit int, isLimit bool) []*Article {
	articleList := make([]*Article, limit)
	l := c.hot.ArticleList
	index := 0
	for e := l.Front(); e != nil; e = e.Next() {
		if isLimit && index == limit {
			// fmt.Println("Returning articles early in hot")
			return articleList
		}
		if strings.Contains(e.Value.(*Article).Category, category) {
			articleList[index] = e.Value.(*Article)
			index++
		}
	}
	l = c.cold.ArticleList
	for e := l.Front(); e != nil; e = e.Next() {
		if isLimit && index == limit {
			// fmt.Println("Returning articles early in cold")
			return articleList
		}
		if strings.Contains(e.Value.(*Article).Category, category) {
			articleList[index] = e.Value.(*Article)
			index++
		}
	}
	// fmt.Printf("returning %d articles\n", index)
	if index == limit {
		return articleList
	}
	return articleList[:index]
}

// we know that since it's a write-through cache, the most recent stuff will definitely be in the cache
// func (c *Cache) GetBreakingNewsArticles(limit int, isLimit bool) []*Article {

// }

func (c *Cache) GetArticlesByFieldQuery(field string, query string, limit int, isLimit bool) []*Article {
	articleList := make([]*Article, limit)
	l := c.hot.ArticleList
	index := 0
	for e := l.Front(); e != nil; e = e.Next() {
		if isLimit && index == limit {
			return articleList
		}
		refArticle := reflect.ValueOf(e.Value.(*Article))
		valOfStringField := reflect.Indirect(refArticle).FieldByName(field)
		strVal := string(valOfStringField.String())
		if strings.Contains(strVal, query) {
			articleList[index] = e.Value.(*Article)
			index++
		}
	}
	// then go through the cold cache
	l = c.cold.ArticleList
	for e := l.Front(); e != nil; e = e.Next() {
		if isLimit && index == limit {
			return articleList
		}
		refArticle := reflect.ValueOf(e.Value.(*Article))
		valOfStringField := reflect.Indirect(refArticle).FieldByName(field)
		strVal := string(valOfStringField.String())
		if strings.Contains(strVal, query) {
			articleList[index] = e.Value.(*Article)
			index++
		}
	}
	if index == limit {
		return articleList
	}
	return articleList[:index]
}

func (c *Cache) Add(a *Article) {
	c.Lock()
	defer c.Unlock()
	if elem, found := c.hot.IdToArticle[a.Id]; found {
		c.hot.Move(elem)
		return
	}
	if elem, found := c.cold.IdToArticle[a.Id]; found {
		c.coldToHot(elem.Value.(*Article))
		return
	}
	if c.hot.Size < c.hot.Capacity {
		c.hot.Add(a)
	} else {
		oldColdElem := c.cold.Add(a)
		if oldColdElem != nil {
			c.outBuffer = append(c.outBuffer, oldColdElem.Value.(*Article))
		}
	}
}
func (c *Cache) ToString() {
	l := c.hot.ArticleList
	fmt.Printf("HOT size: %d\n", c.hot.Size)
	fmt.Println("---------------------")
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Printf("Id: %d -> ", e.Value.(*Article).Id)
	}
	fmt.Println("\n---------------------")
	fmt.Printf("COLD size: %d\n", c.cold.Size)
	fmt.Println("---------------------")
	l = c.cold.ArticleList
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Printf("Id: %d -> ", e.Value.(*Article).Id)
	}
	fmt.Println()
	fmt.Println("---------------------")
}
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
