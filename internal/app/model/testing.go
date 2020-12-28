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
	post := NewPost()

	post.ID = 0
	post.Title = "SomePost"
	post.Content = "{\"data\":123}"

	testUser := TestUser(t)
	post.SetOwnerID(1)
	post.Owner.Email = testUser.Email
	post.Owner.Password = testUser.Password

	return post
}
