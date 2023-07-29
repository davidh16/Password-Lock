package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"password-lock/config"
)

func ConnectToDatabase() *gorm.DB {
	cfg := config.GetConfig()
	db, err := gorm.Open(postgres.Open(cfg.PgUrl), nil)
	if err != nil {
		fmt.Println("Could not connect to database: ", err.Error())
		return nil
	}

	fmt.Println(" Successfully connected to database")
	return db
}
