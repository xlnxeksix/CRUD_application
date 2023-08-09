package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Email    string
	Role     string
	Products []*Product `gorm:"many2many:user_products;constraint:OnDelete:CASCADE"`
}

type Product struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Type     string
	Quantity uint
	Users    []*User `gorm:"many2many:user_products;constraint:OnDelete:CASCADE"`
}
