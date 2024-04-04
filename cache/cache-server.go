package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	c "github.com/EmilioChan27/Dist-Cache/common"
)

var cache *c.Cache

// var serverUrl string = "http://LX-Server:8080/"
var serverUrl string = "http://localhost:8000/"

func main() {
	// keyVals := make(map[string]string)
	cache = c.NewCache(5, 5, 1)
	// articles := make([]*c.Article, 10)
	// for i := 0; i < 10; i++ {
	// 	a := &Article{Id: i, Category: "Human Interest", Content: "Today I went to the park", AuthorId: i}
	// 	cache.Add(a)
	// }

	// http.HandleFunc("/front-page", getFrontPage)
	// http.HandleFunc("/business", getBusinessArticles)
	http.HandleFunc("/human-interest", getHumanInterestArticles)
	// http.HandleFunc("/international-affairs", getInternationalAffairsArticles)
	// http.HandleFunc("/science-technology", getScienceTechnologyArticles)
	// http.HandleFunc("/get-all-articles", GetArticles)
	// http.HandleFunc("/article", getArticleById)
	fmt.Println("Cache server running on port :8080")
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func sendToServer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("sending to server")
	parsedServerUrl, err := url.Parse(serverUrl)
	if err != nil {
		log.Fatalf("error parsing serverUrl: %v\n", err)
	}
	//TODO reset the timer of the cache so that we know when the last database access was
	server := httputil.NewSingleHostReverseProxy(parsedServerUrl)
	server.ServeHTTP(w, r)

}

func getHumanInterestArticles(w http.ResponseWriter, r *http.Request) {
	var articles []*c.Article
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		log.Fatal(err)
	}
	articles = cache.GetArticlesByCategory("Human Interest", limit, false)
	if len(articles) < limit {
		sendToServer(w, r)
	}
}

// func GetArticles(w http.ResponseWriter, r *http.Request) {
// 	limit := r.URL.Query().Get("limit")
// 	var articles []*c.Article
// 	enc := json.NewEncoder(w)
// 	if limit == "" {
// 		articles = lru.GetArticles(false, -1)
// 	} else {
// 		intLimit, err := strconv.Atoi(limit)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		articles = lru.GetArticles(true, intLimit)
// 	}
// 	for _, article := range articles {
// 		enc.Encode(article)
// 	}
// }

// func GetArticleById(w http.ResponseWriter, r *http.Request) {
// 	id := r.URL.Query().Get("id")
// 	var article *c.Article
// 	enc := json.NewEncoder(w)
// 	if id == "" {
// 		log.Fatal("Id couldn't be parsed from url in getArticleById")
// 	} else {
// 		article = lru.GetArticleById(id)
// 		if article == nil {
// 			log.Fatalf("article with id %s couldn't be found\n", id)
// 		}
// 		enc.Encode(article)
// 	}
// }
