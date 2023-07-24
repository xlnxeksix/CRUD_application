package controllers

import (
	"awesomeProject1/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

type UserProductController struct {
	DB *gorm.DB
}

func NewUserProductController(db *gorm.DB) *UserProductController {
	return &UserProductController{DB: db}
}

func (ctrl *UserProductController) CreateUserProductHandler(c *gin.Context) {
	var user_product models.UserProduct
	if err := c.ShouldBindJSON(&user_product); err != nil {
		models.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := ctrl.DB.Create(&user_product)
	if result.Error != nil {
		models.Logger.Error("Error creating user", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	models.Logger.Info("User created successfully")
	c.JSON(http.StatusCreated, user_product)
}
