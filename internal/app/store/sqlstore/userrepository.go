package sqlstore

import (
	"database/sql"

	"github.com/SKQR01/goblog/internal/app/model"
	"github.com/SKQR01/goblog/internal/app/store"
)

//UserRepository user repository struct accept instance of store.
type UserRepository struct {
	store *Store
}

//Create creates user.
func (rep *UserRepository) Create(user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	if err := user.BeforeCreate(); err != nil {
		return err
	}

	return rep.store.db.QueryRow(
		"INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING id",
		user.Email,
		user.EncryptedPassword,
	).Scan(&user.ID)
}

//Find finds user by id.
func (rep *UserRepository) Find(id int) (*model.User, error) {
	user := &model.User{}
	if err := rep.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM users WHERE id = $1",
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return user, nil
}

//FindByEmail finds user by email.
func (rep *UserRepository) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}
	if err := rep.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM users WHERE email = $1",
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return user, nil
}
