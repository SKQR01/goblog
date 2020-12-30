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
		`WITH i AS (INSERT INTO posts (title, content, owner) VALUES ($1, $2, $3) RETURNING posts.id, posts.owner) 
		SELECT i.id, u.email FROM i LEFT JOIN users u ON i.owner = u.id;`,
		post.Title,
		post.Content,
		post.GetOwnerID(),
	).Scan(&post.ID, &post.Owner.Email)
}

//Find finds user by id.
func (rep *PostRepository) Find(id int) (*model.Post, error) {
	post := &model.Post{}
	if err := rep.store.db.QueryRow(
		"SELECT id, title, content FROM posts WHERE id = $1",
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
func (rep *PostRepository) Remove(ids []int, owner int) error {
	if _, err := rep.store.db.Exec("DELETE FROM posts WHERE id = ANY($1) AND owner = $2;", pq.Array(ids), owner); err != nil {
		return err
	}

	return nil
}

//GetRecords ...
func (rep *PostRepository) GetRecords(pageNum int, paginationSize int) ([]*model.Post, error) {

	posts := []*model.Post{}
	singlePost := model.NewPost()

	sqlPosts, err := rep.store.db.Query(
		// "SELECT id, title, content, owner FROM posts ORDER BY id ASC OFFSET $1 LIMIT $2;",

		"SELECT p.id, p.title, p.content, u.id, u.email FROM posts p LEFT JOIN users u ON p.owner = u.id ORDER BY p.id DESC OFFSET $1 LIMIT $2;",
		pageNum*paginationSize,
		paginationSize,
	)

	if err != nil {
		return nil, err
	}

	for sqlPosts.Next() {

		if err := sqlPosts.Scan(&singlePost.ID, &singlePost.Title, &singlePost.Content, &singlePost.Owner.ID, &singlePost.Owner.Email); err != nil {
			return nil, err
		}
		posts = append(posts, singlePost)

		singlePost = model.NewPost()
	}

	return posts, nil
}
