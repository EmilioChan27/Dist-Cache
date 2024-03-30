package common

import (
	"container/list"
	"log"
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

type Cache struct {
	sync.RWMutex
	lru       *LRU
	la        *LatencyAware
	inBuffer  []*Article
	outBuffer []*Article
}

type LRU struct {
	Capacity     int
	Curr_size    int
	Article_list *list.List
	IdToArticle  map[string]*list.Element
}
type LatencyAware struct {
	Capacity int
}

func NewCache(LRUCapacity int, LatencyAwareCapacity int, bufferSize int) *Cache {
	return &Cache{lru: NewLRU(LRUCapacity), la: NewLatencyAware(LatencyAwareCapacity), inBuffer: make([]*Article, bufferSize), outBuffer: make([]*Article, bufferSize)}
}
func NewLatencyAware(capacity int) *LatencyAware {
	return &LatencyAware{Capacity: capacity}
}

func NewLRU(capacity int) *LRU {
	return &LRU{Capacity: capacity, Curr_size: 0, Article_list: list.New(), IdToArticle: make(map[string]*list.Element)}
}

func (lru *LRU) SetObjects(objects []*list.Element) {
	for _, e := range objects {
		lru.Article_list.MoveToFront(e)
	}
}

func (lru *LRU) SetObject(object *list.Element) {
	lru.Article_list.MoveToFront(object)
}

// TODO add the la back in
func (c *Cache) GetArticleById(id string) *Article {
	lru := c.lru
	if val, found := lru.IdToArticle[id]; found {
		lru.Article_list.MoveToFront(val)
		return val.Value.(*Article)
	} else {
		return nil
	}

}

// no need to do any caching because it's accessing the whole cache lmao
// TODO add the la section
func (c *Cache) GetArticles(isLimit bool, limit int) []*Article {
	lru := c.lru
	var arr []*Article = make([]*Article, 0)
	for e := lru.Article_list.Front(); e != nil; e = e.Next() {
		article := e.Value.(*Article)
		arr = append(arr, article)
		if isLimit {
			go func() {
				lru.SetObject(e)
			}()
		}
		if isLimit && len(arr) == limit {
			break
		}
	}

	return arr
}

// TODO add the cache aspect back lol
func (lru *LRU) GetArticlesByTitle(title string, isLimit bool, limit int) []*Article {
	var arr []*Article = make([]*Article, 0)
	for e := lru.Article_list.Front(); e != nil; e = e.Next() {
		article := e.Value.(*Article)
		if strings.Contains(article.Title, title) {
			arr = append(arr, article)
		}
		go func() {
			lru.SetObject(e)
		}()
		if isLimit && len(arr) == limit {
			break
		}
	}
	return arr
}

// TODO add the cache aspect back lol
func (lru *LRU) GetArticlesByKeyword(keyword string, isLimit bool, limit int) []*Article {
	var arr []*Article = make([]*Article, 0)
	for e := lru.Article_list.Front(); e != nil; e = e.Next() {
		article := e.Value.(*Article)
		if strings.Contains(article.Content, keyword) {
			arr = append(arr, article)
		}
		go func() {
			lru.SetObject(e)
		}()
		if isLimit && len(arr) == limit {
			break
		}
	}
	return arr
}

// TODO add the cache aspect back lol
func (lru *LRU) GetArticlesByCategory(category string, isLimit bool, limit int) []*Article {
	var arr []*Article = make([]*Article, 0)
	for e := lru.Article_list.Front(); e != nil; e = e.Next() {
		article := e.Value.(*Article)
		if article.Category == category {
			arr = append(arr, article)
		}
		go func() {
			lru.SetObject(e)
		}()
		if isLimit && len(arr) == limit {
			break
		}
	}
	return arr
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
