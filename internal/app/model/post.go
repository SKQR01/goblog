package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

//Post model scheme.
type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Owner   *User  `json:"owner"`
}

//NewPost ...
func NewPost() *Post {
	return &Post{
		Owner: &User{},
	}
}

//GetOwnerID gets owner id.
func (post *Post) GetOwnerID() int {
	return post.Owner.ID
}

//SetOwnerID gets owner id.
func (post *Post) SetOwnerID(newID int) {
	post.Owner.ID = newID
}

//Validate ...
func (post *Post) Validate() error {
	//ValidateStruct perform cascade validation of embeded structs if they satisfies Validatable interface (have ValidateFunction).
	err := validation.Validate(&post.Owner.ID, validation.Required.Error("owner:(id: cannot be blank.)."))
	if err != nil {
		return err
	}
	return validation.ValidateStruct(
		post,
		validation.Field(&post.Title, validation.Required),
		validation.Field(&post.Content, validation.Required, is.JSON),
	)
}
