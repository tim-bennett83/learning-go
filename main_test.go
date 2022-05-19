package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

type UserPostsTestService struct{}

func (s *UserPostsTestService) GetUserInfo(userId int) (info UserInfo, err error) {
	userInfoMap := map[int]UserInfo{
		42: {"Leanne Graham", "Bret", "Sincere@april.biz"},
	}
	info, found := userInfoMap[userId]
	if !found {
		return info, UserNotFoundError{userId}
	}
	return info, err
}

func TestMain(t *testing.T) {
	t.Run("receives basic user info", func(t *testing.T) {
		response := httptest.NewRecorder()
		handleUserPosts(&UserPostsTestService{})(response, httptest.NewRequest("GET", "/v1/user-posts/42", nil))
		defer response.Result().Body.Close()

		gotBytes, err := ioutil.ReadAll(response.Result().Body)
		if err != nil {
			t.Errorf("got error %v", err)
		}
		var got UserAndPostsInfo
		json.Unmarshal(gotBytes, &got)
		want := UserAndPostsInfo{
			42,
			UserInfo{"Leanne Graham", "Bret", "Sincere@april.biz"},
			make([]string, 0),
		}

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
