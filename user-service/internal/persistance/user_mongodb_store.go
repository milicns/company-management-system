package persistance

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/milicns/company-manager/user-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongoDbStore struct {
	users *mongo.Collection
}

func NewUserStore(users *mongo.Collection) Store {

	return &UserMongoDbStore{
		users: users,
	}
}

func (store *UserMongoDbStore) Create(user model.User) error {
	result, err := store.users.InsertOne(context.Background(), user)

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			if e, ok := err.(mongo.WriteException); ok {
				for _, writeErr := range e.WriteErrors {
					if writeErr.Code == 11000 {
						msg := writeErr.Message
						if contains(msg, "username") {
							return errors.New("user with that username already exists")
						} else if contains(msg, "email") {
							return errors.New("user with that email already exists")
						}
					}
				}
			}
		}
		return err
	}

	log.Printf("Registered user with username %s with ID %s", user.Username, result.InsertedID)
	return nil
}

func (store *UserMongoDbStore) GetByUsername(username string) (*model.User, error) {
	filter := bson.D{{Key: "username", Value: username}}
	var result model.User
	err := store.users.FindOne(context.Background(), filter).Decode(&result)

	return &result, err
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
