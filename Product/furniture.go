package product

type FurnitureStrategy struct{}

func (s *FurnitureStrategy) CalculatePrice(product *Product) float64 {
	packagingPrice := 50.0
	return float64(product.Quantity)*1.1 + packagingPrice
}
