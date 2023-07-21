package models

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
)

func CreateUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			Logger.Error("Error binding JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result := db.Create(&user)
		if result.Error != nil {
			Logger.Error("Error creating user", zap.Error(result.Error))
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		Logger.Info("User created successfully")
		c.JSON(http.StatusCreated, user)
	}
}

func DeleteUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		result := db.Delete(&user, c.Param("id"))
		if result.Error != nil {
			Logger.Error("Error deleting user", zap.Error(result.Error))
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		Logger.Info("User deleted successfully")
		c.JSON(http.StatusOK, gin.H{"message": "User is deleted"})
	}
}

func UpdateUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		if err := db.First(&user, c.Param("id")).Error; err != nil {
			Logger.Error("Error finding user", zap.Error(err))
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		if err := c.ShouldBindJSON(&user); err != nil {
			Logger.Error("Error binding JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result := db.Save(&user)
		if result.Error != nil {
			Logger.Error("Error updating user", zap.Error(result.Error))
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		Logger.Info("User updated successfully")
		c.JSON(http.StatusOK, user)
	}
}

// GetAllUsersHandler handles getting all users
func GetAllUsersHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []User
		result := db.Find(&users)
		if result.Error != nil {
			Logger.Error("Error getting all users", zap.Error(result.Error))
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

// GetSpecificUserHandler handles getting a specific user
func GetSpecificUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		result := db.First(&user, c.Param("id"))
		if result.Error != nil {
			Logger.Error("User not found", zap.Error(result.Error))
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
