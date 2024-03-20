package main

import (
	"container/list"
	"strings"
	"sync"
	"time"
)

type Article struct {
	Id          string // primary key
	Author_id   string // foreign key into author
	Category_id string // foreign key into category
	Title       string
	Image       string
	Content     string
	Tags        string // CSV tags
	Created_at  time.Time
	Updated_at  time.Time
	Likes       int
	Size        int
}

type Author struct {
	Id           string // primary key
	Name         string
	Num_articles int
}

type Category struct {
	Id          string // primary key
	Article_ids string // foreign key into article
	Name        string
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
func (lru *LRU) GetArticlesByTag(tag string, isLimit bool, limit int) []*Article {
	var arr []*Article = make([]*Article, 0)
	for e := lru.Article_list.Front(); e != nil; e = e.Next() {
		article := e.Value.(*Article)
		if strings.Contains(article.Tags, tag) {
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
func (lru *LRU) GetArticlesByCategory(categoryId string, isLimit bool, limit int) []*Article {
	var arr []*Article = make([]*Article, 0)
	for e := lru.Article_list.Front(); e != nil; e = e.Next() {
		article := e.Value.(*Article)
		if article.Category_id == categoryId {
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
