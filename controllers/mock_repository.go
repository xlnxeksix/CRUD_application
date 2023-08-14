package controllers

import (
	"awesomeProject1/models"
)

type MockUserRepository struct {
	CreateUserFn  func(user *models.User) error
	GetUserByIDFn func(userID uint) (*models.User, error)
	DeleteUserFn  func(userID uint) error
	UpdateUserFn  func(user *models.User, existingUID uint) error
	GetAllUsersFn func() ([]models.User, error)
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	return m.CreateUserFn(user)
}

func (m *MockUserRepository) GetUserByID(userID uint) (*models.User, error) {
	return m.GetUserByIDFn(userID)
}

func (m *MockUserRepository) DeleteUser(userID uint) error {
	return m.DeleteUserFn(userID)
}

func (m *MockUserRepository) UpdateUser(user *models.User, existingUID uint) error {
	return m.UpdateUserFn(user, existingUID)
}

func (m *MockUserRepository) GetAllUsers() ([]models.User, error) {
	return m.GetAllUsersFn()
}

type MockProductRepository struct {
	CreateProductFn  func(product *models.Product) error
	GetProductByIDFn func(productID uint) (*models.Product, error)
	DeleteProductFn  func(productID uint) error
	UpdateProductFn  func(product *models.Product, existingPID uint) error
	GetAllProductsFn func() ([]models.Product, error)
}

func (m *MockProductRepository) CreateProduct(product *models.Product) error {
	return m.CreateProductFn(product)
}

func (m *MockProductRepository) GetProductByID(productID uint) (*models.Product, error) {
	return m.GetProductByIDFn(productID)
}

func (m *MockProductRepository) DeleteProduct(productID uint) error {
	return m.DeleteProductFn(productID)
}

func (m *MockProductRepository) UpdateProduct(product *models.Product, existingPID uint) error {
	return m.UpdateProductFn(product, existingPID)
}

func (m *MockProductRepository) GetAllProducts() ([]models.Product, error) {
	return m.GetAllProductsFn()
}
