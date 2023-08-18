package product

type Pricing interface {
	CalculatePrice(p *Product) float64
}
