package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
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
	file, err := os.Create("1.375m_create_delete.txt")
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 25; i++ {
		time.Sleep(83 * time.Second)
		beforeTime := time.Now()
		res, err := http.Get("http://localhost:8080/")
		if err != nil {
			log.Fatal("Error reading Employees: ", err.Error())
		}
		afterTime := time.Now()
		executionTime := afterTime.Sub(beforeTime)
		_, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		str := fmt.Sprintf("Execution time: %v\n", executionTime)
		res.Body.Close()
		file.WriteString(str)
		// fmt.Println(str)
	}
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
