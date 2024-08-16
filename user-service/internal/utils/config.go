package utils

import (
	"os"
)

type AppConfig struct {
	AppPort string
}

type DbConfig struct {
	DbHost         string
	DbPort         string
	DbRootUsername string
	DbRootPassword string
}

func LoadAppConfig() *AppConfig {
	return &AppConfig{
		AppPort: os.Getenv("APP_PORT"),
	}
}

func LoadDbConfig() *DbConfig {
	return &DbConfig{
		DbHost:         os.Getenv("DB_HOST"),
		DbPort:         os.Getenv("DB_PORT"),
		DbRootUsername: os.Getenv("DB_USERNAME"),
		DbRootPassword: os.Getenv("DB_ROOT_PASSWORD"),
	}
}
