package utils

type RegisterDto struct {
	Email             string `json:"email"`
	Username          string `json:"username"`
	PlaintextPassword string `json:"password"`
}

type LoginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
