package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 12
	minFirstNameLen = 3
	minLastNameLen = 3
	minPasswordLen = 7
)

type CreateUserParams struct {
	FirstName string  `json:"firstName"`
	LastName string  `json:"lastName"`
	Email string  `json:"email"`
	Password string  `json:"password"`
}

func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("firstName length should be at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("lastName length should be at least %d characters", minLastNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password length should be at least %d characters", minPasswordLen)
	}
	if !isEmailValid(params.Email) {
		errors["email"] = fmt.Sprintf("email %s is invalid", params.Email)
	}
	return errors
}


func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string  `bson:"firstName" json:"firstName"`
	LastName string  `bson:"lastName" json:"lastName"`
	Email string  `bson:"email" json:"email"`
	Password string  `bson:"password" json:"-"`

}

func NewUserFromParams(params CreateUserParams) (*User, error){
	encPw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}

	return &User {
		FirstName: params.FirstName,
		LastName: params.LastName,
		Email: params.Email,
		Password: string(encPw),
	}, nil
}