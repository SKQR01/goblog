package sqlstore

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/SKQR01/goblog/internal/app/store"
)

//Store for direct interactions with database.
type Store struct {
	db             *sql.DB
	userRepository *UserRepository
	postRepository *PostRepository
}

//New creates store instance.
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

//User adresses to user repository to further actions.
func (store *Store) User() store.UserRepository {
	if store.userRepository != nil {
		return store.userRepository
	}
	store.userRepository = &UserRepository{
		store: store,
	}
	return store.userRepository
}

// Post adresses to post repository to further actions.
func (store *Store) Post() store.PostRepository {
	if store.postRepository != nil {
		return store.postRepository
	}
	store.postRepository = &PostRepository{
		store: store,
	}
	return store.postRepository
}
