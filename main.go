package main

import (
	"awesomeProject1/Authentication"
	"awesomeProject1/Models"
	"awesomeProject1/Product"
	"awesomeProject1/User"
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
	db.AutoMigrate(&user.User{}, &product.Product{})

	if err != nil {
		panic("There is an error when creating table")
	}

	// Create gin-gonic router
	router := gin.Default()

	//Create Logger defined in logger.go
	models.InitLogger()
	models.CloseLogger()

	userRepo := &user.SQLUserRepository{DB: db}
	userController := user.NewUserController(userRepo)

	productRepo := &product.SQLProductRepository{DB: db}
	productController := product.NewProductController(productRepo)

	authRepo := &Authentication.SQLAuthRepository{DB: db}
	authController := Authentication.NewAuthController(authRepo)

	user.SetupUserRoutes(router, authController, userController)
	product.SetupProductRoutes(router, authController, productController)
	// Run the app
	models.Logger.Info("Application started successfully")
	router.Run(":8080")
}
