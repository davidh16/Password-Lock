package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tkanos/gonfig"
	"os"
)

// Config struct can be expanded if more env variables are introduced

type Config struct {
	PgUrl     string `env:"PG_URL"`
	Port      string `env:"PORT"`
	SecretKey string `env:"SECRET_KEY"`
}

func GetConfig() Config {

	configuration := Config{}

	// fileName could be changed dynamically if there are more env files (like production env and development env ), but for the purpose of this app, it will be hardcoded
	//fileName := ".env"
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load .env file")
		os.Exit(1)
	}

	err := gonfig.GetConf("", &configuration)
	if err != nil {
		fmt.Println(err)
		return Config{}
	}

	return configuration
}
