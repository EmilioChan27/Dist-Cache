package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	c "github.com/EmilioChan27/Dist-Cache/common"
	d "github.com/EmilioChan27/Dist-Cache/db"
)

// var numRequests int = 0
var db *d.DB
var cache *c.Cache
var coldCapacity int
var hotCapacity int

func main() {
	db = d.NewDB()
	coldCapacity = 50
	hotCapacity = 50
	cache = c.NewCache(coldCapacity, hotCapacity, 1)
	// // articles := make([]*c.Article, 10)
	for i := 0; i < 100; i++ {
		var a *c.Article
		if i < 50 {
			a = &c.Article{Id: i, Category: "Business", Content: "Today I went to the park", AuthorId: i}
		} else {
			a = &c.Article{Id: i, Category: "Human Interest", Content: "Today I went to the park", AuthorId: i}
		}
		cache.Add(a)
	}
	// for i := 0; i < 10; i++ {

	// 	go func() {
	// 		db.InsertTestArticles("../cohere/output.txt", 99)

	// 	}()
	// }
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

	// newId, err := db.AddArticle(&c.Article{Title: "Title 26", Content: "Random content for this article", AuthorId: 1, ImageUrl: "Random ImageUrl", Category: "Business"})
	// c.CheckErr(err)
	// file.WriteString("Newid: %d\n", newId)
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	numRequests++
	// 	file.WriteString("Received %d requests\n", numRequests)
	// 	id := r.URL.Query().Get("id")
	// 	file.WriteString("processing request %s\n", id)
	// 	fmt.Fprintf(w, "Server response to %s: %s", id, id)
	// })
	// cache.ToString()
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
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
func getHumanInterestArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	articles := cache.GetArticlesByCategory("Human Interest", limit, true)
	if len(articles) < limit {
		articles, err = db.GetArticlesByCategory("Human Interest", limit)
		c.CheckErr(err)
	}
	encodeArticles(w, articles)
	if limit < coldCapacity {
		updateCache(articles)
	} else {
		// fmt.Println("Not updating the cache because the limit is too large")
	}
}
func getInternationalAffairsArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	articles := cache.GetArticlesByCategory("International Affairs", limit, true)
	if len(articles) < limit {
		articles, err = db.GetArticlesByCategory("International Affairs", limit)
		c.CheckErr(err)
	}
	encodeArticles(w, articles)
	if limit < coldCapacity {
		updateCache(articles)
	} else {
		// fmt.Println("Not updating the cache because the limit is too large")
	}
}
func getSportsArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	articles := cache.GetArticlesByCategory("Sports", limit, true)
	if len(articles) < limit {
		articles, err = db.GetArticlesByCategory("Sports", limit)
		c.CheckErr(err)
	}
	encodeArticles(w, articles)
	if limit < coldCapacity {
		updateCache(articles)
	} else {
		// fmt.Println("Not updating the cache because the limit is too large")
	}
}
func getPoliticsArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	articles := cache.GetArticlesByCategory("Politics", limit, true)
	if len(articles) < limit {
		articles, err = db.GetArticlesByCategory("Politics", limit)
		c.CheckErr(err)
	}
	encodeArticles(w, articles)
	if limit < coldCapacity {
		updateCache(articles)
	} else {
		// fmt.Println("Not updating the cache because the limit is too large")
	}
}
func getScienceTechnologyArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	articles := cache.GetArticlesByCategory("Politics", limit, true)
	if len(articles) < limit {
		articles, err = db.GetArticlesByCategory("Politics", limit)
		c.CheckErr(err)
	}
	encodeArticles(w, articles)
	if limit < coldCapacity {
		updateCache(articles)
	} else {
		// fmt.Println("Not updating the cache because the limit is too large")
	}
}
func getBusinessArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	articles := cache.GetArticlesByCategory("Politics", limit, true)
	if len(articles) < limit {
		articles, err = db.GetArticlesByCategory("Politics", limit)
		c.CheckErr(err)
	}
	encodeArticles(w, articles)
	if limit < coldCapacity {
		updateCache(articles)
	} else {
		// fmt.Println("Not updating the cache because the limit is too large")
	}
}
func getFrontPageArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	articles := cache.GetArticlesByCategory("Politics", limit, true)
	if len(articles) < limit {
		articles, err = db.GetArticlesByCategory("Politics", limit)
		c.CheckErr(err)
	}
	encodeArticles(w, articles)
	if limit < coldCapacity {
		updateCache(articles)
	} else {
		// fmt.Println("Not updating the cache because the limit is too large")
	}
}
func getBreakingNewsArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	articles := cache.GetArticlesByCategory("Breaking News", limit, true)
	if len(articles) < limit {
		articles, err = db.GetArticlesByCategory("Breaking News", limit)
		c.CheckErr(err)
	}
	encodeArticles(w, articles)
	if limit < coldCapacity {
		updateCache(articles)
	} else {
		// fmt.Println("Not updating the cache because the limit is too large")
	}
}
func getArtsCultureArticles(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	c.CheckErr(err)
	articles := cache.GetArticlesByCategory("Arts and Culture", limit, true)
	if len(articles) < limit {
		articles, err = db.GetArticlesByCategory("Arts and Culture", limit)
		c.CheckErr(err)
	}
	encodeArticles(w, articles)
	if limit < coldCapacity {
		updateCache(articles)
	} else {
		// fmt.Println("Not updating the cache because the limit is too large")
	}
}

func getArticleById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	c.CheckErr(err)
	article := cache.GetArticleById(id)
	if article == nil {
		article, err = db.GetArticleById(id)
		c.CheckErr(err)
	}
	articles := make([]*c.Article, 1)
	articles[0] = article
	encodeArticles(w, articles)
	updateCache(articles)
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

func updateCache(articles []*c.Article) {
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
