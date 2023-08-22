package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func BasicAuthMiddleware(userRepo UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("Username")
		password := c.GetHeader("Password")

		// Authenticate user using the repository
		user, err := userRepo.AuthenticateUser(username, password)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			c.Abort()
			return
		}

		// Retrieve user's role from the database
		userRole, err := userRepo.GetUserRole(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user role"})
			c.Abort()
			return
		}
		// Store authenticated user and role in context for further use
		c.Set("authenticated_user", user)
		c.Set("user_role", userRole)

		c.Next()
	}
}

func RoleBasedAuthMiddleware(c *gin.Context) {
	// Get the user's role from the context
	userRole, existsRole := c.Get("user_role")
	if !existsRole {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		c.Abort()
		return
	}

	role := userRole.(string)

	// Check the user's role and authorize access
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		c.Abort()
		return
	}

	c.Next()
}

func SetupUserRoutes(r *gin.Engine, userController *Controller, userRepo UserRepository) {
	r.Use(BasicAuthMiddleware(userRepo))
	adminGroup := r.Group("/user")
	adminGroup.Use(RoleBasedAuthMiddleware)
	{
		adminGroup.POST("/", userController.CreateUserHandler)
		adminGroup.GET("/", userController.GetAllUsersHandler)
		adminGroup.DELETE("/:id", userController.DeleteUserHandler)
	}

	userGroup := r.Group("/user")
	{
		userGroup.GET("/:id", userController.GetSpecificUserHandler)
		userGroup.PUT("/:id", userController.UpdateUserHandler)
	}
	sudo := r.Group("/create_admin")
	{
		sudo.POST("/", userController.CreateUserHandler)
	}
}
