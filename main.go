package main

import (
	"fmt"
	"net/http"
	"strings"
)

func handleUserPosts() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pathSegments := strings.Split(r.URL.Path, "/")
		// there's certainly a better, more canonical method for getting the user ID from the path
		userID := pathSegments[len(pathSegments)-1]
		fmt.Fprintf(w, `{"userID": %q}`, userID)
	}
}

func main() {
	http.HandleFunc("/v1/user-posts/", handleUserPosts())
	http.ListenAndServe(":8081", nil)
}
