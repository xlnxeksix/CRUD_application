package product

import user "awesomeProject1/User"

type Product struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Type     string
	Quantity uint
	Users    []*user.User `gorm:"many2many:user_products;constraint:OnDelete:CASCADE"`
}
