package model

//Some assist in testing.
import "testing"

//TestUser creates new test user instance.
func TestUser(t *testing.T) *User {
	return &User{
		Email:    "example@example.com",
		Password: "123456",
	}
}

//TestPost creates new test post instance.
func TestPost(t *testing.T) *Post {
	return &Post{
		Title:   "someTitle",
		Content: `{"data":123}`,
		OwnerID: 1,
	}
}
