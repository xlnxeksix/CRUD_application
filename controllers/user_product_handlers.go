package controllers

import (
	"awesomeProject1/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserProductController struct {
	DB *gorm.DB
}

func NewUserProductController(db *gorm.DB) *UserProductController {
	return &UserProductController{DB: db}
}
func (ctrl *UserProductController) AllocateProductToUser(c *gin.Context) {
	var requestBody struct {
		UserID    uint `json:"user_id"`
		ProductID uint `json:"product_id"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	var user models.User
	if err := ctrl.DB.First(&user, requestBody.UserID).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	var product models.Product
	if err := ctrl.DB.First(&product, requestBody.ProductID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}

	// Check if the product is already allocated to the user
	for _, allocatedProduct := range user.Products {
		if allocatedProduct.ID == requestBody.ProductID {
			c.JSON(409, gin.H{"error": "Product is already allocated to the user"})
			return
		}
	}

	// Allocate the product to the user
	user.Products = append(user.Products, &product)
	product.Users = append(product.Users, &user)
	ctrl.DB.Save(&user)
	ctrl.DB.Save(&product)

	c.JSON(200, gin.H{"message": "Product allocated to the user"})
}

func (ctrl *UserProductController) DeallocateProductFromUser(c *gin.Context) {
	var requestBody struct {
		UserID    uint `json:"user_id"`
		ProductID uint `json:"product_id"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	var user models.User
	if err := ctrl.DB.Preload("Products").First(&user, requestBody.UserID).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Find the product to be removed from the user
	var productToRemove *models.Product
	for _, allocatedProduct := range user.Products {
		if allocatedProduct.ID == requestBody.ProductID {
			productToRemove = allocatedProduct
			break
		}
	}

	if productToRemove == nil {
		c.JSON(404, gin.H{"error": "Product is not allocated to the user"})
		return
	}

	// Remove the product from the user's allocation
	ctrl.DB.Model(&user).Association("Products").Delete(productToRemove)
	ctrl.DB.Model(&productToRemove).Association("Users").Delete(user)

	c.JSON(200, gin.H{"message": "Product deallocated from the user"})
}
