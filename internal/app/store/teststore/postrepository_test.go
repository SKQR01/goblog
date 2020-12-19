package teststore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SKQR01/goblog/internal/app/model"
	"github.com/SKQR01/goblog/internal/app/store/teststore"
)

func TestPostRepository_Create(t *testing.T) {
	s := teststore.New()
	p := model.TestPost(t)
	assert.NoError(t, s.Post().Create(p))
	assert.NotNil(t, p.ID)
}

func TestPostRepository_Find(t *testing.T) {
	s := teststore.New()
	p1 := model.TestPost(t)
	s.Post().Create(p1)
	p2, err := s.Post().Find(p1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, p2)
}

func TestPostRepository_Remove(t *testing.T) {
	s := teststore.New()
	posts := []*model.Post{}
	removeIds := []int{}

	for i := 0; i < 11; i++ {
		post := model.TestPost(t)
		posts = append(posts, post)
		s.Post().Create(post)
		removeIds = append(removeIds, post.ID)
	}

	assert.NoError(t, s.Post().Remove(removeIds))
}
