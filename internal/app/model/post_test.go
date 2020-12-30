package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SKQR01/goblog/internal/app/model"
)

func TestPost_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		post    func() *model.Post
		isValid bool
	}{
		{
			name: "valid",
			post: func() *model.Post {
				return model.TestPost(t)
			},
			isValid: true,
		},
		{
			name: "empty content",
			post: func() *model.Post {
				post := model.TestPost(t)
				post.Content = ""
				return post
			},
			isValid: false,
		},
		{
			name: "empty title",
			post: func() *model.Post {
				post := model.TestPost(t)
				post.Title = ""
				return post
			},
			isValid: false,
		},
		{
			name: "empty title and content",
			post: func() *model.Post {
				post := model.TestPost(t)
				post.Title = ""
				post.Content = ""
				return post
			},
			isValid: false,
		},
		{
			name: "content not json",
			post: func() *model.Post {
				post := model.TestPost(t)
				post.Content = "sdasdass"
				return post
			},
			isValid: false,
		},
		{
			name: "invalid json content",
			post: func() *model.Post {
				post := model.TestPost(t)
				post.Content = "{123}"
				return post
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.post().Validate())
			} else {
				assert.Error(t, tc.post().Validate())
			}
		})
	}
}
