package main

import (
	"fmt"
	"net/http"
)

var numRequests int = 0

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		numRequests++
		id := "1"
		fmt.Printf("processing request %s\n", id)
		fmt.Fprintf(w, "Server response to %s: %s", id, id)
	})
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
