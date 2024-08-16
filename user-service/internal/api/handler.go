package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/milicns/company-manager/user-service/internal/application"
	"github.com/milicns/company-manager/user-service/internal/utils"
)

type Handler struct {
	service       *application.UserService
	authenticator *application.Authenticator
}

func NewHandler(service *application.UserService, authenticator *application.Authenticator) *Handler {
	return &Handler{
		service:       service,
		authenticator: authenticator,
	}
}

func (handler *Handler) Create(writer http.ResponseWriter, req *http.Request) {
	var userDto utils.RegisterDto
	err := json.NewDecoder(req.Body).Decode(&userDto)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = handler.service.Create(userDto)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusExpectationFailed)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
}

func (handler *Handler) Login(writer http.ResponseWriter, req *http.Request) {
	var loginDto utils.LoginDto
	err := json.NewDecoder(req.Body).Decode(&loginDto)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var token string
	token, err = handler.authenticator.Login(loginDto)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if encodeErr := json.NewEncoder(writer).Encode(token); encodeErr != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("error encoding token to JSON: %v", encodeErr)
		return
	}
}

func (handler *Handler) Authorize(writer http.ResponseWriter, req *http.Request) {

	authorizationHeader := req.Header.Get("Authorization")

	if authorizationHeader == "" {
		http.Error(writer, "You are not logged in", http.StatusUnauthorized)
		return
	}
	token, username := parseJwt(authorizationHeader)
	_, err := handler.service.GetByUsername(username)
	if !token.Valid || err != nil {
		http.Error(writer, "Token is not valid", http.StatusUnauthorized)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}

func parseJwt(authorizationHeader string) (*jwt.Token, string) {
	tokenString := strings.TrimSpace(strings.Split(authorizationHeader, "Bearer")[1])
	token, _ := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	}, jwt.WithLeeway(5*time.Second))

	claims, _ := token.Claims.(*utils.Claims)
	username := claims.CustomClaims["username"]
	return token, username
}
