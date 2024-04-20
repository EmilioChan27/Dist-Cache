package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	c "github.com/EmilioChan27/Dist-Cache/common"
	d "github.com/EmilioChan27/Dist-Cache/db"
)

// var numRequests int = 0
var db *d.DB
var cache *c.Cache
var coldCapacity int
var hotCapacity int
var timerDuration time.Duration
var writeChanLen int

func main() {
	db = d.NewDB()
	hotCapacity = 650
	coldCapacity = 350
	timerDuration = 55 * time.Second
	writeChanLen = 300
	cacheFile, err := os.OpenFile("cacheAccessTime.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	c.CheckErr(err)
	cacheFile.WriteString(fmt.Sprintf("%v\n", time.Now()))
	cacheFile.Close()
	dbFile, err := os.OpenFile("dbAccessTimeCacheServer.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	c.CheckErr(err)
	dbFile.WriteString(fmt.Sprintf("%v\n", time.Now()))
	dbFile.Close()
	cacheUpdateFile, err := os.OpenFile("cacheUpdateTime.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	c.CheckErr(err)
	cacheUpdateFile.WriteString(fmt.Sprintf("%v\n", time.Now()))
	cacheUpdateFile.Close()

	// // articles := make([]*c.Article, 10)
	articles, newestId, err := db.GetNewestArticles(coldCapacity + hotCapacity)
	cache = c.NewCache(hotCapacity, coldCapacity, timerDuration, writeChanLen, newestId)
	if err != nil {
		log.Fatal(err)
	}
	for _, article := range articles {
		cache.Add(article)
	}
	cache.ToString()

	// http.HandleFunc("/front-page", getFrontPageArticles)
	http.HandleFunc("/business", getBusinessArticles)
	http.HandleFunc("/human-interest", getHumanInterestArticles)
	http.HandleFunc("/international-affairs", getInternationalAffairsArticles)
	http.HandleFunc("/sports", getSportsArticles)
	http.HandleFunc("/politics", getPoliticsArticles)
	http.HandleFunc("/science-technology", getScienceTechnologyArticles)
	http.HandleFunc("/breaking-news", getBreakingNewsArticles)
	http.HandleFunc("/arts-culture", getArtsCultureArticles)
	http.HandleFunc("/article", getArticleById)
	http.HandleFunc("/start-non-bursty-section", nonBurstySectionHandler)
	http.HandleFunc("/start-bursty-section", burstySectionHandler)
	http.HandleFunc("/write", writeHandler)
	startTimer(cache)
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func nonBurstySectionHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.OpenFile("dbAccessTimeCacheServer.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	c.CheckErr(err)
	file.WriteString(fmt.Sprintf("NON BURSTY SECTION %v\n", time.Now()))
	file.Close()
}

func burstySectionHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.OpenFile("dbAccessTimeCacheServer.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	c.CheckErr(err)
	file.WriteString(fmt.Sprintf("NON BURSTY SECTION %v\n", time.Now()))
	file.Close()
}

func recordCacheExecTime(beforeTime time.Time) {
	file, err := os.OpenFile("cacheAccessTime.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	// defer file.Close()
	c.CheckErr(err)
	file.WriteString(fmt.Sprintf("%v\n", time.Since(beforeTime).Microseconds()))
	file.Close()
}
func recordDBExecTime(beforeTime time.Time) {
	file, err := os.OpenFile("dbAccessTimeCacheServer.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	// defer file.Close()
	c.CheckErr(err)
	file.WriteString(fmt.Sprintf("%v\n", time.Since(beforeTime).Microseconds()))
	file.Close()
}

func recordCacheUpdateTime(beforeTime time.Time) {
	file, err := os.OpenFile("cacheUpdateTime.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	// defer file.Close()
	c.CheckErr(err)
	file.WriteString(fmt.Sprintf("%v\n", time.Since(beforeTime).Microseconds()))
	file.Close()
}

func startTimer(cache *c.Cache) {
	go func(cache *c.Cache) {
		for {
			<-cache.DbTimer.C
			write := cache.GetWrite()
			if write != nil {
				if write.Operation == "create" {
					id, err := db.AddArticle(write.Article)
					c.CheckErr(err)
					if int(id) > cache.NewestId {
						cache.SetNewestId(int(id))
					}
				} else {
					log.Fatal("Something went wrong - the operation isn't create lol")
				}
			} else {
				article, err := db.GetArticleById(cache.NewestId)
				c.CheckErr(err)
				cache.Add(article)
				fmt.Println("received no write")
			}
			cache.DbTimer.Reset(timerDuration)
		}
	}(cache)
}

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
	a.Id = cache.NewestId + 1
	go func(cache *c.Cache, db *d.DB) {
		cache.SetNewestId(cache.NewestId + 1)
		if r.Method == "PUT" {
			oldWrite := cache.AddWrite(&c.Write{Operation: "Edit", Article: a})
			if oldWrite != nil && oldWrite.Operation == "create" {
				db.AddArticle(oldWrite.Article)
			}
		} else if r.Method == "POST" {
			// fmt.Println("am about to create a write")
			cache.Add(a)
			beforeTime := time.Now()
			oldWrite := cache.AddWrite(&c.Write{Operation: "create", Article: a})
			recordCacheUpdateTime(beforeTime)
			if oldWrite != nil && oldWrite.Operation == "create" {
				id, err := db.AddArticle(oldWrite.Article)
				c.CheckErr(err)
				intId := int(id)
				if intId > a.Id {
					cache.SetNewestId(intId)
				}
			}
		} else {
			log.Fatalf("Actually, the method was %v\n", r.Method)
		}
	}(cache, db)
	enc := json.NewEncoder(w)
	enc.Encode(cache.NewestId)
	// w.WriteHeader(201)
}

func getHumanInterestArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	beforeTime := time.Now()
	articles := cache.GetArticlesByCategory("Human Interest", limit, true)
	if len(articles) < limit {
		beforeTime = time.Now()
		articles, err = db.GetArticlesByCategory("Human Interest", limit)
		recordDBExecTime(beforeTime)
		c.CheckErr(err)
		cache.ResetTimer(timerDuration)
	} else {
		recordCacheExecTime(beforeTime)
	}
	// go func(limit int, coldCapacity int, cache *c.Cache, articles []*c.Article) {
	// }(limit, coldCapacity, cache, articles)
	encodeArticles(w, articles)
	if limit < coldCapacity {
		beforeTime = time.Now()
		updateCache(cache, articles)
		recordCacheUpdateTime(beforeTime)
	}

}
func getInternationalAffairsArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	beforeTime := time.Now()
	articles := cache.GetArticlesByCategory("International Affairs", limit, true)
	if len(articles) < limit {
		beforeTime = time.Now()
		articles, err = db.GetArticlesByCategory("International Affairs", limit)
		c.CheckErr(err)
		recordDBExecTime(beforeTime)
		cache.ResetTimer(timerDuration)
	} else {
		recordCacheExecTime(beforeTime)
	}
	encodeArticles(w, articles)
	if limit < coldCapacity {
		beforeTime = time.Now()
		updateCache(cache, articles)
		recordCacheUpdateTime(beforeTime)
	}
	// go func(limit int, coldCapacity int, cache *c.Cache, articles []*c.Article) {
	// 	if limit < coldCapacity {
	// 		updateCache(cache, articles)
	// 	}
	// }(limit, coldCapacity, cache, articles)

}
func getSportsArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	beforeTime := time.Now()
	articles := cache.GetArticlesByCategory("Sports", limit, true)
	if len(articles) < limit {
		beforeTime = time.Now()
		articles, err = db.GetArticlesByCategory("Sports", limit)
		c.CheckErr(err)
		recordDBExecTime(beforeTime)
		cache.ResetTimer(timerDuration)
	} else {
		recordCacheExecTime(beforeTime)
	}
	encodeArticles(w, articles)
	if limit < coldCapacity {
		beforeTime = time.Now()
		updateCache(cache, articles)
		recordCacheUpdateTime(beforeTime)
	}
	// go func(limit int, coldCapacity int, cache *c.Cache, articles []*c.Article) {
	// 	if limit < coldCapacity {
	// 		updateCache(cache, articles)
	// 	}
	// }(limit, coldCapacity, cache, articles)

}
func getPoliticsArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	beforeTime := time.Now()
	articles := cache.GetArticlesByCategory("Politics", limit, true)
	if len(articles) < limit {
		beforeTime = time.Now()
		articles, err = db.GetArticlesByCategory("Politics", limit)
		c.CheckErr(err)
		recordDBExecTime(beforeTime)
		cache.ResetTimer(timerDuration)
	} else {
		recordCacheExecTime(beforeTime)
	}
	encodeArticles(w, articles)
	if limit < coldCapacity {
		beforeTime = time.Now()
		updateCache(cache, articles)
		recordCacheUpdateTime(beforeTime)
	}
	// go func(limit int, coldCapacity int, cache *c.Cache, articles []*c.Article) {
	// 	if limit < coldCapacity {
	// 		updateCache(cache, articles)
	// 	}
	// }(limit, coldCapacity, cache, articles)

}
func getScienceTechnologyArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	beforeTime := time.Now()
	articles := cache.GetArticlesByCategory("Politics", limit, true)
	if len(articles) < limit {
		beforeTime = time.Now()
		articles, err = db.GetArticlesByCategory("Politics", limit)
		c.CheckErr(err)
		recordDBExecTime(beforeTime)
		cache.ResetTimer(timerDuration)
	} else {
		recordCacheExecTime(beforeTime)
	}
	encodeArticles(w, articles)
	if limit < coldCapacity {
		beforeTime = time.Now()
		updateCache(cache, articles)
		recordCacheUpdateTime(beforeTime)
	}
	// go func(limit int, coldCapacity int, cache *c.Cache, articles []*c.Article) {
	// 	if limit < coldCapacity {
	// 		updateCache(cache, articles)
	// 	}
	// }(limit, coldCapacity, cache, articles)

}
func getBusinessArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	beforeTime := time.Now()
	articles := cache.GetArticlesByCategory("Politics", limit, true)
	if len(articles) < limit {
		beforeTime = time.Now()
		articles, err = db.GetArticlesByCategory("Politics", limit)
		c.CheckErr(err)
		recordDBExecTime(beforeTime)
		cache.ResetTimer(timerDuration)
	} else {
		recordCacheExecTime(beforeTime)
	}
	encodeArticles(w, articles)
	if limit < coldCapacity {
		beforeTime = time.Now()
		updateCache(cache, articles)
		recordCacheUpdateTime(beforeTime)
	}
	// go func(limit int, coldCapacity int, cache *c.Cache, articles []*c.Article) {
	// 	if limit < coldCapacity {
	// 		updateCache(cache, articles)
	// 	}
	// }(limit, coldCapacity, cache, articles)

}

