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
	DbTimer   *time.Timer
	writes    chan (*Write)
	inBuffer  []*Article
	outBuffer []*Article
	stats     *Stats
	NewestId  int // the id of the newest item in the cache
}

type Write struct {
	Operation string   // create or edit or delete
	Article   *Article // the articles in question
}

type LRU struct {
	Capacity    int
	Size        int
	ArticleList *list.List
	IdToArticle *SyncMap
}

func NewCache(hotCapacity int, coldCapacity int, bufferSize int, timerDuration time.Duration, writeChanLen int, newestId int) *Cache {
	return &Cache{hot: NewLRU(hotCapacity), cold: NewLRU(coldCapacity), inBuffer: make([]*Article, bufferSize), outBuffer: make([]*Article, bufferSize), writes: make(chan *Write, writeChanLen), DbTimer: time.NewTimer(timerDuration), stats: &Stats{}, NewestId: newestId}
}

func NewLRU(capacity int) *LRU {
	return &LRU{Capacity: capacity, Size: 0, ArticleList: list.New(), IdToArticle: NewSyncMap()}
}

func (lru *LRU) Move(object *list.Element) {
	lru.ArticleList.MoveToFront(object)
}

func (c *Cache) ResetTimer(duration time.Duration) {
	c.DbTimer.Reset(duration)
}

func (c *Cache) GetWrite() *Write {
	if len(c.writes) == 0 {
		return nil
	}
	return <-c.writes
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
	lru.IdToArticle.Store(a.Id, elem)
	lru.Size++
	return outgoingElem
}

func (lru *LRU) Remove(a *Article) {
	if elem, found := lru.IdToArticle.Load(a.Id); found {
		lru.IdToArticle.Delete(a.Id)
		lru.ArticleList.Remove(elem.(*list.Element))
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
	if elem, found := hot.IdToArticle.Load(id); found {
		hot.Move(elem.(*list.Element))
		return elem.(*list.Element).Value.(*Article)
	}
	cold := c.cold
	if elem, found := cold.IdToArticle.Load(id); found {
		c.coldToHot(elem.(*list.Element).Value.(*Article))
		return elem.(*list.Element).Value.(*Article)
	}
	return nil
}

func (c *Cache) ModifyArticle(a *Article) bool {
	hot := c.hot
	if elem, found := hot.IdToArticle.Load(a.Id); found {
		elem.(*list.Element).Value = a
		hot.Move(elem.(*list.Element))
		return true
	}
	cold := c.cold
	if elem, found := cold.IdToArticle.Load(a.Id); found {
		elem.(*list.Element).Value = a
		c.coldToHot(elem.(*list.Element).Value.(*Article))
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

func (c *Cache) SetNewestId(id int) {
	c.Lock()
	c.NewestId = id
	c.Unlock()
}

// func (c *Cache) GetNewestArticles(limit int) []*Article {
// 	articleList := make([]*Article, limit)
// 	l := c.hot.ArticleList
// 	index := 0
// 	for
// }

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
	if elem, found := c.hot.IdToArticle.Load(a.Id); found {
		c.hot.Move(elem.(*list.Element))
		return
	}
	if elem, found := c.cold.IdToArticle.Load(a.Id); found {
		c.coldToHot(elem.(*list.Element).Value.(*Article))
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

func (c *Cache) AddWrite(w *Write) *Write {
	fmt.Println("about to add a write inside of caceh")
	var write *Write
	if len(c.writes) == cap(c.writes) {
		write = <-c.writes
	}
	c.writes <- w
	// fmt.Println("just added a write inside of cache")
	return write
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
