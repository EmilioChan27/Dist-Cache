package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	c "github.com/EmilioChan27/Dist-Cache/common"
	d "github.com/EmilioChan27/Dist-Cache/db"
)

// var numRequests int = 0
var db *d.DB

func main() {
	db = d.NewDB()
	// db.InsertTestArticles("../cohere/output.txt", 99)
	http.HandleFunc("/front-page", getFrontPage)
	http.HandleFunc("/business", getBusinessArticles)
	http.HandleFunc("/human-interest", getHumanInterestArticles)
	http.HandleFunc("/international-affairs", getInternationalAffairsArticles)
	http.HandleFunc("/science-technology", getScienceTechnologyArticles)
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
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
func getFrontPage(w http.ResponseWriter, r *http.Request) {

}

func getBusinessArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := db.GetArticlesByCategory("Business")
	c.CheckErr(err)
	for _, a := range articles {
		writeArticleToFile(a, "output.txt")
	}
}
func getHumanInterestArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := db.GetArticlesByCategory("Human Interest")
	c.CheckErr(err)
	for _, a := range articles {
		writeArticleToFile(a, "output.txt")
	}
}
func getInternationalAffairsArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := db.GetArticlesByCategory("International Affairs")
	c.CheckErr(err)
	for _, a := range articles {
		writeArticleToFile(a, "output.txt")
	}
}
func getScienceTechnologyArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := db.GetArticlesByCategory("Science and Technology")
	c.CheckErr(err)
	for _, a := range articles {
		writeArticleToFile(a, "output.txt")
	}
}

func writeArticleToFile(a *c.Article, filename string) {
	var file *os.File
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		file, err = os.Create(filename)
	} else {
		file, err = os.Open(filename)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString("----------------------")
	file.WriteString(fmt.Sprintf("Title: %s\n", a.Title))
	file.WriteString(fmt.Sprintf("Category: %s\n", a.Category))
	file.WriteString(fmt.Sprintf("Author: %d\n", a.AuthorId))
	file.WriteString(fmt.Sprintf("Content Preview: %s\n", a.Content[:10]))
	file.WriteString(fmt.Sprintf("Created At: %v\n", a.CreatedAt))
	file.WriteString("-------------------------------")
}
