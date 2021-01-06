package sqlstore

import (
	"database/sql"
	"fmt"

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
		`WITH i AS (INSERT INTO posts (title, content, owner) VALUES ($1, $2, $3) RETURNING posts.id, posts.owner) 
		SELECT i.id, u.email FROM i LEFT JOIN users u ON i.owner = u.id;`,
		post.Title,
		post.Content,
		post.GetOwnerID(),
	).Scan(&post.ID, &post.Owner.Email)
}

//Find finds user by id.
func (rep *PostRepository) Find(id int) (*model.Post, error) {
	post := model.NewPost()
	if err := rep.store.db.QueryRow(
		"SELECT p.id, p.title, p.content, u.id, u.email FROM posts p LEFT JOIN users u ON p.owner = u.id WHERE p.id = $1",
		id,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.Owner.ID,
		&post.Owner.Email,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return post, nil
}

//Remove removes post by ids.
func (rep *PostRepository) Remove(ids []int, owner int) error {
	if _, err := rep.store.db.Exec("DELETE FROM posts WHERE id = ANY($1) AND owner = $2;", pq.Array(ids), owner); err != nil {
		return err
	}

	return nil
}

//GetRecords ...
func (rep *PostRepository) GetRecords(pageNum int, paginationSize int, userID int) ([]*model.Post, error) {
	byUserID := ""
	if userID >= 0 {
		byUserID = fmt.Sprintf("WHERE owner=%d", userID)
	}

	posts := []*model.Post{}
	singlePost := model.NewPost()

	query := fmt.Sprintf("SELECT p.id, p.title, u.id, u.email FROM posts p LEFT JOIN users u ON p.owner = u.id  %s ORDER BY p.id DESC OFFSET $1 LIMIT $2;", byUserID)

	sqlPosts, err := rep.store.db.Query(
		// "SELECT p.id, p.title, u.id, u.email FROM posts p LEFT JOIN users u ON p.owner = u.id ORDER BY p.id DESC OFFSET $1 LIMIT $2;",
		query,
		(pageNum-1)*paginationSize,
		paginationSize,
	)
	if err != nil {
		return nil, err
	}

	for sqlPosts.Next() != false {
		if err := sqlPosts.Scan(&singlePost.ID, &singlePost.Title, &singlePost.Owner.ID, &singlePost.Owner.Email); err != nil {
			return nil, err
		}
		posts = append(posts, singlePost)
		singlePost = model.NewPost()
	}
	sqlPosts.Close()
	return posts, nil
}
