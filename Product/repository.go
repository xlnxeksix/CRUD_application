package product

import (
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(user *Product) error
	GetProductByID(userID uint) (*Product, error)
	DeleteProduct(userID uint) error
	UpdateProduct(user *Product, existingUID uint) error
	GetAllProducts() ([]Product, error)
}

type SQLProductRepository struct {
	DB *gorm.DB
}

func NewSQLProductRepository(db *gorm.DB) ProductRepository {
	return &SQLProductRepository{DB: db}
}

func (repo *SQLProductRepository) CreateProduct(product *Product) error {
	insertQuery := "INSERT INTO products (name, type, quantity, shipping_price) VALUES (?, ?, ?, ?)"
	if err := repo.DB.Exec(insertQuery, product.Name, product.Type, product.Quantity, product.ShippingPrice).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SQLProductRepository) GetProductByID(productID uint) (*Product, error) {
	var product Product
	if err := repo.DB.First(&product, productID).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (repo *SQLProductRepository) DeleteProduct(productID uint) error {
	deleteQuery := "DELETE FROM products WHERE id = ?"
	if err := repo.DB.Exec(deleteQuery, productID).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SQLProductRepository) UpdateProduct(product *Product, existingPID uint) error {
	updateQuery := "UPDATE products SET name = ?, type = ?, quantity = ? WHERE id = ?"
	if err := repo.DB.Exec(updateQuery, product.Name, product.Type, product.Quantity, existingPID).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SQLProductRepository) GetAllProducts() ([]Product, error) {
	var products []Product
	if err := repo.DB.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
