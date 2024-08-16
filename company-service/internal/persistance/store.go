package persistance

import (
	"github.com/milicns/company-manager/company-service/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Store interface {
	Create(company model.Company) (primitive.ObjectID, error)
	GetOne(id primitive.ObjectID) (*model.Company, error)
	Patch(patchData map[string]interface{}) error
	Delete(id primitive.ObjectID) (*model.Company, error)
}
