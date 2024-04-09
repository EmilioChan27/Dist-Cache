package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	c "github.com/EmilioChan27/Dist-Cache/common"
)

func main() {
	getArticleById(2)
	// latencyTest()
}

// func editArticleById(a *c.Article) {
// 	serverUrl := "http://localhost:8080/"
// 	fmt.Println("------------------")
// 	var res *http.Response
// 	a = &c.Article{Id: a.Id, AuthorId: a.AuthorId, Content: "new-content", Category: a.Category, Title: a.Title, ImageUrl: a.ImageUrl, Likes: a.Likes, Size: a.Size}
// 	// values := map[string]interface{}{"Id": fmt.Sprintf("%d", a.Id), "AuthorId": fmt.Sprintf("%d", a.AuthorId), "Content": a.Content, "Category": a.Category, "Title": a.Title, "ImageUrl": a.ImageUrl, "Likes": fmt.Sprintf("%d", a.Likes), "Size": fmt.Sprintf("%d", a.Size)}
// 	values := map[string]interface{}{"Id": a.Id, "AuthorId": a.AuthorId, "Content": a.Content, "Category": a.Category, "Title": a.Title, "ImageUrl": a.ImageUrl, "Likes": a.Likes, "Size": a.Size}
// 	jsonValue, _ := json.Marshal(values)
// 	res, err := http.(serverUrl+"write", "applicatio/json", bytes.NewBuffer(jsonValue))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("res: %v\n", res)
// }

func getArticleById(id int) {
	serverUrl := "http://localhost:8080/"
	fmt.Println("------------------")
	var res *http.Response
	urlEnd := fmt.Sprintf("article?id=%d", id)
	res, err := http.Get(serverUrl + urlEnd)
	fmt.Println("Get By ID")

	if err != nil {
		log.Fatal(err)
	}
	dec := json.NewDecoder(res.Body)
	defer res.Body.Close()
	// articles := make([]*c.Article, 0)
	numArticles := 0
	for dec.More() {
		var a *c.Article
		err = dec.Decode(&a)
		fmt.Println(a)
		if err != nil {
			log.Fatal(err)
		}
		numArticles++
		// fmt.Printf("Id: %d\n", a.Id)
		// articles = append(articles, a)
	}
	fmt.Printf("Num articles: %d\n", numArticles)

}

func latencyTest() {
	serverUrl := "http://localhost:8080/"
	fmt.Println("------------------")
	var res *http.Response
	a := &c.Article{Id: 50000, AuthorId: 1, Content: "content", Category: "International Affairs", Title: "Testing article", ImageUrl: "Random Image", Likes: 0, Size: 25}
	// values := map[string]interface{}{"Id": fmt.Sprintf("%d", a.Id), "AuthorId": fmt.Sprintf("%d", a.AuthorId), "Content": a.Content, "Category": a.Category, "Title": a.Title, "ImageUrl": a.ImageUrl, "Likes": fmt.Sprintf("%d", a.Likes), "Size": fmt.Sprintf("%d", a.Size)}
	values := map[string]interface{}{"Id": a.Id, "AuthorId": a.AuthorId, "Content": a.Content, "Category": a.Category, "Title": a.Title, "ImageUrl": a.ImageUrl, "Likes": a.Likes, "Size": a.Size}
	jsonValue, _ := json.Marshal(values)
	res, err := http.Post(serverUrl+"write", "applicatio/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("res: %v\n", res)
	// fmt.Println("Waiting 5 minutes")
	// fmt.Println("------------------")
	// for i := 0; i < 11; i++ {
	// 	fmt.Printf("In the %dth operation\n", i)
	// 	time.Sleep(450 * time.Second)
	// 	// beforeTime := time.Now()
	// 	var res *http.Response
	// 	var err error
	// 	if i%9 == 0 {
	// 		res, err = http.Get(serverUrl + "human-interest?limit=100")
	// 		fmt.Println("human interest")
	// 	} else if i%8 == 0 {
	// 		res, err = http.Get(serverUrl + "business?limit=100")
	// 		fmt.Println("business")
	// 	} else if i%7 == 0 {
	// 		res, err = http.Get(serverUrl + "international-affairs?limit=100")
	// 		fmt.Println("international affairs")
	// 	} else if i%6 == 0 {
	// 		res, err = http.Get(serverUrl + "science-technology?limit=100")
	// 		fmt.Println("science and technology")
	// 	} else if i%5 == 0 {
	// 		res, err = http.Get(serverUrl + "arts-culture?limit=100")
	// 		fmt.Println("arts and culture")
	// 	} else if i%4 == 0 {
	// 		res, err = http.Get(serverUrl + "politics?limit=100")
	// 		fmt.Println("politics")
	// 	} else if i%3 == 0 {
	// 		res, err = http.Get(serverUrl + "sports?limit=100")
	// 		fmt.Println("sports")
	// 	} else if i%2 == 0 {
	// 		res, err = http.Get(serverUrl + "breaking-news?limit=100")
	// 		fmt.Println("Breaking News")
	// 	} else {
	// 		res, err = http.Get(serverUrl + "article?id=3")
	// 		fmt.Println("Get By ID")
	// 	}
	// 	// execTime := time.Since(beforeTime)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	dec := json.NewDecoder(res.Body)
	// 	defer res.Body.Close()
	// 	// articles := make([]*c.Article, 0)
	// 	numArticles := 0
	// 	for dec.More() {
	// 		var a c.Article
	// 		err = dec.Decode(&a)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		numArticles++
	// 		// fmt.Printf("Id: %d\n", a.Id)
	// 		// articles = append(articles, a)
	// 	}
	// 	fmt.Printf("Num articles: %d\n", numArticles)
	// }
}
