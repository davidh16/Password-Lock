package repository

import (
	"fmt"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) Testiram() {
	fmt.Println("aaaaa")
}
