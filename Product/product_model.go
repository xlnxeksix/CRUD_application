package product

type Product struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Type     string
	Quantity uint
}
