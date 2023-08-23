package Authentication

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	Repo AuthRepository
}

func NewAuthController(repo AuthRepository) *Controller {
	return &Controller{Repo: repo}
}

func (ctrl Controller) BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetHeader("Username")
		password := c.GetHeader("Password")

		role, err := ctrl.Repo.AuthenticateUser(username, password)
		if err != nil || role == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			c.Abort()
			return
		}
		c.Set("user_role", role)
		c.Next()
	}
}

func (ctrl Controller) AdminAuthMiddleware(c *gin.Context) {
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
func (ctrl Controller) UserAuthMiddleware(c *gin.Context) {
	userRole, existsRole := c.Get("user_role")
	if !existsRole {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		c.Abort()
		return
	}
	role := userRole.(string)
	if role != "user" || role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		c.Abort()
		return
	}
	c.Next()
}
