package product

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name          string
	Type          string
	Quantity      uint
	ShippingPrice *float64
}
