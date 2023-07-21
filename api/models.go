package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Email    string
}

type Product struct {
	ID       uint `gorm:"primaryKey"`
	Type     string
	Quantity uint
	Name     string
}

type UserProduct struct {
	UserID    uint
	AddressID uint
}
