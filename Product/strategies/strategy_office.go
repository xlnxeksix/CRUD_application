package product

import product "awesomeProject1/Product"

type OfficeStrategy struct{}

func (s *OfficeStrategy) CalculatePrice(product *product.Product) float64 {
	packagingPrice := 5.0
	return float64(product.Quantity)*0.75 + packagingPrice
}
