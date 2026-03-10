package user

import (
	"errors"
	"fmt"
	"lab_fullcyle-auction_go/cmd/internal/entity/user_entity"
	"lab_fullcyle-auction_go/cmd/internal/internal_error"
	"lab_fullcyle-auction_go/configuration/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserEntityMongo struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) *UserRepository {
	return &UserRepository{
		Collection: database.Collection("users"),
	}
}
func (r *UserRepository) FindUserByID(id string) (*user_entity.User, *internal_error.InternalError) {
	var user UserEntityMongo
	filter := bson.M{"_id": id}
	err := r.Collection.FindOne(nil, filter).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error(fmt.Sprintf("User not found with this ID = %s", id), err)
			return nil, internal_error.NewNotFoundError("user not found with id: " + id)
		}
		logger.Error(fmt.Sprintf("Error while finding user with ID = %s", id), err)
		return nil, internal_error.NewInternalServerError("failed to find user")
	}

	user_entity := &user_entity.User{
		ID:   user.ID,
		Name: user.Name,
	}
	return user_entity, nil
}
