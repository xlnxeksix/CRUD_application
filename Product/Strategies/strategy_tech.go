package pricing

type TechStrategy struct{}

func (s *TechStrategy) CalculatePrice(quantity uint) float64 {
	packagingPrice := 20.0
	return float64(quantity)*0.95 + packagingPrice
}
