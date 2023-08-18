package product

type TechStrategy struct{}

func (s *TechStrategy) CalculatePrice(product *Product) float64 {
	packagingPrice := 20.0
	return float64(product.Quantity)*0.95 + packagingPrice
}
