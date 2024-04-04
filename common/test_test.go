package common

import (
	"fmt"
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	cache := NewCache(5, 5, 1)
	// articles := make([]*c.Article, 10)
	for i := 0; i < 10; i++ {
		a := &Article{Id: i}
		cache.Add(a)
	}
	cache.ToString()
}
func TestCacheUpdating(t *testing.T) {
	cache := NewCache(5, 5, 1)
	// articles := make([]*c.Article, 10)
	for i := 0; i < 10; i++ {
		a := &Article{Id: i, Category: "Human Interest", Content: "Today I went to the park", AuthorId: i}
		cache.Add(a)
	}
	// cache.ToString()
	// cache.GetArticleById(1)
	// cache.ToString()
	// cache.GetArticleById(0)
	cache.ToString()
	cache.GetArticleById(25)
	cache.ToString()
	articles := cache.GetArticlesByFieldQuery("Content", " went to the park", 25, true)
	for _, article := range articles {
		fmt.Printf("Id: %d\n", article.Id)
	}
	cache.ToString()
}

func TestInsertion(t *testing.T) {
	cache := NewCache(50, 50, 1)
	// articles := make([]*c.Article, 10)
	for i := 0; i < 10; i++ {
		a := &Article{Id: i, Category: "Human Interest", Content: "Today I went to the park", AuthorId: i}
		cache.Add(a)
	}
	cache.ToString()
	cache.Add(&Article{Id: 11, Category: "Human Interest", Content: "Today I went to the park", AuthorId: 11})
	cache.ToString()
}

func TestDupInsertion(t *testing.T) {
	cache := NewCache(10000, 10000, 1)
	// articles := make([]*c.Article, 10)
	for i := 0; i < 20000; i++ {
		a := &Article{Id: i, Category: "Human Interest", Content: "Today I went to the park", AuthorId: i}
		cache.Add(a)
	}
	cache.ToString()
	beforeTime := time.Now()
	for i := 0; i < 20000; i++ {
		a := &Article{Id: i + 20000, Category: "Human Interest", Content: "Today I went to the park", AuthorId: i + 20000}
		cache.Add(a)
	}
	execTime := time.Since(beforeTime)
	fmt.Printf("Execution time: %v\n", execTime)
	// cache.ToString()
}
