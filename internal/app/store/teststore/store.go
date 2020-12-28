package teststore

import (
	"github.com/SKQR01/goblog/internal/app/model"
	"github.com/SKQR01/goblog/internal/app/store"
)

// Store test store without db
type Store struct {
	userRepository *UserRepository
	postRepository *PostRepository
}

// New ...
func New() *Store {
	return &Store{}
}

// User ...
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[int]*model.User),
	}

	return s.userRepository
}

// Post ...
func (s *Store) Post() store.PostRepository {
	if s.postRepository != nil {
		return s.postRepository
	}

	s.postRepository = &PostRepository{
		store: s,
		posts: make(map[int]*model.Post),
	}

	return s.postRepository
}
