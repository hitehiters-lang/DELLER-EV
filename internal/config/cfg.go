package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseHost string
	DatabasePort string
	DatabaseUser string
	DatabasePass string
	DatabaseName string
	ServerPort   string
	JwtSecret    string
}

var cfg Config

func LoadENV() {
	godotenv.Load(".env")

	cfg.DatabaseHost = os.Getenv("HOST")
	cfg.DatabasePort = os.Getenv("DB_PORT")
	cfg.DatabaseUser = os.Getenv("DB_USER")
	cfg.DatabasePass = os.Getenv("PASSWORD")
	cfg.DatabaseName = os.Getenv("DB_NAME")
	cfg.ServerPort = os.Getenv("SERVER_PORT")
	cfg.JwtSecret = os.Getenv("JWT_SECRET")

	if cfg.DatabaseHost == "" {
		cfg.DatabaseHost = "localhost"
	}
	if cfg.DatabasePort == "" {
		cfg.DatabasePort = "5432"
	}
	if cfg.DatabaseUser == "" {
		cfg.DatabaseUser = "postgres"
	}
	if cfg.DatabasePass == "" {
		cfg.DatabasePass = "123"
	}
	if cfg.DatabaseName == "" {
		cfg.DatabaseName = "test"
	}
	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080"
	}
}

func GetSQL() (string, string, string, string, string) {
	return cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseUser,
		cfg.DatabasePass,
		cfg.DatabaseName
}
func GetServer() string {
	return cfg.ServerPort
}
func GetJwt() string {
	return cfg.ServerPort
}
