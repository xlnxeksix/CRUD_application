package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Email    string
	Role     string
}

func main() {
	// Connect to the database
	dsn := "host=localhost user=postgres password=docker dbname=user_database port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Unable to connect to the database")
	}

	//Create table automatically
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("There is an error when creating table")
	}

	// Create gin-gonic router
	router := gin.Default()

	// Create a file output
	file, err := os.OpenFile("app.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic("Failed to create log file")
	}
	defer file.Close()

	// Create a Zap logger configuration
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Customize the time format if needed
	config.OutputPaths = []string{file.Name()}
	config.ErrorOutputPaths = []string{file.Name()}

	// Create the logger
	logger, err := config.Build()
	if err != nil {
		panic("Failed to create logger")
	}
	defer logger.Sync()

	// ...
	// Creating a user
	router.POST("/users", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			logger.Error("Failed to bind JSON data", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result := db.Create(&user)
		if result.Error != nil {
			logger.Error("Failed to create user", zap.Error(result.Error))
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		logger.Info("User created", zap.Any("user", user))
		c.JSON(http.StatusCreated, user)
	})

	// Get all users
	router.GET("/users", func(c *gin.Context) {
		var users []User
		result := db.Find(&users)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		logger.Info("All users are returned")
		c.JSON(http.StatusOK, users)
	})

	// Get a spesific user
	router.GET("/users/:id", func(c *gin.Context) {
		var user User
		result := db.First(&user, c.Param("id"))
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	})

	// Update a user
	router.PUT("/users/:id", func(c *gin.Context) {
		var user User
		if err := db.First(&user, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result := db.Save(&user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	})

	// Delete a user
	router.DELETE("/users/:id", func(c *gin.Context) {
		var user User
		result := db.Delete(&user, c.Param("id"))
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User is deleted"})
	})

	// Run the app
	logger.Info("Application version 1 started")
	router.Run(":8080")

}
