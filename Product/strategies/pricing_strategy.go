package product

import product "awesomeProject1/Product"

type Pricing interface {
	CalculatePrice(p *product.Product) float64
}
