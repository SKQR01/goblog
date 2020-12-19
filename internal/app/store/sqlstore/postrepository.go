package sqlstore

import (
	"database/sql"

	"github.com/lib/pq"

	"github.com/SKQR01/goblog/internal/app/model"
	"github.com/SKQR01/goblog/internal/app/store"
)

//PostRepository post repository struct accept instance of store.
type PostRepository struct {
	store *Store
}

//Create creates post.
func (rep *PostRepository) Create(post *model.Post) error {
	if err := post.Validate(); err != nil {
		return err
	}

	return rep.store.db.QueryRow(
		"INSERT INTO posts (title, content, owner) VALUES ($1, $2, $3) RETURNING id",
		post.Title,
		post.Content,
		post.OwnerID,
	).Scan(&post.ID)
}

//Find finds user by id.
func (rep *PostRepository) Find(id int) (*model.Post, error) {
	post := &model.Post{}
	if err := rep.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM users WHERE id = $1",
		id,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return post, nil
}

//Remove removes post by ids.
func (rep *PostRepository) Remove(ids []int) error {
	if _, err := rep.store.db.Exec("DELETE FROM posts WHERE id = ANY($1);", pq.Array(ids)); err != nil {
		return err
	}

	return nil
}
