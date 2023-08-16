package user

import product "awesomeProject1/Product"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Email    string
	Role     string
	Products []*product.Product `gorm:"many2many:user_products;constraint:OnDelete:CASCADE"`
}
