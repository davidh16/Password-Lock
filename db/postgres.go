package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"password-lock/config"
)

func ConnectToDatabase() *gorm.DB {
	cfg := config.GetConfig()
	db, err := gorm.Open(postgres.Open(cfg.PgUrl), nil)
	if err != nil {
		log.Panic("Could not connect to database: ", err.Error())
	}

	fmt.Println(" Successfully connected to database")
	return db
}

const (
	USERS_TABLE                   = "users"
	ENTITIES_TABLE                = "entities"
	TOKENS_TABLE                  = "tokens"
	PERSONAL_QUESTIONS_TABLE      = "personal_questions"
	USER_PERSONAL_QUESTIONS_TABLE = "user_personal_questions"
)
