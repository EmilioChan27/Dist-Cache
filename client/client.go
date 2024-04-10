package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"

	c "github.com/EmilioChan27/Dist-Cache/common"
)

func main() {
	actualTest(25, 15*time.Minute)
	// getArticleById(2)
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

func getArticleById(id int) *http.Response {
	serverUrl := "http://LX-Server:8080/"
	// fmt.Println("------------------")
	var res *http.Response
	urlEnd := fmt.Sprintf("article?id=%d", id)
	res, err := http.Get(serverUrl + urlEnd)
	// fmt.Println("Get By ID")

	if err != nil {
		log.Fatal(err)
	}
	return res
	// dec := json.NewDecoder(res.Body)

	// defer res.Body.Close()
	// // articles := make([]*c.Article, 0)
	// numArticles := 0
	// for dec.More() {
	// 	var a *c.Article
	// 	err = dec.Decode(&a)
	// 	fmt.Println(a)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	numArticles++
	// 	// fmt.Printf("Id: %d\n", a.Id)
	// 	// articles = append(articles, a)
	// }
	// fmt.Printf("Num articles: %d\n", numArticles)

}

func latencyTest() {
	serverUrl := "http://localhost:8080/"
	fmt.Println("------------------")
	var res *http.Response
	a := &c.Article{Id: 50000, AuthorId: 1, Content: "contentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontent", Category: "International Affairs", Title: "Testing article", ImageUrl: "Random Image", Likes: 0, Size: 25}
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

func actualTest(numClients int, testDuration time.Duration) {
	clients := make(chan int, numClients)
	writes := make(chan int, 100)
	overallTimer := time.NewTimer(testDuration)
	maxId := 51476
	src := rand.NewSource(int64(maxId))
	zipf := rand.NewZipf(rand.New(src), 1.5, 8, uint64(maxId))
	actualNumClients := 0
	file, err := os.Create(fmt.Sprintf("%dclients-%vduration-pause2+65s-cache.txt", numClients, testDuration))
	c.CheckErr(err)
	file.WriteString("overallTimer := time.NewTimer(testDuration)\nmaxId := 51476\nsrc := rand.NewSource(int64(maxId))\nzipf := rand.NewZipf(rand.New(src), 1.5, 8, uint64(maxId))\n")
	file.WriteString("Hello????")
outerlabel:
	for {
		select {
		case <-overallTimer.C:
			break outerlabel
		case <-clients:
			go func(zipf *rand.Zipf, maxId int, file *os.File) {
				waitTime := int(math.Abs(rand.NormFloat64()*2 + 2))
				id := maxId - int(zipf.Uint64())
				for i := 0; i < waitTime; i++ {
					time.Sleep(1 * time.Second)
				}
				beforeTime := time.Now()
				_ = getArticleById(id)
				execTime := time.Since(beforeTime).Milliseconds()
				execTimeString := fmt.Sprintf("%v\n", execTime)
				// if execTimeString[len(execTimeString)-2:] != "ms" && execTimeString[len(execTimeString)-2:] != "μs" {
				// 	ms, err := strconv.ParseFloat(execTimeString[:len(execTimeString)-2], 32)
				// 	c.CheckErr(err)
				// 	seconds := ms * 1000
				// 	execTimeString = fmt.Sprintf("%v\n", seconds)
				// } else if execTimeString[len(execTimeString)-2:] == "μs" {
				// 	us, err := strconv.ParseFloat(execTimeString[:len(execTimeString)-2], 32)
				// 	c.CheckErr(err)
				// 	ms := us / 1000.0
				// 	execTimeString = fmt.Sprintf("%v\n", ms)
				// } else {
				// 	execTimeString = execTimeString[:len(execTimeString)-2] + "\n"
				// }
				// _ = insertArticle()
				file.WriteString(execTimeString)
				clients <- 1
			}(zipf, maxId, file)
		case <-writes:
			maxId++
		default:
			if actualNumClients < numClients {
				time.Sleep(10 * time.Millisecond)
				clients <- 1
				actualNumClients++
				fmt.Printf("Current numClients: %d\n", actualNumClients)
			}
		}
	}

}

// func zipfTest() {
// 	writeChan := make(chan int, 1000)
// 	longQueryChan := make(chan int, 1000)
// 	file, err := os.Create("cache-basicZipf-51475-1.5-8-51475-250u-5m-5articles-cold.txt")
// 	c.CheckErr(err)
// 	for i := 0; i < 250; i++ {
// 		go func(zipf *rand.Zipf, maxId uint64, i int, file *os.File) {
// 			localMaxId := maxId
// 			waitTime := int(math.Abs(rand.NormFloat64()*6 + 2))
// 			for k := 0; k < waitTime; k++ {
// 				time.Sleep(15 * time.Second)
// 			}
// 			numArticles := int(math.Abs(rand.NormFloat64()*5+2)) + 1
// 			// fmt.Printf("Num articles: %d\n", numArticles)
// 			for j := 0; j < numArticles; j++ {
// 				id := localMaxId - zipf.Uint64()
// 				// fmt.Printf("Id: %d\n", id)
// 				beforeTime := time.Now()
// 				_ = getArticleById(int(id))
// 				if id < 48000 {
// 					_ = getSectionArticles()
// 					longQueryChan <- 1
// 					fmt.Println("adding a long query")
// 				}
// 				if id < 5000 {
// 					insertArticle()
// 					writeChan <- 1
// 					localMaxId++
// 					fmt.Println("adding a write")
// 				}
// 				file.WriteString(fmt.Sprintf("%v\n", time.Since(beforeTime)))
// 				mult := int(math.Abs(rand.NormFloat64()*30 + 60))
// 				// fmt.Printf("Sleep time: %d\n", mult)
// 				for k := 0; k < mult; k++ {
// 					time.Sleep(time.Second)
// 				}
// 			}

// 		}(zipf, uint64(maxId), i, file)
// 	}
// 	time.Sleep(6 * time.Minute)
// 	file.WriteString(fmt.Sprintf("Num Writes: %d\n", len(writeChan)))
// 	file.WriteString(fmt.Sprintf("Num long queries: %d\n", len(longQueryChan)))
// }

func insertArticle() *http.Response {
	serverUrl := "http://localhost:8080/"
	a := &c.Article{Id: 1, AuthorId: 1, Content: "contentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontentcontent", Category: "Breaking News", Title: "Testing article", ImageUrl: "Random Image", Likes: 0, Size: 25}
	// values := map[string]interface{}{"Id": fmt.Sprintf("%d", a.Id), "AuthorId": fmt.Sprintf("%d", a.AuthorId), "Content": a.Content, "Category": a.Category, "Title": a.Title, "ImageUrl": a.ImageUrl, "Likes": fmt.Sprintf("%d", a.Likes), "Size": fmt.Sprintf("%d", a.Size)}
	values := map[string]interface{}{"Id": 1, "AuthorId": a.AuthorId, "Content": a.Content, "Category": a.Category, "Title": a.Title, "ImageUrl": a.ImageUrl, "Likes": a.Likes, "Size": a.Size}
	jsonValue, _ := json.Marshal(values)
	res, err := http.Post(serverUrl+"write", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("res: %v\n", res)
	return res
}

func getSectionArticles() *http.Response {
	serverUrl := "http://localhost:8080/"
	var res *http.Response
	var err error
	i := rand.Intn(10)
	if i%9 == 0 {
		res, err = http.Get(serverUrl + "human-interest?limit=75")
		// fmt.Println("human interest")
	} else if i%8 == 0 {
		res, err = http.Get(serverUrl + "business?limit=75")
		// fmt.Println("business")
	} else if i%7 == 0 {
		res, err = http.Get(serverUrl + "international-affairs?limit=75")
		// fmt.Println("international affairs")
	} else if i%6 == 0 {
		res, err = http.Get(serverUrl + "science-technology?limit=75")
		// fmt.Println("science and technology")
	} else if i%5 == 0 {
		res, err = http.Get(serverUrl + "arts-culture?limit=75")
		// fmt.Println("arts and culture")
	} else if i%4 == 0 {
		res, err = http.Get(serverUrl + "politics?limit=75")
		// fmt.Println("politics")
	} else if i%3 == 0 {
		res, err = http.Get(serverUrl + "sports?limit=75")
		// fmt.Println("sports")
	} else {
		res, err = http.Get(serverUrl + "breaking-news?limit=75")
		// fmt.Println("Breaking News")
	}
	c.CheckErr(err)
	return res
}
