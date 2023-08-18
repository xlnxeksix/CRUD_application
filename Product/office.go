package product

type OfficeStrategy struct{}

func (s *OfficeStrategy) CalculatePrice(product *Product) float64 {
	packagingPrice := 5.0
	return float64(product.Quantity)*0.75 + packagingPrice
}
