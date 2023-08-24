package pricing

type FurnitureStrategy struct{}

func (s *FurnitureStrategy) CalculatePrice(quantity uint) float64 {
	packagingPrice := 50.0
	return float64(quantity)*1.1 + packagingPrice
}
