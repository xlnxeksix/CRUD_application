package product

type MockProductRepository struct {
	CreateProductFn  func(product *Product) error
	GetProductByIDFn func(productID uint) (*Product, error)
	DeleteProductFn  func(productID uint) error
	UpdateProductFn  func(product *Product, existingPID uint) error
	GetAllProductsFn func() ([]Product, error)
}

func (m *MockProductRepository) CreateProduct(product *Product) error {
	return m.CreateProductFn(product)
}

func (m *MockProductRepository) GetProductByID(productID uint) (*Product, error) {
	return m.GetProductByIDFn(productID)
}

func (m *MockProductRepository) DeleteProduct(productID uint) error {
	return m.DeleteProductFn(productID)
}

func (m *MockProductRepository) UpdateProduct(product *Product, existingPID uint) error {
	return m.UpdateProductFn(product, existingPID)
}

func (m *MockProductRepository) GetAllProducts() ([]Product, error) {
	return m.GetAllProductsFn()
}