//	func getFrontPageArticles(w http.ResponseWriter, r *http.Request) {
//		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
//		c.CheckErr(err)
//		cache.
//		articles, _, err := db.GetNewestArticles(limit)
//		c.CheckErr(err)
//		encodeArticles(w, articles)
//		if limit < coldCapacity {
//			updateCache(articles)
//		} else {
//			// fmt.Println("Not updating the cache because the limit is too large")
//		}
//	}
func getBreakingNewsArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	beforeTime := time.Now()
	articles := cache.GetArticlesByCategory("Breaking News", limit, true)
	if len(articles) < limit {
		beforeTime = time.Now()
		articles, err = db.GetArticlesByCategory("Breaking News", limit)
		c.CheckErr(err)
		recordDBExecTime(beforeTime)
		cache.ResetTimer(timerDuration)

	} else {
		recordCacheExecTime(beforeTime)
	}
	encodeArticles(w, articles)
	// go func(limit int, coldCapacity int, cache *c.Cache, articles []*c.Article) {
	// 	if limit < coldCapacity {
	// 		updateCache(cache, articles)
	// 	}
	// }(limit, coldCapacity, cache, articles)
	if limit < coldCapacity {
		beforeTime = time.Now()
		updateCache(cache, articles)
		recordCacheUpdateTime(beforeTime)
	}

}
func getArtsCultureArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	beforeTime := time.Now()
	articles := cache.GetArticlesByCategory("Arts and Culture", limit, true)
	if len(articles) < limit {
		beforeTime = time.Now()
		articles, err = db.GetArticlesByCategory("Arts and Culture", limit)
		c.CheckErr(err)
		recordDBExecTime(beforeTime)
		cache.ResetTimer(timerDuration)
	} else {
		recordCacheExecTime(beforeTime)
	}
	encodeArticles(w, articles)
	// go func(limit int, coldCapacity int, cache *c.Cache, articles []*c.Article) {
	if limit < coldCapacity {
		beforeTime = time.Now()
		updateCache(cache, articles)
		recordCacheUpdateTime(beforeTime)
	}
	// }(limit, coldCapacity, cache, articles)
}

func getArticleById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	c.CheckErr(err)
	beforeTime := time.Now()
	article := cache.GetArticleById(id)
	if article == nil {
		beforeTime = time.Now()
		article, err = db.GetArticleById(id)
		recordDBExecTime(beforeTime)
		c.CheckErr(err)
		cache.ResetTimer(timerDuration)
	} else {
		recordCacheExecTime(beforeTime)
	}
	articles := make([]*c.Article, 1)
	articles[0] = article
	encodeArticles(w, articles)
	// go func(cache *c.Cache, articles []*c.Article) {
	beforeTime = time.Now()
	updateCache(cache, articles)
	recordCacheUpdateTime(beforeTime)
	// }(cache, articles)
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

func updateCache(cache *c.Cache, articles []*c.Article) {
	for i := len(articles) - 1; i >= 0; i-- {
		cache.Add(articles[i])
	}
}

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
