package utils

import (
	"os"
)

type AppConfig struct {
	AppPort string
}

type KafkaConfig struct {
	KafkaHost string
	KafkaPort string
}

type DbConfig struct {
	DbHost         string
	DbPort         string
	DbRootUsername string
	DbRootPassword string
}

type UserServiceConfig struct {
	UserServiceHost string
	UserServicePort string
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		AppPort: os.Getenv("APP_PORT"),
	}
}

func LoadKafkaConfig() KafkaConfig {
	return KafkaConfig{
		KafkaHost: os.Getenv("KAFKA_HOST"),
		KafkaPort: os.Getenv("KAFKA_PORT"),
	}
}

func LoadDbConfig() DbConfig {
	return DbConfig{
		DbHost:         os.Getenv("DB_HOST"),
		DbPort:         os.Getenv("DB_PORT"),
		DbRootUsername: os.Getenv("DB_USERNAME"),
		DbRootPassword: os.Getenv("DB_ROOT_PASSWORD"),
	}
}

func LoadUserServiceConfig() UserServiceConfig {
	return UserServiceConfig{
		UserServiceHost: os.Getenv("USER_SERVICE_HOST"),
		UserServicePort: os.Getenv("USER_SERVICE_PORT"),
	}
}
