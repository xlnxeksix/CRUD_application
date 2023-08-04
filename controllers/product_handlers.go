package controllers

import (
	"awesomeProject1/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type ProductController struct {
	DB *gorm.DB
}

func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{DB: db}
}

func (ctrl *ProductController) CreateProductHandler(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		models.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	insertQuery := "INSERT INTO products (id, name, type, quantity) VALUES (?, ?, ?, ?)"

	ctrl.DB.Raw(insertQuery).Row()

	if err := ctrl.DB.Exec(insertQuery, product.ID, product.Name, product.Type, product.Quantity).Error; err != nil {
		models.Logger.Error("Error creating user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	models.Logger.Info("Product created successfully")
	c.JSON(http.StatusCreated, product)
}

func (ctrl *ProductController) DeleteProductHandler(c *gin.Context) {
	idStr := c.Param("id")
	_, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var count int64
	err = ctrl.DB.Model(&models.User{}).Where("id = ?", c.Param("id")).Count(&count).Error
	if err != nil {
		models.Logger.Error("Error checking if user exists", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to check user existence"})
		return
	}
	if count == 0 {
		// User not found, return an error
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	deleteQuery := "DELETE FROM products WHERE id = ?"
	err = ctrl.DB.Exec(deleteQuery, c.Param("id")).Error
	if err != nil {
		models.Logger.Error("Error deleting user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed ti delete products"})
		return
	}

	models.Logger.Info("product deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Product is deleted"})

}

// UpdateProductHandler handles updating a specific product
func (ctrl *ProductController) UpdateProductHandler(c *gin.Context) {
	var product models.Product

	query := "SELECT * FROM products WHERE id = ?"
	err := ctrl.DB.Exec(query, c.Param("id")).Error

	if err != nil {
		models.Logger.Error("Error updating product", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update user"})
		return
	}

	row := ctrl.DB.Raw(query, c.Param("id")).Row()

	if err := row.Scan(&product.ID, &product.Name, &product.Type, &product.Quantity); err != nil {
		models.Logger.Error("Error finding product", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		models.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateQuery := "UPDATE products SET ID = ?, name = ?, type = ?, quantity = ? WHERE id = ?"
	err = ctrl.DB.Exec(updateQuery, product.ID, product.Name, product.Type, product.Quantity, c.Param("id")).Error

	if err != nil {
		models.Logger.Error("Error updating product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
		return
	}

	models.Logger.Info("Product updated successfully")
	c.JSON(http.StatusOK, product)
}

// GetAllproductsHandler handles getting all products
func (ctrl *ProductController) GetAllProductsHandler(c *gin.Context) {
	var products []models.Product

	query := "SELECT * FROM products"

	rows, err := ctrl.DB.Raw(query).Rows()
	if err != nil {
		models.Logger.Error("Error getting all products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get products"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Type, &product.Quantity, &product.Name); err != nil {
			models.Logger.Error("Error scanning product row", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get products"})
			return
		}
		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}

// GetSpecificproductHandler handles getting a specific product
func (ctrl *ProductController) GetSpecificproductHandler(c *gin.Context) {
	var product models.Product
	query := "SELECT * FROM products WHERE id = ?"
	rows, err := ctrl.DB.Raw(query, c.Param("id")).Rows()
	if err != nil {
		models.Logger.Error("Error finding product", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID must be an integer"})
		return
	}
	defer rows.Close()

	found := false // Variable to keep track if the user was found or not

	for rows.Next() {
		err := rows.Scan(&product.ID, &product.Name, &product.Type, &product.Quantity)
		if err != nil {
			models.Logger.Error("Error scanning product row", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get product"})
			return
		}
		found = true // User found, set the flag to true
		break        // Assuming that there will be only one row since ID is unique
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}
