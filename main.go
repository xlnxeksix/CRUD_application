package main

import (
	"awesomeProject1/api"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Connect to the database
	dsn := "host=localhost user=postgres password=docker dbname=user_database port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Unable to connect to the database")
	}

	//Create table automatically
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic("There is an error when creating table")
	}

	// Create gin-gonic router
	router := gin.Default()

	// Creating a user
	createUserHandler := models.CreateUserHandler(db)
	router.POST("/users", createUserHandler)

	//Get user by ID
	getSpecificUserHandler := models.GetSpecificUserHandler(db)
	router.GET("/users/:id", getSpecificUserHandler)

	// Get all users
	getAllUsersHandler := models.GetAllUsersHandler(db)
	router.GET("/users", getAllUsersHandler)

	// Update a user
	updateUserHandler := models.UpdateUserHandler(db)
	router.PUT("/users/:id", updateUserHandler)

	// Delete a user
	deleteUserHandler := models.DeleteUserHandler(db)
	router.DELETE("/users/:id", deleteUserHandler)

	// Run the app
	router.Run(":8080")

}
