package main

import (
	"awesomeProject1/controllers"
	"awesomeProject1/models"
	"awesomeProject1/routers"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Connect to the database
	dsn := "host=localhost user=postgres password=docker dbname=CRUD-db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Unable to connect to the database")
	}

	//Create table automatically
	db.AutoMigrate(&models.User{}, &models.Product{}, &models.UserProduct{})

	if err != nil {
		panic("There is an error when creating table")
	}

	// Create gin-gonic router
	router := gin.Default()

	//Create Logger defined in logger.go
	models.InitLogger()
	models.CloseLogger()

	userController := controllers.NewUserController(db)
	productController := controllers.NewProductController(db)
	userProductController := controllers.NewUserProductController(db)

	routers.SetupRoutes(router, userController, productController, userProductController)
	// Run the app
	models.Logger.Info("Application started succesfully")
	router.Run(":8080")

}
