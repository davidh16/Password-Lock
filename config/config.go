package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tkanos/gonfig"
	"log"
	"os"
)

// Config struct can be expanded if more env variables are introduced

type Config struct {
	PgUrl                 string `env:"PG_URL"`
	Port                  string `env:"PORT"`
	UserSecretKey         string `env:"USER_SECRET_KEY"`
	EntitySecretKey       string `env:"ENTITY_SECRET_KEY"`
	EntitySecretVector    string `env:"ENTITY_SECRET_VECTOR"`
	ResponseSecretKey     string `env:"RESPONSE_SECRET_KEY"`
	ResponseSecretVector  string `env:"RESPONSE_SECRET_VECTOR"`
	StorageBucket         string `env:"STORAGE_BUCKET"`
	SmtpPort              string `env:"SMTP_PORT"`
	SmtpHost              string `env:"SMTP_HOST"`
	SmtpFrom              string `env:"SMTP_FROM"`
	FirebaseAppPassword   string `env:"FIREBASE_APP_PASSWORD"`
	BaseUrl               string `env:"BASE_URL"`
	DefaultEntityIconPath string `env:"DEFAULT_ENTITY_ICON_PATH"`
}

func GetConfig() *Config {

	configuration := &Config{}

	// fileName could be changed dynamically if there are more env files (like production env and development env ), but for the purpose of this app, it will be hardcoded
	//fileName := ".env"
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load .env file")
		os.Exit(1)
	}

	err := gonfig.GetConf("", configuration)
	if err != nil {
		log.Panic(err.Error())
	}

	return configuration
}
