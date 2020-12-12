package sqlstore

import (
	"database/sql"
	"github.com/SKQR01/goblog/internal/app/store"
	_ "github.com/lib/pq"
)

//Store ...
type Store struct {
	db *sql.DB
	userRepository *UserRepository
}

//New ...
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

//User ...
func (store *Store) User() store.UserRepository {
	if store.userRepository != nil {
		return store.userRepository
	}
	store.userRepository = &UserRepository{
		store: store,
	}
	return store.userRepository
}
