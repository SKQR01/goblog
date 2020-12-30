package teststore

import (
	"github.com/SKQR01/goblog/internal/app/model"
	"github.com/SKQR01/goblog/internal/app/store"
)

// PostRepository ...
type PostRepository struct {
	store *Store
	posts map[int]*model.Post //database table imitation
}

// Create ...
func (r *PostRepository) Create(p *model.Post) error {
	if err := p.Validate(); err != nil {
		return err
	}

	p.ID = len(r.posts) + 1
	r.posts[p.ID] = p

	return nil
}

// Find ...
func (r *PostRepository) Find(id int) (*model.Post, error) {
	p, ok := r.posts[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return p, nil
}

//TODO: finish it for test

// Remove ...
func (r *PostRepository) Remove(ids []int, owner int) error {
	for _, id := range ids {
		if r.posts[id] != nil && r.posts[id].GetOwnerID() == owner {
			delete(r.posts, id)
		}
	}

	return nil
}

//GetRecords ...
func (r *PostRepository) GetRecords(pageNum int, paginationSize int) ([]*model.Post, error) {
	return []*model.Post{}, nil
}
