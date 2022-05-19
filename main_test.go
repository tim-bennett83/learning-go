package main

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestMain(t *testing.T) {
	t.Run("mew", func(t *testing.T) {
		response := httptest.NewRecorder()
		handleUserPosts()(response, httptest.NewRequest("GET", "/v1/user-posts/42", nil))
		defer response.Result().Body.Close()

		got, err := ioutil.ReadAll(response.Result().Body)
		want := `{"userID": "42"}`

		if err != nil {
			t.Errorf("got error %v", err)
		}
		if string(got) != want {
			t.Errorf("expected %v got %v", want, string(got))
		}
	})
}
