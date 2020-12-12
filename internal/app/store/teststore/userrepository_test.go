package teststore_test

import (
	"github.com/SKQR01/goblog/internal/app/model"
	"github.com/SKQR01/goblog/internal/app/store"
	"github.com/SKQR01/goblog/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	st := teststore.New()
	user := model.TestUser(t)

	assert.NoError(t, st.User().Create(user))
	assert.NotNil(t, user)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	st := teststore.New()
	email := "example@example.com"
	_, err := st.User().FindByEmail(email)

	//TODO: найти по этой фигне инфу
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	err = st.User().Create(model.TestUser(t))
	assert.NoError(t, err)

	user, err := st.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

//TestUserRepository_Find ...
func TestUserRepository_Find(t *testing.T) {
	s := teststore.New()
	u1 := model.TestUser(t)
	s.User().Create(u1)
	u2, err := s.User().Find(u1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}
