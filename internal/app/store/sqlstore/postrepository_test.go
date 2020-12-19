package sqlstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SKQR01/goblog/internal/app/model"
	"github.com/SKQR01/goblog/internal/app/store/sqlstore"
)

func TestPostRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("posts, users")

	s := sqlstore.New(db)

	u := model.TestUser(t)
	s.User().Create(u)

	p := model.TestPost(t)
	p.OwnerID = u.ID

	assert.NoError(t, s.Post().Create(p))
	assert.NotNil(t, p.ID)
}

func TestPostRepository_Find(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("posts")

	s := sqlstore.New(db)
	u1 := model.TestUser(t)
	s.User().Create(u1)
	u2, err := s.User().Find(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

func TestPostRepository_Remove(t *testing.T)  {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("posts, users")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	s.User().Create(u)

	p := model.TestPost(t)
	ids := []int{}
	for i := 0; i < 20; i++{
		p.OwnerID = u.ID
		s.Post().Create(p)
		ids = append(ids, i)
	}
	assert.NoError(t, s.Post().Remove(ids)) 
}