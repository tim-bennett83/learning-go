package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

// setup to allow a range of test users
var userInfoMap = map[int]UserInfo{
	42: {"Leanne Graham", "Bret", "Sincere@april.biz"},
}
var postsMap = map[int][]PostInfo{
	42: {{57, "Dissertation on the Weave", "work in progress..."}, {91, "Hobgoblins and You", "Goblinoids are fascinating creatures..."}},
}

// mock implementation of the service interface
type UserPostsTestService struct{}

func (s *UserPostsTestService) GetUserInfo(userId int, userChan chan chanResult[UserInfo]) {
	info, found := userInfoMap[userId]
	if !found {
		userChan <- chanResult[UserInfo]{UserInfo{}, UserNotFoundError{userId}}
		return
	}
	userChan <- chanResult[UserInfo]{info, nil}
}

func (s *UserPostsTestService) GetPostsForUser(userId int, postsChan chan chanResult[[]PostInfo]) {
	posts, found := postsMap[userId]
	if !found {
		postsChan <- chanResult[[]PostInfo]{nil, UserNotFoundError{userId}}
		return
	}
	postsChan <- chanResult[[]PostInfo]{posts, nil}
}

func TestMain(t *testing.T) {
	t.Run("receives user info with posts", func(t *testing.T) {
		response := httptest.NewRecorder()
		handleUserPosts(&UserPostsTestService{})(response, httptest.NewRequest("GET", "/v1/user-posts/42", nil))
		defer response.Result().Body.Close()

		gotBytes, err := ioutil.ReadAll(response.Result().Body)
		if err != nil {
			t.Errorf("got error %v", err)
		}
		var got UserAndPostsInfo
		json.Unmarshal(gotBytes, &got)
		want := UserAndPostsInfo{42, userInfoMap[42], postsMap[42]}

		if got.UserId != want.UserId || got.UserInfo != want.UserInfo {
			t.Errorf("expected %v got %v", want, got)
		}
	})

	t.Run("returns a 404 if upstream service returns a not found error", func(t *testing.T) {
		response := httptest.NewRecorder()
		handleUserPosts(&UserPostsTestService{})(response, httptest.NewRequest("GET", "/v1/user-posts/43", nil))
		defer response.Result().Body.Close()

		if response.Code != 404 {
			t.Errorf("Expected response code 404, got %d", response.Code)
		}
	})

	t.Run("returns a 400 if the user ID is not a number", func(t *testing.T) {
		response := httptest.NewRecorder()
		handleUserPosts(&UserPostsTestService{})(response, httptest.NewRequest("GET", "/v1/user-posts/asdf", nil))
		defer response.Result().Body.Close()

		if response.Code != 400 {
			t.Errorf("Expected response code 404, got %d", response.Code)
		}
	})
}
