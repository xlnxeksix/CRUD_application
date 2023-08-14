package controllers

import (
	"awesomeProject1/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(userID uint) (*models.User, error)
	DeleteUser(userID uint) error
	UpdateUser(user *models.User, existingUID uint) error
	GetAllUsers() ([]models.User, error)
}

type SQLUserRepository struct {
	DB *gorm.DB
}

func NewSQLUserRepository(db *gorm.DB) UserRepository {
	return &SQLUserRepository{DB: db}
}
func (repo *SQLUserRepository) CreateUser(user *models.User) error {
	insertQuery := "INSERT INTO users (id, username, email, role) VALUES (?, ?, ?, ?)"
	if err := repo.DB.Exec(insertQuery, user.ID, user.Username, user.Email, user.Role).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SQLUserRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := repo.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (repo *SQLUserRepository) DeleteUser(userID uint) error {
	deleteQuery := "DELETE FROM users WHERE id = ?"
	if err := repo.DB.Exec(deleteQuery, userID).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SQLUserRepository) UpdateUser(user *models.User, existingUID uint) error {
	updateQuery := "UPDATE users SET id = ?, username = ?, email = ?, role = ? WHERE id = ?"
	if err := repo.DB.Exec(updateQuery, user.ID, user.Username, user.Email, user.Role, existingUID).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SQLUserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := repo.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

type ProductRepository interface {
	CreateProduct(user *models.Product) error
	GetProductByID(userID uint) (*models.Product, error)
	DeleteProduct(userID uint) error
	UpdateProduct(user *models.Product, existingUID uint) error
	GetAllProducts() ([]models.Product, error)
}

type SQLProductRepository struct {
	DB *gorm.DB
}

func NewSQLProductRepository(db *gorm.DB) ProductRepository {
	return &SQLProductRepository{DB: db}
}

func (repo *SQLProductRepository) CreateProduct(product *models.Product) error {
	insertQuery := "INSERT INTO products (id, name, type, quantity) VALUES (?, ?, ?, ?)"
	if err := repo.DB.Exec(insertQuery, product.ID, product.Name, product.Type, product.Quantity).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SQLProductRepository) GetProductByID(productID uint) (*models.Product, error) {
	var product models.Product
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

func (repo *SQLProductRepository) UpdateProduct(product *models.Product, existingPID uint) error {
	updateQuery := "UPDATE products SET id = ?, name = ?, type = ?, quantity = ? WHERE id = ?"
	if err := repo.DB.Exec(updateQuery, product.ID, product.Name, product.Type, product.Quantity, existingPID).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SQLProductRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	if err := repo.DB.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

type UserProductRepository interface {
	AllocateProduct(userID uint, productID uint) error
	DeallocateProduct(userID uint, productID uint) error
}

type SQLUserProductRepository struct {
	DB          *gorm.DB
	UserRepo    UserRepository
	ProductRepo ProductRepository
}

func NewSQLUserProductRepository(db *gorm.DB) UserProductRepository {
	return &SQLUserProductRepository{DB: db}
}

func (repo *SQLUserProductRepository) AllocateProduct(userID uint, productID uint) error {
	var user *models.User
	var product *models.Product
	user, _ = repo.UserRepo.GetUserByID(userID)
	product, _ = repo.ProductRepo.GetProductByID(productID)
	user.Products = append(user.Products, product)
	product.Users = append(product.Users, user)
	if err := repo.DB.Save(&user).Error; err != nil {
		return err
	}
	if err := repo.DB.Save(&product).Error; err != nil {
		return err
	}
	return nil

}

func (repo *SQLUserProductRepository) DeallocateProduct(userID uint, productID uint) error {
	var user *models.User
	var product *models.Product
	user, _ = repo.UserRepo.GetUserByID(userID)
	product, _ = repo.ProductRepo.GetProductByID(productID)

	// Remove the product from the user's allocated products
	var updatedProducts []*models.Product
	for _, allocatedProduct := range user.Products {
		if allocatedProduct.ID != productID {
			updatedProducts = append(updatedProducts, allocatedProduct)
		}
	}
	user.Products = updatedProducts

	// Remove the user from the product's users
	var updatedUsers []*models.User
	for _, allocatedUser := range product.Users {
		if allocatedUser.ID != userID {
			updatedUsers = append(updatedUsers, allocatedUser)
		}
	}
	product.Users = updatedUsers

	// Save the changes to the database
	if err := repo.DB.Save(&user).Error; err != nil {
		return err
	}
	if err := repo.DB.Save(&product).Error; err != nil {
		return err
	}

	return nil
}
