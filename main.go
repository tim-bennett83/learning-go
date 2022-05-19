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

type UserPostsServiceImpl struct{}

type UserNotFoundError struct{
	userId int
}
func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("User %d not found", e.userId)
}

func (s *UserPostsServiceImpl) GetUserInfo(userId int) (userInfo UserInfo, err  error) {
	res, err := http.Get(fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%d", userId))
	if err != nil {
		log.Printf("error fetching user: %v", err)
		return UserInfo{}, err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return UserInfo{}, UserNotFoundError{userId}
	}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&userInfo)
	return userInfo, err
}

func handleUserPosts(service UserPostsService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pathSegments := strings.Split(r.URL.Path, "/")
		// there must be a better, more canonical method for getting the user ID from the path
		userId, err := strconv.Atoi(pathSegments[len(pathSegments)-1])
		if err != nil {
			log.Printf("error when attempting to parse user ID from URI: %v", err)
			w.WriteHeader(400)
			return
		}

		userInfo, err := service.GetUserInfo(userId)

		if err != nil {
			switch err.(type) {
			case UserNotFoundError:
				w.WriteHeader(404)
			default:
				log.Printf("unexpected error: %v", err)
				w.WriteHeader(500)
			}
			return
		}

		js, _ := json.Marshal(UserAndPostsInfo{
			userId,
			userInfo,
			make([]string, 0),
		})
		fmt.Fprint(w, string(js))
	}
}

func main() {
	http.HandleFunc("/v1/user-posts/", handleUserPosts(&UserPostsServiceImpl{}))
	log.Fatal(http.ListenAndServe(":8081", nil))
}
