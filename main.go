package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type UserInfo struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserAndPostsInfo struct {
	UserId   int      `json:"id"`
	UserInfo UserInfo `json:"userInfo"`
	Posts    []string `json:"posts"`
}

type UserPostsService interface {
	GetUserInfo(userId int) (UserInfo, error)
}

func handleUserPosts(service UserPostsService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pathSegments := strings.Split(r.URL.Path, "/")
		// there's certainly a better, more canonical method for getting the user ID from the path
		userId, err := strconv.Atoi(pathSegments[len(pathSegments)-1])
		if err != nil {
			log.Printf("error when attempting to parse user ID from URI: %v", err)
			// TODO: add response
			return
		}

		userInfo, err := service.GetUserInfo(userId)

		js, err := json.Marshal(UserAndPostsInfo{
			userId,
			userInfo,
			make([]string, 0),
		})
		fmt.Fprint(w, string(js))
	}
}

func main() {
	http.HandleFunc("/v1/user-posts/", handleUserPosts(nil))
	log.Fatal(http.ListenAndServe(":8081", nil))
}
