package pricing

type Pricing interface {
	CalculatePrice(q uint) float64
}
