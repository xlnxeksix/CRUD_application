package controllers

import (
	"awesomeProject1/models"
	"github.com/gin-gonic/gin"
)

type UserProductController struct {
	UserRepo     UserRepository
	ProductRepo  ProductRepository
	UserProdRepo UserProductRepository
}

func NewUserProductController(userRepo UserRepository, productRepo ProductRepository, userProdRepo UserProductRepository) *UserProductController {
	return &UserProductController{
		UserRepo:     userRepo,
		ProductRepo:  productRepo,
		UserProdRepo: userProdRepo,
	}
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

	var user *models.User
	user, usererr := ctrl.UserRepo.GetUserByID(requestBody.UserID)
	if usererr != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	if _, err := ctrl.ProductRepo.GetProductByID(requestBody.ProductID); err != nil {
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
	if err := ctrl.UserProdRepo.AllocateProduct(requestBody.UserID, requestBody.ProductID); err != nil {
		c.JSON(500, gin.H{"error": "Internal error allocating product"})
		return
	}

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

	var user *models.User
	user, usererr := ctrl.UserRepo.GetUserByID(requestBody.UserID)
	if usererr != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	if _, err := ctrl.ProductRepo.GetProductByID(requestBody.ProductID); err != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
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
	// Allocate the product to the user
	if err := ctrl.UserProdRepo.DeallocateProduct(requestBody.UserID, requestBody.ProductID); err != nil {
		c.JSON(500, gin.H{"error": "Internal error deallocating product"})
		return
	}

	c.JSON(200, gin.H{"message": "Product deallocated from the user"})
}
