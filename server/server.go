package main

import (
	"encoding/json"
	"fmt"
	"os"

	// "log"
	"net/http"
	"strconv"
	"time"

	c "github.com/EmilioChan27/Dist-Cache/common"
	d "github.com/EmilioChan27/Dist-Cache/db"
)

// var numRequests int = 0
var db *d.DB
var file *os.File

func main() {
	db = d.NewDB()
	file, err := os.Create("dbAccessTime.txt")
	c.CheckErr(err)
	file.WriteString(fmt.Sprintf("%v\n", time.Now()))
	// coldCapacity = 350
	// hotCapacity = 350
	// timerDuration = 60 * time.Second
	// writeChanLen = 75
	// cache = c.NewCache(coldCapacity, hotCapacity, 1, timerDuration, writeChanLen)
	// // // articles := make([]*c.Article, 10)
	// articles, err := db.GetNewestArticles(750)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, article := range articles {
	// 	cache.Add(article)
	// }
	// cache.ToString()

	http.HandleFunc("/front-page", getFrontPageArticles)
	http.HandleFunc("/business", getBusinessArticles)
	http.HandleFunc("/human-interest", getHumanInterestArticles)
	http.HandleFunc("/international-affairs", getInternationalAffairsArticles)
	http.HandleFunc("/sports", getSportsArticles)
	http.HandleFunc("/politics", getPoliticsArticles)
	http.HandleFunc("/science-technology", getScienceTechnologyArticles)
	http.HandleFunc("/breaking-news", getBreakingNewsArticles)
	http.HandleFunc("/arts-culture", getArtsCultureArticles)
	http.HandleFunc("/article", getArticleById)
	http.HandleFunc("/write", writeHandler)
	// startTimer(cache)
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func recordDBExecTime(beforeTime time.Time) {
	file.WriteString(fmt.Sprintf("%v\n", time.Since(beforeTime).Microseconds()))
}

// func startTimer(cache *c.Cache) {
// 	go func(cache *c.Cache) {
// 		for {
// 			<-cache.DbTimer.C
// 			write := cache.GetWrite()
// 			if write != nil {
// 				if write.Operation == "create" {
// 					db.AddArticle(write.Article)
// 				} else {
// 					log.Fatal("Something went wrong - the operation isn't create lol")
// 				}
// 			} else {
// 				fmt.Println("received no write")
// 			}
// 			cache.DbTimer.Reset(timerDuration)
// 		}
// 	}(cache)
// }

// func getFrontPage(w http.ResponseWriter, r *http.Request) {

// }

//	func getBusinessArticles(w http.ResponseWriter, r *http.Request) {
//		articles, err := db.GetArticlesByCategory("Business")
//		c.CheckErr(err)
//		for _, a := range articles {
//			writeArticleToFile(a, "output.txt")
//		}
//	}
func writeHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var a *c.Article
	err := dec.Decode(&a)
	// fmt.Println(*a)
	c.CheckErr(err)
	// if r.Method == "PUT" {
	// 	oldWrite := cache.AddWrite(&c.Write{Operation: "Edit", Article: a})
	// 	if oldWrite != nil && oldWrite.Operation == "create" {
	// 		db.AddArticle(oldWrite.Article)
	// 	}
	// } else if r.Method == "POST" {
	// 	fmt.Println("am about to create a write")
	// 	cache.Add(a)
	// 	oldWrite := cache.AddWrite(&c.Write{Operation: "create", Article: a})
	// 	if oldWrite != nil && oldWrite.Operation == "create" {
	// 		db.AddArticle(oldWrite.Article)
	// 	}
	// } else {
	// 	log.Fatalf("Actually, the method was %v\n", r.Method)
	// }
	beforeTime := time.Now()
	newId, err := db.AddArticle(a)
	recordDBExecTime(beforeTime)
	// a.Id = int(newId)
	c.CheckErr(err)
	enc := json.NewEncoder(w)

	enc.Encode(newId)
	// w.WriteHeader(201)
}

func getHumanInterestArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	// articles := cache.GetArticlesByCategory("Human Interest", limit, true)
	// if len(articles) < limit {
	// 	cache.ResetTimer(timerDuration)

	// }
	beforeTime := time.Now()
	articles, err := db.GetArticlesByCategory("Human Interest", limit)
	recordDBExecTime(beforeTime)
	c.CheckErr(err)
	encodeArticles(w, articles)
	// if limit < coldCapacity {
	// 	updateCache(articles)
	// } else {
	// 	// fmt.Println("Not updating the cache because the limit is too large")
	// }
}
func getInternationalAffairsArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	// articles := cache.GetArticlesByCategory("International Affairs", limit, true)
	// if len(articles) < limit {
	// 	cache.ResetTimer(timerDuration)

	// }
	beforeTime := time.Now()
	articles, err := db.GetArticlesByCategory("International Affairs", limit)
	recordDBExecTime(beforeTime)

	c.CheckErr(err)
	encodeArticles(w, articles)
	// if limit < coldCapacity {
	// 	updateCache(articles)
	// } else {
	// 	// fmt.Println("Not updating the cache because the limit is too large")
	// }
}
func getSportsArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	// articles := cache.GetArticlesByCategory("Sports", limit, true)
	// if len(articles) < limit {
	// 	cache.ResetTimer(timerDuration)

	// }
	beforeTime := time.Now()
	articles, err := db.GetArticlesByCategory("Sports", limit)
	recordDBExecTime(beforeTime)

	c.CheckErr(err)
	encodeArticles(w, articles)
	// if limit < coldCapacity {
	// 	updateCache(articles)
	// } else {
	// 	// fmt.Println("Not updating the cache because the limit is too large")
	// }
}
func getPoliticsArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	// articles := cache.GetArticlesByCategory("Politics", limit, true)
	// if len(articles) < limit {
	// 	cache.ResetTimer(timerDuration)

	// }
	beforeTime := time.Now()
	articles, err := db.GetArticlesByCategory("Politics", limit)
	recordDBExecTime(beforeTime)

	c.CheckErr(err)
	encodeArticles(w, articles)
	// if limit < coldCapacity {
	// 	updateCache(articles)
	// } else {
	// 	// fmt.Println("Not updating the cache because the limit is too large")
	// }
}
func getScienceTechnologyArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	// articles := cache.GetArticlesByCategory("Politics", limit, true)
	// if len(articles) < limit {
	// 	cache.ResetTimer(timerDuration)

	// }
	beforeTime := time.Now()
	articles, err := db.GetArticlesByCategory("Politics", limit)
	recordDBExecTime(beforeTime)

	c.CheckErr(err)
	encodeArticles(w, articles)
	// if limit < coldCapacity {
	// 	updateCache(articles)
	// } else {
	// 	// fmt.Println("Not updating the cache because the limit is too large")
	// }
}
func getBusinessArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	// articles := cache.GetArticlesByCategory("Politics", limit, true)
	// if len(articles) < limit {
	// 	cache.ResetTimer(timerDuration)

	// }
	beforeTime := time.Now()
	articles, err := db.GetArticlesByCategory("Politics", limit)
	recordDBExecTime(beforeTime)

	c.CheckErr(err)
	encodeArticles(w, articles)
	// if limit < coldCapacity {
	// 	updateCache(articles)
	// } else {
	// 	// fmt.Println("Not updating the cache because the limit is too large")
	// }
}
func getFrontPageArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	// articles := cache.GetArticlesByCategory("Politics", limit, true)
	// if len(articles) < limit {
	// 	cache.ResetTimer(timerDuration)

	// }
	beforeTime := time.Now()
	articles, err := db.GetArticlesByCategory("Politics", limit)
	recordDBExecTime(beforeTime)

	c.CheckErr(err)
	encodeArticles(w, articles)
	// if limit < coldCapacity {
	// 	updateCache(articles)
	// } else {
	// 	// fmt.Println("Not updating the cache because the limit is too large")
	// }
}
func getBreakingNewsArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	// articles := cache.GetArticlesByCategory("Breaking News", limit, true)
	// if len(articles) < limit {
	// 	cache.ResetTimer(timerDuration)

	// }
	beforeTime := time.Now()
	articles, err := db.GetArticlesByCategory("Breaking News", limit)
	recordDBExecTime(beforeTime)

	c.CheckErr(err)
	encodeArticles(w, articles)
	// if limit < coldCapacity {
	// 	updateCache(articles)
	// } else {
	// 	// fmt.Println("Not updating the cache because the limit is too large")
	// }
}
func getArtsCultureArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	// articles := cache.GetArticlesByCategory("Arts and Culture", limit, true)
	// if len(articles) < limit {
	// 	cache.ResetTimer(timerDuration)

	// }
	beforeTime := time.Now()
	articles, err := db.GetArticlesByCategory("Arts and Culture", limit)
	recordDBExecTime(beforeTime)

	c.CheckErr(err)
	encodeArticles(w, articles)
	// if limit < coldCapacity {
	// 	updateCache(articles)
	// } else {
	// 	// fmt.Println("Not updating the cache because the limit is too large")
	// }
}

func getArticleById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	c.CheckErr(err)
	// article := cache.GetArticleById(id)
	// if article == nil {
	// 	cache.ResetTimer(timerDuration)
	// }
	beforeTime := time.Now()
	article, err := db.GetArticleById(id)
	recordDBExecTime(beforeTime)
	c.CheckErr(err)
	articles := make([]*c.Article, 1)
	articles[0] = article
	encodeArticles(w, articles)
	// updateCache(articles)
}

func encodeArticles(w http.ResponseWriter, articles []*c.Article) {
	enc := json.NewEncoder(w)
	numArticlesEncoded := 0
	for _, a := range articles {
		if a != nil {
			enc.Encode(a)
			numArticlesEncoded++
		}
		// log.Fatal("BROHER AN ARTICEL IS NIL")
	}
	// fmt.Printf("Num articles Encoded: %d\n", numArticlesEncoded)
}

// func updateCache(articles []*c.Article) {
// 	for i := len(articles) - 1; i >= 0; i-- {
// 		cache.Add(articles[i])
// 	}
// }

// func writeArticleToFile(a *c.Article, filename string) {
// 	var file *os.File
// 	_, err := os.Stat(filename)
// 	if os.IsNotExist(err) {
// 		file, err = os.Create(filename)
// 	} else {
// 		file, err = os.Open(filename)
// 	}
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()
// 	file.WriteString("----------------------")
// 	file.WriteString(fmt.Sprintf("Title: %s\n", a.Title))
// 	file.WriteString(fmt.Sprintf("Category: %s\n", a.Category))
// 	file.WriteString(fmt.Sprintf("Author: %d\n", a.AuthorId))
// 	file.WriteString(fmt.Sprintf("Content Preview: %s\n", a.Content[:10]))
// 	file.WriteString(fmt.Sprintf("Created At: %v\n", a.CreatedAt))
// 	file.WriteString("-------------------------------")
// }
