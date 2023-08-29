package SIEM

import "gorm.io/gorm"

type SIEMRepository interface {
}

type SQLProductRepository struct {
	DB *gorm.DB
}

func NewSQLProductRepository(db *gorm.DB) SIEMRepository {
	return &SQLProductRepository{DB: db}
}
