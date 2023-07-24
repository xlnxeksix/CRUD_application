package controllers

import (
	"awesomeProject1/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

type ProductController struct {
	DB *gorm.DB
}

func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{DB: db}
}

func (ctrl *ProductController) CreateproductHandler(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		models.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := ctrl.DB.Create(&product)
	if result.Error != nil {
		models.Logger.Error("Error creating product", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	models.Logger.Info("product created successfully")
	c.JSON(http.StatusCreated, product)
}

func (ctrl *ProductController) DeleteproductHandler(c *gin.Context) {
	var product models.Product
	result := ctrl.DB.Delete(&product, c.Param("id"))
	if result.Error != nil {
		models.Logger.Error("Error deleting product", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	models.Logger.Info("product deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "product is deleted"})
}

func (ctrl *ProductController) UpdateproductHandler(c *gin.Context) {
	var product models.Product
	if err := ctrl.DB.First(&product, c.Param("id")).Error; err != nil {
		models.Logger.Error("Error finding product", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		models.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := ctrl.DB.Save(&product)
	if result.Error != nil {
		models.Logger.Error("Error updating product", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	models.Logger.Info("product updated successfully")
	c.JSON(http.StatusOK, product)
}

// GetAllproductsHandler handles getting all products
func (ctrl *ProductController) GetAllproductsHandler(c *gin.Context) {
	var products []models.Product
	result := ctrl.DB.Find(&products)
	if result.Error != nil {
		models.Logger.Error("Error getting all products", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// GetSpecificproductHandler handles getting a specific product
func (ctrl *ProductController) GetSpecificproductHandler(c *gin.Context) {
	var product models.Product
	result := ctrl.DB.First(&product, c.Param("id"))
	if result.Error != nil {
		models.Logger.Error("product not found", zap.Error(result.Error))
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}
