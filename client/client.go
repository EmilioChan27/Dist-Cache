package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	c "github.com/EmilioChan27/Dist-Cache/common"
)

func main() {
	// params := url.Values{}
	// params.Add("id", fmt.Sprint(1))
	// url := "http://LX-Server:8080/"
	// // + "?" + params.Encode()
	// fmt.Println("About to make request")
	// res, err := http.Get(url)
	// fmt.Println("Just made request")
	// if err != nil {
	// 	fmt.Println("Error", err)
	// 	return
	// }
	// defer res.Body.Close()
	// body, err := io.ReadAll(res.Body)
	// if err != nil {
	// 	fmt.Println("something went wrong, ", err)
	// } else {
	// 	fmt.Println(string(body))
	// }
	file, err := os.Create("2m_output.txt")
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		fmt.Printf("In the %dth operation\n", i)
		time.Sleep(1 * time.Second)
		beforeTime := time.Now()
		var res *http.Response
		var err error
		if i%9 == 0 {
			res, err = http.Get("http://localhost:8080/human-interest?limit=100")
			fmt.Println("human interest")
		} else if i%8 == 0 {
			res, err = http.Get("http://localhost:8080/business?limit=100")
			fmt.Println("business")
		} else if i%7 == 0 {
			res, err = http.Get("http://localhost:8080/international-affairs?limit=100")
			fmt.Println("international affairs")
		} else if i%6 == 0 {
			res, err = http.Get("http://localhost:8080/science-technology?limit=100")
			fmt.Println("science and technology")
		} else if i%5 == 0 {
			res, err = http.Get("http://localhost:8080/arts-culture?limit=100")
			fmt.Println("arts and culture")
		} else if i%4 == 0 {
			res, err = http.Get("http://localhost:8080/politics?limit=100")
			fmt.Println("politics")
		} else if i%3 == 0 {
			res, err = http.Get("http://localhost:8080/sports?limit=100")
			fmt.Println("sports")
		} else if i%2 == 0 {
			res, err = http.Get("http://localhost:8080/breaking-news?limit=100")
			fmt.Println("Breaking News")
		} else {
			res, err = http.Get("http://localhost:8080/article?id=3")
			fmt.Println("Get By ID")
		}
		execTime := time.Since(beforeTime)
		file.WriteString(fmt.Sprintf("%v", execTime))
		if err != nil {
			log.Fatal(err)
		}
		dec := json.NewDecoder(res.Body)
		defer res.Body.Close()
		// articles := make([]*c.Article, 0)
		numArticles := 0
		for dec.More() {
			var a c.Article
			err = dec.Decode(&a)
			if err != nil {
				log.Fatal(err)
			}
			numArticles++
			// fmt.Printf("Id: %d\n", a.Id)
			// articles = append(articles, a)
		}
		fmt.Printf("Num articles: %d\n", numArticles)
	}
	// for _, a := range articles {
	// 	fmt.Printf("Id: %d\n", a.Id)
	// }
	// for i := 0; i < 10; i++ {
	// 	time.Sleep(6 * time.Second)
	// 	_, err := http.Get("http://localhost:8080/business")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	time.Sleep(6 * time.Second)
	// 	_, err = http.Get("http://localhost:8080/human-interest")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	time.Sleep(6 * time.Second)
	// 	_, err = http.Get("http://localhost:8080/science-technology")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	time.Sleep(6 * time.Second)
	// 	_, err = http.Get("http://localhost:8080/international-affairs")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	//	file, err := os.Create("1.375m_create_delete.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for i := 0; i < 25; i++ {
	// 	time.Sleep(83 * time.Second)
	// 	beforeTime := time.Now()
	// 	res, err := http.Get("http://localhost:8080/")
	// 	if err != nil {
	// 		log.Fatal("Error reading Employees: ", err.Error())
	// 	}
	// 	afterTime := time.Now()
	// 	executionTime := afterTime.Sub(beforeTime)
	// 	_, err = io.ReadAll(res.Body)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	str := fmt.Sprintf("Execution time: %v\n", executionTime)
	// 	res.Body.Close()
	// 	file.WriteString(str)
	// 	// fmt.Println(str)
	// }
	// resChanLen := 100
	// resChan := make(chan string, resChanLen)
	// for i := 0; i < resChanLen; i++ {
	// 	go func(index int) {
	// 		params := url.Values{}
	// 		params.Add("id", fmt.Sprint(index))
	// 		url := "http://LX-Cache:8080/" + "?" + params.Encode()
	// 		res, err := http.Get(url)
	// 		if err != nil {
	// 			fmt.Println("Error", err)
	// 			return
	// 		}
	// 		defer res.Body.Close()
	// 		body, err := io.ReadAll(res.Body)
	// 		if err != nil {
	// 			fmt.Println("something went wrong, ", err)
	// 		}
	// 		// fmt.Println(string(body))
	// 		resChan <- string(body)
	// 		// fmt.Printf("received output from server from request %d\n", index)
	// 	}(i)
	// }
	// for i := 0; i < resChanLen; i++ {
	// 	s := <-resChan
	// 	fmt.Println(s)
	// }
	// for s := range resChan {
	// 	fmt.Println(s)
	// 	if len(resChan) == 0 {
	// 		break
	// 	}
	// }
	// resReader := bufio.NewReader(res.Body)
	// for {
	// 	line, err := resReader.ReadBytes('\n')

	// 	if err != nil {
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		fmt.Println("Error reading line:", err)
	// 		break
	// 	}
	// 	fmt.Println("Received Event: ", string(line))
	// }
}
