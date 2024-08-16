package persistance

import (
	"context"
	"fmt"
	"log"

	"github.com/milicns/company-manager/user-service/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetConn(config *utils.DbConfig) (*mongo.Client, func()) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/", config.DbRootUsername, config.DbRootPassword, config.DbHost, config.DbPort)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	return client, func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("failed to close MongoDB connection: %v", err)
		}
		log.Println("closing MongoDB connection")
	}
}

func GetCollection(client *mongo.Client) *mongo.Collection {
	users := client.Database("company").Collection("users")
	uniqueUsernameIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := users.Indexes().CreateOne(context.Background(), uniqueUsernameIndex)
	if err != nil {
		log.Println("failed to create unique index for username: %w", err)
	}

	uniqueEmailIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = users.Indexes().CreateOne(context.Background(), uniqueEmailIndex)
	if err != nil {
		log.Println("failed to create unique index for email: %w", err)
	}

	return users
}
