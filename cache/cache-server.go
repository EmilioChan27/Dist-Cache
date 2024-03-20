package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
)

var lru *LRU

func main() {
	// keyVals := make(map[string]string)
	keyVals := new(sync.Map)
	lru = NewLRU(15)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Access-Control-Allow-Origin", "*")

		id := r.URL.Query().Get("id")
		val, found := keyVals.Load(id)
		if found {

			fmt.Printf("responding to request with number %s\n", id)
			fmt.Fprintf(w, "cached response to %s: %s", id, val)
		} else {
			keyVals.Store(id, id)
			fmt.Printf("Forwarding request with number %s\n", id)
			serverUrl, err := url.Parse("http://LX-Server:8080/")
			if err != nil {
				fmt.Println("Error was not nil")
			}
			server := httputil.NewSingleHostReverseProxy(serverUrl)
			server.ServeHTTP(w, r)
		}
	})
	http.HandleFunc("/get-all-articles", GetArticles)
	fmt.Println("Cache server running on port :8080")
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func GetArticles(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	var articles []*Article
	enc := json.NewEncoder(w)
	if limit == "" {
		articles = lru.GetArticles(false, -1)
	} else {
		intLimit, err := strconv.Atoi(limit)
		if err != nil {
			log.Fatal(err)
		}
		articles = lru.GetArticles(true, intLimit)
	}
	for _, article := range articles {
		enc.Encode(article)
	}
}

func GetArticleById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var article *Article
	enc := json.NewEncoder(w)
	if id == "" {
		log.Fatal("Id couldn't be parsed from url in getArticleById")
	} else {
		article = lru.GetArticleById(id)
		if article == nil {
			log.Fatalf("article with id %s couldn't be found\n", id)
		}
		enc.Encode(article)
	}
}
