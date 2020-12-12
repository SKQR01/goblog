package teststore

import (
	"github.com/SKQR01/goblog/internal/app/model"
	"github.com/SKQR01/goblog/internal/app/store"
)

//Store ...
type Store struct {
	userRepository *UserRepository
}

//New ...
func New() *Store {
	return &Store{}
}

//User ...
func (store *Store) User() store.UserRepository {
	if store.userRepository != nil {
		return store.userRepository
	}
	store.userRepository = &UserRepository {
		store: store,
		users: make(map[int]*model.User),
	}
	return store.userRepository
}
