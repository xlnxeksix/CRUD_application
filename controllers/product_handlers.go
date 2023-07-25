package controllers

import (
	"awesomeProject1/models"
	"fmt"
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
	insertQuery := fmt.Sprintf("INSERT INTO products (id, name, type, quantity) VALUES (%d, '%s', '%s', %d)",
		product.ID, product.Name, product.Type, product.Quantity)

	ctrl.DB.Raw(insertQuery).Row()
	/*
		if err := row.Scan(&product.ID); err != nil {
			models.Logger.Error("Error creating product", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
			return
		}*/

	models.Logger.Info("Product created successfully")
	c.JSON(http.StatusCreated, product)
}

func (ctrl *ProductController) DeleteProductHandler(c *gin.Context) {
	idParam := c.Param("id")
	productID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		models.Logger.Error("Invalid ID provided", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID provided"})
		return
	}

	deleteQuery := fmt.Sprintf("DELETE FROM products WHERE id = '%d'", productID)
	ctrl.DB.Exec(deleteQuery)
	/*
		if err := ; err != nil {
			models.Logger.Error("Error deleting product")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
			return
		}
	*/
	models.Logger.Info("Product deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "product is deleted"})
}

// UpdateProductHandler handles updating a specific product
func (ctrl *ProductController) UpdateProductHandler(c *gin.Context) {
	var product models.Product
	idParam := c.Param("id")
	productID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		models.Logger.Error("Invalid ID provided", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID provided"})
		return
	}

	query := fmt.Sprintf("SELECT * FROM products WHERE id = '%d' FOR UPDATE", productID)

	row := ctrl.DB.Raw(query).Row()
	row.Scan(&product.ID, &product.Name, &product.Type, &product.Quantity)
	/*
		if err := r; err != nil {
			models.Logger.Error("Error finding product", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
	*/
	if err := c.ShouldBindJSON(&product); err != nil {
		models.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateQuery := fmt.Sprintf("UPDATE products SET id = %d, name = '%s', type = '%s', quantity = %d WHERE id = '%d'",
		product.ID, product.Name, product.Type, product.Quantity, productID)
	ctrl.DB.Exec(updateQuery)
	/*
		if err := ; err != nil {
			models.Logger.Error("Error updating product")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
			return
		}
	*/
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
	idParam := c.Param("id")
	productID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		models.Logger.Error("Invalid ID provided", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID provided"})
		return
	}
	query := fmt.Sprintf("SELECT * FROM products WHERE id = '%d'", productID)
	row := ctrl.DB.Raw(query).Row()
	row.Scan(&product.ID, &product.Type, &product.Quantity, &product.Name)
	/*
		if err :=  err != nil {
			models.Logger.Error("Product not found", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "user-product association not found"})
			return
		}*/
	c.JSON(http.StatusOK, product)
}
