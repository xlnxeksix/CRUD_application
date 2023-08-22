package product

type Pricing interface {
	CalculatePrice(t *Product) float64
}
