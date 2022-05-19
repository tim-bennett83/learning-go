package main

import (
	"fmt"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/v1/user-posts/", func(w http.ResponseWriter, r *http.Request) {
		pathSegments := strings.Split(r.URL.Path, "/")
		// there's certainly a better, more canonical method for getting the user ID from the path
		userID := pathSegments[len(pathSegments) - 1]
		fmt.Fprintf(w, `{"userID": %q}`, userID)
	})

	http.ListenAndServe(":8081", nil)
}
