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

	query := fmt.Sprintf("INSERT INTO user_products (id, user_id, product_id) VALUES (%d, %d, %d)", userProduct.ID, userProduct.UserID, userProduct.ProductID)

	row := ctrl.DB.Raw(query).Row()
	row.Scan(&userProduct.ID)
	/*
		if err :=  err != nil {
			models.Logger.Error("Error creating user-product association", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user-product association"})
			return
		}
	*/
	models.Logger.Info("User-product association created successfully")
	c.JSON(http.StatusCreated, userProduct)
}

// DeleteUserProductHandler handles deleting a specific user-product association
func (ctrl *UserProductController) DeleteUserProductHandler(c *gin.Context) {
	idParam := c.Param("id")
	userProductID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		models.Logger.Error("Invalid ID provided", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID provided"})
		return
	}
	query := fmt.Sprintf("DELETE FROM user_products WHERE id = '%d'", userProductID)

	ctrl.DB.Exec(query)
	/*
		if err != nil {
			models.Logger.Error("Error deleting product")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user-product association"})
			return
		}
	*/
	models.Logger.Info("product deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "product is deleted"})
}

func (ctrl *UserProductController) UpdateUserProductHandler(c *gin.Context) {
	var userProduct models.UserProduct
	idParam := c.Param("id")
	userProductID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		models.Logger.Error("Invalid ID provided", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID provided"})
		return
	}
	query := fmt.Sprintf("SELECT * FROM user_products WHERE id = '%d' FOR UPDATE", userProductID)

	row := ctrl.DB.Raw(query).Row()
	if err := row.Scan(&userProduct.ID, &userProduct.UserID, &userProduct.ProductID); err != nil {
		models.Logger.Error("user-product association not found", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "user-product association not found"})
		return
	}

	if err := c.ShouldBindJSON(&userProduct); err != nil {
		models.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateQuery := fmt.Sprintf("UPDATE user_products SET id = '%d', user_id = %d, product_id = %d WHERE id = '%d'",
		userProduct.ID, userProduct.UserID, userProduct.ProductID, userProductID)
	ctrl.DB.Exec(updateQuery)
	/*
		if err := ; err != nil {
			models.Logger.Error("Error updating product")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user-product association"})
			return
		}
	*/
	models.Logger.Info("product updated successfully")
	c.JSON(http.StatusOK, userProduct)
}

// GetAllproductsHandler handles getting all products
func (ctrl *UserProductController) GetAllUserProductsHandler(c *gin.Context) {
	var user_products []models.UserProduct

	query := "SELECT * FROM user_products"

	rows, err := ctrl.DB.Raw(query).Rows()
	if err != nil {
		models.Logger.Error("Error getting all products", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get products"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user_product models.UserProduct
		rows.Scan(&user_product.ID, &user_product.UserID, &user_product.ProductID)
		/*
			if err := ; err != nil {
				models.Logger.Error("Error scanning product row", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get products"})
				return
			}*/
		user_products = append(user_products, user_product)
	}

	c.JSON(http.StatusOK, user_products)
}

// Get Specific product Handler handles getting a specific product
func (ctrl *UserProductController) GetSpecificUserProductHandler(c *gin.Context) {
	var userProduct models.UserProduct
	idParam := c.Param("id")
	userproductID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		models.Logger.Error("Invalid ID provided", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID provided"})
		return
	}
	query := fmt.Sprintf("SELECT * FROM user_products WHERE id = '%d'", userproductID)

	row := ctrl.DB.Raw(query).Row()
	row.Scan(&userProduct.ID, &userProduct.UserID, &userProduct.ProductID)
	/*
		if err := ; err != nil {
			models.Logger.Error("user-product association not found", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "user-product association not found"})
			return
		}*/
	c.JSON(http.StatusOK, userProduct)
}
