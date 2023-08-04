package controllers

import (
	"awesomeProject1/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type UserProductController struct {
	DB *gorm.DB
}

func NewUserProductController(db *gorm.DB) *UserProductController {
	return &UserProductController{DB: db}
}

func (ctrl *UserProductController) CreateUserProductHandler(c *gin.Context) {
	var userProduct models.UserProduct

	if err := c.ShouldBindJSON(&userProduct); err != nil {
		models.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO user_products (id, user_id, product_id) VALUES (?, ?, ?)"

	if err := ctrl.DB.Exec(query, userProduct.ID, userProduct.UserID, userProduct.ProductID).Error; err != nil {
		models.Logger.Error("Error creating user-product association", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user-product association"})
		return
	}

	models.Logger.Info("User-product association created successfully")
	c.JSON(http.StatusCreated, userProduct)
}

// DeleteUserProductHandler handles deleting a specific user-product association
func (ctrl *UserProductController) DeleteUserProductHandler(c *gin.Context) {
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
	query := "DELETE FROM user_products WHERE id = ?"
	err = ctrl.DB.Exec(query, c.Param("id")).Error

	if err != nil {
		models.Logger.Error("Error deleting association", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete user-product association"})
		return
	}

	models.Logger.Info("user-pproduct association deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "user-product association is deleted"})
}

func (ctrl *UserProductController) UpdateUserProductHandler(c *gin.Context) {
	var userproduct models.UserProduct

	query := "SELECT * FROM user_products WHERE id = ?"
	err := ctrl.DB.Exec(query, c.Param("id")).Error
	if err != nil {
		models.Logger.Error("Error updating user-product associations", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user-product associations"})
		return
	}
	row := ctrl.DB.Raw(query, c.Param("id")).Row()

	if err := row.Scan(&userproduct.ID, userproduct.UserID, userproduct.ProductID); err != nil {
		models.Logger.Error("Error finding user-product associations", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "user-product associations not found"})
		return
	}

	if err := c.ShouldBindJSON(&userproduct); err != nil {
		models.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateQuery := "UPDATE user_products SET ID = ?, userID = ?, productID = ? WHERE id = ?"
	err = ctrl.DB.Exec(updateQuery, userproduct.ID, userproduct.UserID, userproduct.ProductID, c.Param("id")).Error

	if err != nil {
		models.Logger.Error("Error updating user-product associations")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user-product associations"})
		return
	}

	models.Logger.Info("user-product associations updated successfully")
	c.JSON(http.StatusOK, userproduct)
}

// GetAllproductsHandler handles getting all products
func (ctrl *UserProductController) GetAllUserProductsHandler(c *gin.Context) {
	var user_products []models.UserProduct

	query := "SELECT * FROM user_products"

	rows, err := ctrl.DB.Raw(query).Rows()
	if err != nil {
		models.Logger.Error("Error getting all user-product associations", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user-product associations"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user_product models.UserProduct

		if err := rows.Scan(&user_product.ID, &user_product.UserID, &user_product.ProductID); err != nil {
			models.Logger.Error("Error scanning user-product row", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user-products"})
			return

			user_products = append(user_products, user_product)
		}

		c.JSON(http.StatusOK, user_products)
	}
}

// Get Specific product Handler handles getting a specific product
func (ctrl *UserProductController) GetSpecificUserProductHandler(c *gin.Context) {
	var userproduct models.UserProduct
	query := "SELECT * FROM user_products WHERE id = ?"
	rows, err := ctrl.DB.Raw(query, c.Param("id")).Rows()
	if err != nil {
		models.Logger.Error("Error finding user-product", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "User-product ID must be an integer"})
		return
	}
	defer rows.Close()

	found := false // Variable to keep track if the user was found or not

	for rows.Next() {
		err := rows.Scan(&userproduct.ID, &userproduct.UserID, &userproduct.ProductID)
		if err != nil {
			models.Logger.Error("Error scanning user-product row", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user-product"})
			return
		}
		found = true // User found, set the flag to true
		break        // Assuming that there will be only one row since ID is unique
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "User-product not found"})
		return
	}

	c.JSON(http.StatusOK, userproduct)
}
