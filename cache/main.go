package main

import (
	"fmt"
	"net/http"
)

func main() {
	// keyVals := make(map[string]string)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Access-Control-Allow-Origin", "*")

		id := "1"
		// val, found := keyVals[id]
		// if found {
		val := "1"
		fmt.Printf("responding to request with number %s\n", id)
		fmt.Fprintf(w, "cached response to %s: %s", id, val)
		// } else {
		// 	// keyVals[id] = id
		// 	fmt.Printf("Forwarding request with number %s\n", id)
		// 	serverUrl, err := url.Parse("http://localhost:8080/")
		// 	if err != nil {
		// 		fmt.Println("Error was not nil")
		// 	}
		// 	server := httputil.NewSingleHostReverseProxy(serverUrl)
		// 	server.ServeHTTP(w, r)
		// }
	})
	fmt.Println("Cache server running on port :80")
	http.ListenAndServe(":80", nil)
}
