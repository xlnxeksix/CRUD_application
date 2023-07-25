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

	insertQuery := fmt.Sprintf("INSERT INTO users (id, username, email, role) VALUES (%d, '%s', '%s', '%s')",
		user.ID, user.Username, user.Email, user.Role)

	if err := ctrl.DB.Exec(insertQuery).Error; err != nil {
		models.Logger.Error("Error creating user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	models.Logger.Info("User created successfully")
	c.JSON(http.StatusCreated, user)
}

func (ctrl *UserController) DeleteUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		models.Logger.Error("Invalid ID provided", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID provided"})
		return
	}

	// Convert the userID to a string before using it in the SQL query
	deleteQuery := fmt.Sprintf("DELETE FROM users WHERE id = '%d'", userID)

	if err := ctrl.DB.Exec(deleteQuery).Error; err != nil {
		models.Logger.Error("Error deleting user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	models.Logger.Info("User deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "User is deleted"})
}

func (ctrl *UserController) UpdateUserHandler(c *gin.Context) {
	var user models.User
	idParam := c.Param("id")
	userID, _ := strconv.ParseUint(idParam, 10, 64)
	query := fmt.Sprintf("SELECT * FROM users WHERE id = '%d' FOR UPDATE", userID)

	row := ctrl.DB.Raw(query).Row()
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

	updateQuery := fmt.Sprintf("UPDATE users SET ID = %d, username = '%s', email = '%s', role = '%s' WHERE id = '%d'",
		user.ID, user.Username, user.Email, user.Role, userID)
	ctrl.DB.Exec(updateQuery)

	/*if err, _ :=  err != nil {
		models.Logger.Error("Error updating user")
		models.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}*/

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
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	query := fmt.Sprintf("SELECT * FROM users WHERE id = '%d'", userID)

	rows, err := ctrl.DB.Raw(query).Rows()
	if err != nil {
		models.Logger.Error("Error finding user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role)
		if err != nil {
			models.Logger.Error("Error scanning user row", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			return
		}
		break // Assuming that there will be only one row since ID is unique
	}

	c.JSON(http.StatusOK, user)
}
