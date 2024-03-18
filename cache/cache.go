package main

import (
	"container/list"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Article struct {
	Id          uuid.UUID // primary key
	Author_id   uuid.UUID // foreign key into author
	Category_id uuid.UUID // foreign key into category
	Title       string
	Image       string
	Content     string
	Created_at  time.Time
	Updated_at  time.Time
	Likes       int
	Size 		int
}

type Author struct {
	Id           uuid.UUID // primary key
	Name         string
	Num_articles int
}

type Category struct {
	Id          uuid.UUID // primary key
	Article_ids uuid.UUID // foreign key into article
	Name        string
}

type LRU struct {
	sync.Mutex
	Capacity     int
	Curr_size	 int
	Article_list *list.List
	IdToArticle  map[uuid.UUID]*Article

}

func (*LRU) addToLRU(article *Article) []*Article {
	for ()
}

type LatencyAware struct {

}
