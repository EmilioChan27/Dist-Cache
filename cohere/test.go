package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	output, err := os.Open("output.txt")
	if err != nil {
		log.Fatal(err)
	}
	fileInfo, err := output.Stat()
	if err != nil {
		log.Fatal(err)
	}
	length := fileInfo.Size()
	byteOutput := make([]byte, length)
	_, err = output.Read(byteOutput)
	if err != nil {
		log.Fatal(err)
	}
	strOutput := string(byteOutput)
	strArrOutput := strings.Split(strOutput, "\\")
	for _, str := range strArrOutput {
		fmt.Println("------------------------------")
		titleUnTruncated := strings.Split(strings.Split(str, ":")[1], "\n")[0]
		fmt.Printf("Title: %s\n", titleUnTruncated[:len(titleUnTruncated)-4])
		fmt.Printf("Category: %s\n", strings.Split(str, ":")[0][4:])
		// fmt.Println(str)
		fmt.Println("------------------------------")
	}
}
