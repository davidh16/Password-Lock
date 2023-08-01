package repository

import (
	"gorm.io/gorm"
)

type EntityRepository interface {
	Db() *gorm.DB
}

type entityRepository struct {
	db *gorm.DB
}

func NewEntityRepository(db *gorm.DB) *entityRepository {
	return &entityRepository{
		db: db,
	}
}

func (r entityRepository) Db() *gorm.DB {
	return r.db
}
