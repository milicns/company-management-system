package application

import (
	"github.com/milicns/company-manager/user-service/internal/model"
	"github.com/milicns/company-manager/user-service/internal/persistance"
	"github.com/milicns/company-manager/user-service/internal/utils"
)

type UserService struct {
	store persistance.Store
}

func NewService(store persistance.Store) *UserService {
	return &UserService{
		store: store,
	}
}

func (service *UserService) Create(userDto utils.RegisterDto) error {

	hash, err := utils.Hash(userDto.PlaintextPassword)
	if err != nil {
		return err
	}

	user := model.User{
		Email:    userDto.Email,
		Username: userDto.Username,
		Password: hash,
	}

	return service.store.Create(user)
}

func (service *UserService) GetByUsername(username string) (*model.User, error) {
	return service.store.GetByUsername(username)
}
