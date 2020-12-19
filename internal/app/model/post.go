package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

//Post post model scheme.
type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	OwnerID int    `json:"owner"`
}

//Validate check for correctness of incoming user data.
func (post *Post) Validate() error {
	return validation.ValidateStruct(
		post,
		validation.Field(&post.Title, validation.Required),
		validation.Field(&post.OwnerID, validation.Required),
		validation.Field(&post.Content, validation.Required, is.JSON),
	)
}
