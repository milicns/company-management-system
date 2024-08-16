package persistance

import (
	"context"
	"errors"
	"log"

	"github.com/milicns/company-manager/company-service/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CompanyMongoDbStore struct {
	companies *mongo.Collection
}

func NewCompanyStore(companies *mongo.Collection) Store {
	return &CompanyMongoDbStore{companies: companies}
}

func (store *CompanyMongoDbStore) Create(company model.Company) (primitive.ObjectID, error) {

	result, err := store.companies.InsertOne(context.Background(), company)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return primitive.ObjectID{}, errors.New("company with that name already exists")
		}
		return primitive.ObjectID{}, err
	}
	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.ObjectID{}, errors.New("internal server error")
	}
	log.Printf("Created company named %s with ID %s", company.Name, id)
	return id, nil
}

func (store *CompanyMongoDbStore) GetOne(id primitive.ObjectID) (*model.Company, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	var result model.Company
	err := store.companies.FindOne(context.Background(), filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("company was not found")
		}
		return nil, err
	}

	return &result, nil
}

func (store *CompanyMongoDbStore) Patch(patchData map[string]interface{}) error {
	filter := bson.M{"_id": patchData["id"]}
	delete(patchData, "id")
	update := bson.M{"$set": patchData}

	result, err := store.companies.UpdateOne(context.Background(), filter, update)

	if result.MatchedCount == 0 {
		return errors.New("company was not found")
	}

	return err
}

func (store *CompanyMongoDbStore) Delete(id primitive.ObjectID) (*model.Company, error) {
	var deleted model.Company

	filter := bson.D{{Key: "_id", Value: id}}
	err := store.companies.FindOneAndDelete(context.Background(), filter).Decode(&deleted)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("company was not found")
		}
		return nil, err
	}
	log.Println("Deleted company named ", deleted.Name)

	return &deleted, nil
}
