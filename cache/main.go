package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	keyVals := make(map[string]string)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		val, found := keyVals[id]
		if found {
			fmt.Printf("responding to request with number %s\n", id)
			fmt.Fprintf(w, "cached response to %s: %s", id, val)
		} else {
			keyVals[id] = id
			fmt.Printf("Forwarding request with number %s\n", id)
			serverUrl, err := url.Parse("http://172.208.52.232:8080/")
			if err != nil {
				fmt.Println("Error was not nil")
			}
			server := httputil.NewSingleHostReverseProxy(serverUrl)
			server.ServeHTTP(w, r)
		}
	})
	fmt.Println("Cache server running on port :8888")
	http.ListenAndServe("0.0.0.0:8888", nil)
}
