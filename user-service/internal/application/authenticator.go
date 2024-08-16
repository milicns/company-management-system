package application

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/milicns/company-manager/user-service/internal/persistance"
	"github.com/milicns/company-manager/user-service/internal/utils"
)

type Authenticator struct {
	store persistance.Store
}

func NewAuthenticator(store persistance.Store) *Authenticator {
	return &Authenticator{
		store: store,
	}
}

func (auth *Authenticator) Login(loginDto utils.LoginDto) (string, error) {
	err := auth.validateLogin(loginDto)
	var token string
	if err == nil {
		token, _ = generateJwt(loginDto.Username)
		return token, err
	}

	return token, err
}

func generateJwt(username string) (string, error) {

	var secretKey = []byte(os.Getenv("SECRET_KEY"))

	claims := utils.Claims{
		CustomClaims: map[string]string{
			"username": username,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	}

	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenString.SignedString(secretKey)

	return token, err
}

func (auth *Authenticator) validateLogin(loginDto utils.LoginDto) error {
	user, err := auth.store.GetByUsername(loginDto.Username)
	if err != nil {
		return errors.New("username doesn't exist")
	}

	hash := user.Password
	err = utils.Matches(hash, loginDto.Password)
	if err != nil {
		return err
	}
	return nil
}
