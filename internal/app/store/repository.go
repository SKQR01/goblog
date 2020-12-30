package store

import "github.com/SKQR01/goblog/internal/app/model"

//UserRepository users repository for database interactions.
type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}

//PostRepository post repository for database interactions.
type PostRepository interface {
	Create(*model.Post) error
	Find(int) (*model.Post, error)
	Remove([]int, int) error
	//page and pagination
	GetRecords(pageNum int, paginationSize int) ([]*model.Post, error)
}
