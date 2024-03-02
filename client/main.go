package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func main() {
	resChanLen := 100
	resChan := make(chan string, resChanLen)
	for i := 0; i < resChanLen; i++ {
		go func(index int) {
			params := url.Values{}
			params.Add("id", fmt.Sprint(index))
			url := "http://localhost:8888/" + "?" + params.Encode()
			res, err := http.Get(url)
			if err != nil {
				fmt.Println("Error", err)
				return
			}
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Println("something went wrong, ", err)
			}
			// fmt.Println(string(body))
			resChan <- string(body)
			// fmt.Printf("received output from server from request %d\n", index)
		}(1)
	}
	for i := 0; i < resChanLen; i++ {
		s := <-resChan
		fmt.Println(s)
	}
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
