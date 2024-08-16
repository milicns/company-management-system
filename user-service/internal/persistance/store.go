package persistance

import "github.com/milicns/company-manager/user-service/internal/model"

type Store interface {
	Create(user model.User) error
	GetByUsername(username string) (*model.User, error)
}
