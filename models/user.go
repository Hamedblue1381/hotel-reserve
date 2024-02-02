package models

import (
	"os"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserParams struct {
	FirstName string `bson:"firstName" json:"firstName" validate:"required"`
	LastName  string `bson:"lastName" json:"lastName" validate:"required"`
	Email     string `bson:"email" json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}
type UpdateUserParams struct {
	FirstName string `bson:"firstName,omitempty" json:"firstName" validate:""`
	LastName  string `bson:"lastName,omitempty" json:"lastName" validate:""`
	Email     string `bson:"email,omitempty" json:"email" validate:"email"`
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"EncryptedPassword" json:"-"`
}

func IsValidPassword(encpw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	cost, _ := strconv.Atoi(os.Getenv("BCRYPT_COST"))
	encryptedPw, err := bcrypt.GenerateFromPassword([]byte(params.Password), cost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encryptedPw),
	}, nil
}

func ToDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}
