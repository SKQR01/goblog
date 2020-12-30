package apiserver

//Common testing of server.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"

	"github.com/SKQR01/goblog/internal/app/model"
	"github.com/SKQR01/goblog/internal/app/store/teststore"
)

//TestServerHandlePostsCreate ...
func TestServerHandlePostsCreate(t *testing.T) {
	store := teststore.New()
	srv := newServer(store, sessions.NewCookieStore([]byte("secret")))

	testUser := model.TestUser(t)

	store.User().Create(testUser)

	testCases := []struct {
		name         string
		payload      func() interface{}
		cookieValue  map[interface{}]interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: func() interface{} {
				return model.TestPost(t)
			},
			cookieValue: map[interface{}]interface{}{
				"user_id": testUser.ID,
			},
			expectedCode: http.StatusCreated,
		},
		{
			//TODO:валидатор реагирует на чисто числовые и строковые значения, в том числе "1234", "sdsads", как на JSON
			name: "invalid content",
			payload: func() interface{} {
				post := model.TestPost(t)
				post.Content = "12321s3"
				return post
			},
			cookieValue: map[interface{}]interface{}{
				"user_id": testUser.ID,
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
		{
			name: "invalid payload",
			payload: func() interface{} {
				return "invalid payload"
			},
			cookieValue: map[interface{}]interface{}{
				"user_id": testUser.ID,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "without token",
			payload: func() interface{} {
				post := model.TestPost(t)
				return post
			},
			cookieValue:  nil,
			expectedCode: http.StatusUnauthorized,
		},
	}

	secretKey := []byte("secret")
	sc := securecookie.New(secretKey, nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload())
			rec := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPost, "/private/posts/create", b)
			cookieStr, _ := sc.Encode(sessionName, tc.cookieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServerHandlePostsDelete(t *testing.T) {
	store := teststore.New()
	srv := newServer(store, sessions.NewCookieStore([]byte("secret")))

	type reqStruct struct {
		postsToRemoveIds interface{}
	}

	testUser := model.TestUser(t)

	store.User().Create(testUser)

	testCases := []struct {
		name         string
		payload      func() interface{}
		cookieValue  map[interface{}]interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: func() interface{} {
				return &reqStruct{
					postsToRemoveIds: []int{1, 2, 3},
				}
			},
			cookieValue: map[interface{}]interface{}{
				"user_id": testUser.ID,
			},
			expectedCode: http.StatusOK,
		},
		//Dont know how to test server decodes data without the error only in test case
		//https://www.json.org/json-en.html
		// {
		// 	name: "not an array",
		// 	payload: func() interface{} {
		// 		return &reqStruct{
		// 			postsToRemoveIds: "a1a1234asdas",
		// 		}
		// 	},
		// 	cookieValue: map[interface{}]interface{}{
		// 		"user_id": testUser.ID,
		// 	},
		// 	expectedCode: http.StatusBadRequest,
		// },
		{
			name: "empty array",
			payload: func() interface{} {
				return &reqStruct{
					postsToRemoveIds: []int{},
				}
			},
			cookieValue: map[interface{}]interface{}{
				"user_id": testUser.ID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "without token",
			payload: func() interface{} {
				return &reqStruct{
					postsToRemoveIds: []int{1, 2, 3},
				}
			},
			cookieValue:  nil,
			expectedCode: http.StatusUnauthorized,
		},
	}

	secretKey := []byte("secret")
	sc := securecookie.New(secretKey, nil)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload())
			rec := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodPost, "/private/posts/remove", b)
			cookieStr, _ := sc.Encode(sessionName, tc.cookieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionName, cookieStr))
			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
