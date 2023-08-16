package product

import (
	"awesomeProject1/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type ProductController struct {
	Repo ProductRepository
}

func NewProductController(repo ProductRepository) *ProductController {
	return &ProductController{Repo: repo}
}

func (ctrl *ProductController) CreateProductHandler(c *gin.Context) {
	var product Product

	if err := c.ShouldBindJSON(&product); err != nil {
		models.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the product exists
	if existingProduct, _ := ctrl.Repo.GetProductByID(product.ID); existingProduct != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product with same ID exists"})
		return
	}

	if err := ctrl.Repo.CreateProduct(&product); err != nil {
		models.Logger.Error("Error creating product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	models.Logger.Info("Product created successfully")
	c.JSON(http.StatusCreated, product)
}

func (ctrl *ProductController) DeleteProductHandler(c *gin.Context) {
	productIDStr := c.Param("id")
	productID32, err := strconv.ParseUint(productIDStr, 10, 64)
	productID := uint(productID32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Check if the product exists
	product, err := ctrl.Repo.GetProductByID(productID)
	if err != nil {
		models.Logger.Error("Error retrieving product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Delete the product
	if err := ctrl.Repo.DeleteProduct(productID); err != nil {
		models.Logger.Error("Error deleting product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	models.Logger.Info("Product deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (ctrl *ProductController) UpdateProductHandler(c *gin.Context) {
	var updatedProduct Product

	// Check if the product exists
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Fetch the existing product from the repository
	existingProduct, err := ctrl.Repo.GetProductByID(uint(productID))
	if err != nil {
		models.Logger.Error("Error finding product", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Bind JSON data to the updatedProduct struct
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		models.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update only the necessary fields
	existingProduct.ID = updatedProduct.ID
	existingProduct.Name = updatedProduct.Name
	existingProduct.Type = updatedProduct.Type
	existingProduct.Quantity = updatedProduct.Quantity

	// Update the product in the repository
	err = ctrl.Repo.UpdateProduct(existingProduct, uint(productID))
	if err != nil {
		models.Logger.Error("Error updating product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	models.Logger.Info("Product updated successfully")
	c.JSON(http.StatusOK, existingProduct)
}

func (ctrl *ProductController) GetAllProductsHandler(c *gin.Context) {
	var products []Product

	// Fetch all products from the repository
	allProducts, err := ctrl.Repo.GetAllProducts()
	if err != nil {
		models.Logger.Error("Error getting all products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get products"})
		return
	}

	products = allProducts

	c.JSON(http.StatusOK, products)
}

func (ctrl *ProductController) GetSpecificProductHandler(c *gin.Context) {
	productIDStr := c.Param("id")
	productID32, err := strconv.ParseUint(productIDStr, 10, 64)
	productID := uint(productID32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Check if the product exists
	product, err := ctrl.Repo.GetProductByID(productID)
	if err != nil {
		models.Logger.Error("Error retrieving product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}
