package persistance

import (
	"context"
	"fmt"
	"log"

	"github.com/milicns/company-manager/company-service/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetConn(config utils.DbConfig) (*mongo.Client, func()) {
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
	companies := client.Database("company").Collection("companies")
	uniqueNameIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := companies.Indexes().CreateOne(context.Background(), uniqueNameIndex)
	if err != nil {
		fmt.Println("failed to create unique field: %w", err)
	}
	return companies
}
