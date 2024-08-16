package utils

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	CustomClaims map[string]string `json:"custom_claims"`
	jwt.RegisteredClaims
}
