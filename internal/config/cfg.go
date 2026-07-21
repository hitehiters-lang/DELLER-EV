package config

import (
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseHost   string
	DatabasePort   string
	DatabaseUser   string
	DatabasePass   string
	DatabaseName   string
	ServerPort     string
	JwtSecret      string
	LibreSecret    string
	GmailMailTo    string
	GmailMail      string
	GmailPass      string
	GigaChatAPIKey string
	GigaChatModel  string
	GigaChatAPIURL string
	GigaChatScope  string
	AIAuthorID     string
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
	cfg.LibreSecret = os.Getenv("LIBRE_SECRET")
	cfg.GmailMailTo = os.Getenv("GMAIL_MAIL_TO")
	cfg.GmailMail = os.Getenv("GMAIL_MAIL")
	cfg.GmailPass = os.Getenv("GMAIL_PASS")
	cfg.GigaChatAPIKey = os.Getenv("GIGACHAT_AUTH_KEY")
	cfg.GigaChatModel = os.Getenv("GIGACHAT_MODEL")
	cfg.GigaChatAPIURL = os.Getenv("GIGACHAT_API_URL")
	cfg.GigaChatScope = os.Getenv("GIGACHAT_SCOPE")
	cfg.AIAuthorID = os.Getenv("AI_AUTHOR_ID")

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
	if cfg.GigaChatModel == "" {
		cfg.GigaChatModel = "GigaChat-2-Pro"
	}
	if cfg.GigaChatAPIURL == "" {
		cfg.GigaChatAPIURL = "https://gigachat.devices.sberbank.ru/api/v1"
	}
	if cfg.GigaChatScope == "" {
		cfg.GigaChatScope = "GIGACHAT_API_PERS"
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
	return cfg.JwtSecret
}

func GetLibre() string {
	return cfg.LibreSecret
}

func GetMail() (string, string, string) {
	return cfg.GmailMail, cfg.GmailPass, cfg.GmailMailTo
}

func GetGigaChatAPIKey() string {
	return cfg.GigaChatAPIKey
}

func GetGigaChatModel() string {
	return cfg.GigaChatModel
}

func GetGigaChatAPIURL() string {
	return cfg.GigaChatAPIURL
}

func GetGigaChatScope() string {
	return cfg.GigaChatScope
}

func GetAIAuthorID() uuid.UUID {
	id, _ := uuid.Parse(cfg.AIAuthorID)
	return id
}
