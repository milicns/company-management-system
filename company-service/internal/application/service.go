package application

import (
	"github.com/milicns/company-manager/company-service/internal/model"
	"github.com/milicns/company-manager/company-service/internal/persistance"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CompanyService struct {
	store persistance.Store
}

func NewService(store persistance.Store) *CompanyService {
	return &CompanyService{
		store: store,
	}
}

func (service *CompanyService) Create(company model.Company) (primitive.ObjectID, error) {

	if err := validateCreate(company); err != nil {
		return primitive.ObjectID{}, *err
	}

	return service.store.Create(company)
}

func (service *CompanyService) Delete(id primitive.ObjectID) (*model.Company, error) {
	return service.store.Delete(id)
}

func (service *CompanyService) GetOne(id primitive.ObjectID) (*model.Company, error) {
	return service.store.GetOne(id)
}

func (service *CompanyService) Patch(patchData map[string]interface{}) error {
	if err := validatePatchData(patchData); err != nil {
		return *err
	}

	return service.store.Patch(patchData)
}
