package db

import (
	_ "cloud.google.com/go/cloudsqlconn"
	_ "cloud.google.com/go/cloudsqlconn/postgres/pgxv5"
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"password-lock/config"
)

func ConnectToDatabase(cfg *config.Config) *gorm.DB {

	if cfg.Environment != config.LOCAL_ENVIRONEMNT {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DbHost, cfg.DbUser, cfg.DbPassword, cfg.DbName)
		db, err := sql.Open(
			"cloudsql-postgres",
			dsn,
		)
		if err != nil {
			log.Println("Could not connect to database: ", err.Error())
		} else {
			log.Println(" Successfully connected to database")
		}

		gormDb, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), nil)
		if err != nil {
			log.Println("Could not connect to database: ", err.Error())
		} else {
			log.Println(" Successfully connected to database")
		}

		return gormDb
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
