package db

import (
	"fmt"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"password-lock/config"
)

func ConnectToDatabase(cfg *config.Config) *gorm.DB {

	if cfg.Environment != config.LOCAL_ENVIRONEMNT {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName)
		db, err := gorm.Open(postgres.New(postgres.Config{
			DriverName: "cloudsql-postgres",
			DSN:        dsn,
		}))
		if err != nil {
			log.Println("Could not connect to database: ", err.Error())
		} else {
			log.Println(" Successfully connected to database")
		}

		return db
	}

	dsn := fmt.Sprintf("postgres://%s:%s@postgres:5432/%s?sslmode=disable", cfg.DbUser, cfg.DbPassword, cfg.DbName)
	db, err := gorm.Open(postgres.Open(dsn), nil)
	if err != nil {
		log.Println("Could not connect to database: ", err.Error())
	} else {
		log.Println(" Successfully connected to database")
	}

	return db
}

const (
	USERS_TABLE                   = "users"
	ENTITIES_TABLE                = "entities"
	TOKENS_TABLE                  = "tokens"
	PERSONAL_QUESTIONS_TABLE      = "personal_questions"
	USER_PERSONAL_QUESTIONS_TABLE = "user_personal_questions"
)
