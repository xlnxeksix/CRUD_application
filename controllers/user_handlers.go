package controllers

import (
	"awesomeProject1/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

func (ctrl *UserController) CreateUserHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		models.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := ctrl.DB.Create(&user)
	if result.Error != nil {
		models.Logger.Error("Error creating user", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	models.Logger.Info("User created successfully")
	c.JSON(http.StatusCreated, user)
}

func (ctrl *UserController) DeleteUserHandler(c *gin.Context) {
	var user models.User
	result := ctrl.DB.Delete(&user, c.Param("id"))
	if result.Error != nil {
		models.Logger.Error("Error deleting user", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	models.Logger.Info("User deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "User is deleted"})
}

func (ctrl *UserController) UpdateUserHandler(c *gin.Context) {
	var user models.User
	if err := ctrl.DB.First(&user, c.Param("id")).Error; err != nil {
		models.Logger.Error("Error finding user", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		models.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := ctrl.DB.Save(&user)
	if result.Error != nil {
		models.Logger.Error("Error updating user", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	models.Logger.Info("User updated successfully")
	c.JSON(http.StatusOK, user)
}

// GetAllUsersHandler handles getting all users
func (ctrl *UserController) GetAllUsersHandler(c *gin.Context) {
	var users []models.User
	result := ctrl.DB.Find(&users)
	if result.Error != nil {
		models.Logger.Error("Error getting all users", zap.Error(result.Error))
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetSpecificUserHandler handles getting a specific user
func (ctrl *UserController) GetSpecificUserHandler(c *gin.Context) {
	var user models.User
	result := ctrl.DB.First(&user, c.Param("id"))
	if result.Error != nil {
		models.Logger.Error("User not found", zap.Error(result.Error))
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
