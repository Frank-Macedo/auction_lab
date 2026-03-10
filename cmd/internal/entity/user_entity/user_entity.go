package user_entity

import (
	"lab_fullcyle-auction_go/cmd/internal/internal_error"
)

type User struct {
	ID   string `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
}

type UserRepositoryInterface interface {
	FindUserByID(id string) (*User, *internal_error.InternalError)
}
