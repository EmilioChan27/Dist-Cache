package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

func main() {
	// keyVals := make(map[string]string)
	keyVals := new(sync.Map)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Access-Control-Allow-Origin", "*")

		id := r.URL.Query().Get("id")
		val, found := keyVals.Load(id)
		if found {

			fmt.Printf("responding to request with number %s\n", id)
			fmt.Fprintf(w, "cached response to %s: %s", id, val)
		} else {
			keyVals.Store(id, id)
			fmt.Printf("Forwarding request with number %s\n", id)
			serverUrl, err := url.Parse("http://LX-Server:8080/")
			if err != nil {
				fmt.Println("Error was not nil")
			}
			server := httputil.NewSingleHostReverseProxy(serverUrl)
			server.ServeHTTP(w, r)
		}
	})
	fmt.Println("Cache server running on port :8080")
	http.ListenAndServe("0.0.0.0:8080", nil)
}
