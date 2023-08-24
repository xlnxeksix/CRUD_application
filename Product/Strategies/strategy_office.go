package pricing

type OfficeStrategy struct{}

func (s *OfficeStrategy) CalculatePrice(quantity uint) float64 {
	packagingPrice := 5.0
	return float64(quantity)*0.75 + packagingPrice
}
