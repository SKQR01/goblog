package model

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"golang.org/x/crypto/bcrypt"
)

//User ...
type User struct {
	ID                int		`json:"id"`
	Email             string	`json:"email"`
	Password          string	`json:"password, omitempty"`
	EncryptedPassword string	`json:"-"`
}

func (user *User) BeforeCreate() error {
	if len(user.Password) > 0 {
		enc, err := encryptString(user.Password)
		if err != nil {
			return err
		}
		user.EncryptedPassword = enc
	}
	return nil
}

//Sanitaze ...
func (user *User) Sanitaze()  {
	user.Password = ""
}

//Validate ...
func (user *User) Validate() error {
	return validation.ValidateStruct(
		user,
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.By(requiredIf(user.EncryptedPassword == "")), validation.Length(6, 100)),
	)
}

func encryptString(s string) (string, error) {
	bcr, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(bcr), nil
}

func (user *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(password)) == nil
}
