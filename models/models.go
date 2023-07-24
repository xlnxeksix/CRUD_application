package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Email    string
}

type Product struct {
	gorm.Model
	ID       uint `gorm:"primaryKey"`
	Type     string
	Quantity uint
	Name     string
}

type UserProduct struct {
	gorm.Model
	UserID    uint
	AddressID uint
}
