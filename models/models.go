package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Email    string
	Role     string
}

type Product struct {
	ID       uint `gorm:"primaryKey"`
	Type     string
	Quantity uint
	Name     string
}

type UserProduct struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	ProductID uint
}
