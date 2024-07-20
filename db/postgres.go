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

		dbURI := fmt.Sprintf("user=%s password=%s database=%s host=%s",
			cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbHost)

		// dbPool is the pool of database connections.
		dbPool, err := sql.Open("pgx", dbURI)
		if err != nil {
			log.Fatalf("sql.Open: %w", err)
		}

		db, err := gorm.Open(postgres.New(postgres.Config{
			Conn: dbPool,
		}), &gorm.Config{})
		if err != nil {
			log.Fatalf("Could not connect to database: %s", err.Error())
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
