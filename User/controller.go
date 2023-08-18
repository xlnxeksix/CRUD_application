package user

import (
	"awesomeProject1/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Controller struct {
	Repo UserRepository
}

func NewUserController(repo UserRepository) *Controller {
	return &Controller{Repo: repo}
}

// CreateUserHandler handles creating a new user
func (ctrl *Controller) CreateUserHandler(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		models.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Check if the user exists
	if existingUser, _ := ctrl.Repo.GetUserByID(user.ID); existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with same ID exists"})
		return
	}
	if err := ctrl.Repo.CreateUser(&user); err != nil {
		models.Logger.Error("Error creating user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	models.Logger.Info("User created successfully")
	c.JSON(http.StatusCreated, user)
}

// DeleteUserHandler handles deleting a user
func (ctrl *Controller) DeleteUserHandler(c *gin.Context) {
	userIDStr := c.Param("id")
	userID32, err := strconv.ParseUint(userIDStr, 10, 64)
	userID := uint(userID32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	// Check if the user exists
	user, err := ctrl.Repo.GetUserByID(userID)
	if err != nil {
		models.Logger.Error("Error retrieving user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Delete the user
	if err := ctrl.Repo.DeleteUser(userID); err != nil {
		models.Logger.Error("Error deleting user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	models.Logger.Info("User deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (ctrl *Controller) UpdateUserHandler(c *gin.Context) {
	var updatedUser User

	// Check if the user exists
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Fetch the existing user from the repository
	existingUser, err := ctrl.Repo.GetUserByID(uint(userID))
	if err != nil {
		models.Logger.Error("Error finding user", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Bind JSON data to the updatedUser struct
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		models.Logger.Error("Error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update only the necessary fields
	existingUser.Username = updatedUser.Username
	existingUser.Email = updatedUser.Email
	existingUser.Role = updatedUser.Role

	// Update the user in the repository
	err = ctrl.Repo.UpdateUser(existingUser, uint(userID))
	if err != nil {
		models.Logger.Error("Error updating user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	models.Logger.Info("User updated successfully")
	c.JSON(http.StatusOK, existingUser)
}

func (ctrl *Controller) GetAllUsersHandler(c *gin.Context) {
	var users []User

	// Fetch all users from the repository
	allUsers, err := ctrl.Repo.GetAllUsers()
	if err != nil {
		models.Logger.Error("Error getting all users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	users = allUsers

	c.JSON(http.StatusOK, users)
}

// GetSpecificUserHandler handles getting a specific user

func (ctrl *Controller) GetSpecificUserHandler(c *gin.Context) {

	userIDStr := c.Param("id")
	userID32, err := strconv.ParseUint(userIDStr, 10, 64)
	userID := uint(userID32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	// Check if the user exists
	user, err := ctrl.Repo.GetUserByID(userID)
	if err != nil {
		models.Logger.Error("Error retrieving user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
