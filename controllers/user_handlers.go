package controllers

import (
	"awesomeProject1/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

// CreateUserHandler handles creating a new user
func (ctrl *UserController) CreateUserHandler(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		models.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	insertQuery := "INSERT INTO users (id, username, email, role) VALUES (?, ?, ?, ?)"

	if err := ctrl.DB.Exec(insertQuery, user.ID, user.Username, user.Email, user.Role).Error; err != nil {
		models.Logger.Error("Error creating user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	models.Logger.Info("User created successfully")
	c.JSON(http.StatusCreated, user)
}

func (ctrl *UserController) DeleteUserHandler(c *gin.Context) {
	// Check if the user exists before attempting to delete
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

	deleteQuery := "DELETE FROM users WHERE id = ?"
	err = ctrl.DB.Exec(deleteQuery, c.Param("id")).Error
	if err != nil {
		models.Logger.Error("Error deleting user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	models.Logger.Info("User deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "User is deleted"})
}

func (ctrl *UserController) UpdateUserHandler(c *gin.Context) {
	var user models.User

	query := "SELECT * FROM users WHERE id = ?"
	err := ctrl.DB.Exec(query, c.Param("id")).Error
	if err != nil {
		models.Logger.Error("Error updating user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}
	row := ctrl.DB.Raw(query, c.Param("id")).Row()

	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role); err != nil {
		models.Logger.Error("Error finding user", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		models.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateQuery := "UPDATE users SET ID = ?, username = ?, email = ?, role = ? WHERE id = ?"
	err = ctrl.DB.Exec(updateQuery, user.ID, user.Username, user.Email, user.Role, c.Param("id")).Error

	if err != nil {
		models.Logger.Error("Error updating user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	models.Logger.Info("User updated successfully")
	c.JSON(http.StatusOK, user)
}

// GetAllUsersHandler handles getting all users
func (ctrl *UserController) GetAllUsersHandler(c *gin.Context) {
	var users []models.User

	query := "SELECT * FROM users"

	rows, err := ctrl.DB.Raw(query).Rows()
	if err != nil {
		models.Logger.Error("Error getting all users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role); err != nil {
			models.Logger.Error("Error scanning user row", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

// GetSpecificUserHandler handles getting a specific user

func (ctrl *UserController) GetSpecificUserHandler(c *gin.Context) {
	var user models.User
	query := "SELECT * FROM users WHERE id = ?"
	rows, err := ctrl.DB.Raw(query, c.Param("id")).Rows()
	if err != nil {
		models.Logger.Error("Error finding user", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID must be an integer"})
		return
	}
	defer rows.Close()

	found := false // Variable to keep track if the user was found or not

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role)
		if err != nil {
			models.Logger.Error("Error scanning user row", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			return
		}
		found = true // User found, set the flag to true
		break        // Assuming that there will be only one row since ID is unique
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Preload the allocated products for the user
	ctrl.DB.Model(&user).Preload("Products").Find(&user)
	c.JSON(http.StatusOK, user)
}
