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

type PostInfo struct {
	PostId int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type UserAndPostsInfo struct {
	UserId   int        `json:"id"`
	UserInfo UserInfo   `json:"userInfo"`
	Posts    []PostInfo `json:"posts"`
}

type chanResult[T any] struct {
	value T
	err   error
}

type UserPostsService interface {
	GetUserInfo(userId int, userChan chan chanResult[UserInfo])
	GetPostsForUser(userId int, postsChan chan chanResult[[]PostInfo])
}

type UserPostsServiceImpl struct{}

type UserNotFoundError struct {
	userId int
}

func (e UserNotFoundError) Error() string {
	return fmt.Sprintf("User %d not found", e.userId)
}

func (s *UserPostsServiceImpl) GetUserInfo(userId int, userChan chan chanResult[UserInfo]) {
	defer close(userChan)
	res, err := http.Get(fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%d", userId))
	if err != nil {
		log.Printf("error fetching user: %v", err)
		userChan <- chanResult[UserInfo]{UserInfo{}, err}
		return
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		userChan <- chanResult[UserInfo]{UserInfo{}, UserNotFoundError{userId}}
		return
	}

	dec := json.NewDecoder(res.Body)
	var userInfo UserInfo
	err = dec.Decode(&userInfo)
	if err != nil {
		userChan <- chanResult[UserInfo]{UserInfo{}, err}
		return
	}
	userChan <- chanResult[UserInfo]{userInfo, nil}
}

func (s *UserPostsServiceImpl) GetPostsForUser(userId int, postsChan chan chanResult[[]PostInfo]) {
	defer close(postsChan)
	res, err := http.Get(fmt.Sprintf("https://jsonplaceholder.typicode.com/posts?userId=%d", userId))
	if err != nil {
		log.Printf("error fetching posts for user: %v", err)
		postsChan <- chanResult[[]PostInfo]{nil, err}
		return
	}
	defer res.Body.Close()

	dec := json.NewDecoder(res.Body)
	var posts []PostInfo
	err = dec.Decode(&posts)
	if err != nil {
		postsChan <- chanResult[[]PostInfo]{nil, err}
		return
	}
	postsChan <- chanResult[[]PostInfo]{posts, nil}
}

func handleServiceError(err error, w http.ResponseWriter) {
	switch err.(type) {
	case UserNotFoundError:
		w.WriteHeader(404)
	default:
		log.Printf("unexpected error: %v", err)
		w.WriteHeader(500)
	}
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

		userChan := make(chan chanResult[UserInfo])
		postsChan := make(chan chanResult[[]PostInfo])

		go service.GetUserInfo(userId, userChan)
		go service.GetPostsForUser(userId, postsChan)

		userRes := <-userChan
		postsRes := <-postsChan
		userInfo := userRes.value
		posts := postsRes.value

		if userRes.err != nil {
			handleServiceError(userRes.err, w)
			return
		}
		if postsRes.err != nil {
			handleServiceError(postsRes.err, w)
			return
		}

		js, _ := json.Marshal(UserAndPostsInfo{
			userId,
			userInfo,
			posts,
		})
		fmt.Fprint(w, string(js))
	}
}

func main() {
	http.HandleFunc("/v1/user-posts/", handleUserPosts(&UserPostsServiceImpl{}))
	log.Fatal(http.ListenAndServe(":8081", nil))
}
