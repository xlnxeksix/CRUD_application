package database

import (
	"awesomeProject1/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type DBController struct {
	DB *gorm.DB
}

func NewDBController(db *gorm.DB) *DBController {
	return &DBController{DB: db}
}
func (ctrl *DBController) GetUsersWithUsername(c *gin.Context) {
	var users []models.User

	// SQL query to fetch users who have the product "pen"
	query := `SELECT *
			  FROM users u
			  WHERE u.Username = 'first_user'`

	// Execute the raw SQL query and scan the results into the users slice
	if err := ctrl.DB.Raw(query).Scan(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}
