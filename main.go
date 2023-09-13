package main

import (
	"awesomeProject1/Authentication"
	"awesomeProject1/Models"
	"awesomeProject1/Product"
	"awesomeProject1/SIEM"
	"awesomeProject1/SIEM/Model"
	"awesomeProject1/User"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	// Connect to the database
	dsn := "host=localhost user=postgres password=docker dbname=CRUD-db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Unable to connect to the database")
	}

	//Create table automatically
	db.AutoMigrate(&user.User{}, &product.Product{}, &Model.RuleForm{}, &Model.FlattenedRule{}, &Model.InsightType{})

	if err != nil {
		panic("There is an error when creating table")
	}
	//rule_insight.AddInsightsToDatabase(db)
	// Create gin-gonic router
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	//Create Logger defined in logger.go
	models.InitLogger()
	models.CloseLogger()

	userRepo := &user.SQLUserRepository{DB: db}
	userController := user.NewUserController(userRepo)

	productRepo := &product.SQLProductRepository{DB: db}
	productController := product.NewProductController(productRepo)

	SIEMRepo := &SIEM.SQLRuleRepository{DB: db}
	SIEMController := SIEM.NewSIEMController(SIEMRepo)

	authRepo := &Authentication.SQLAuthRepository{DB: db}
	authController := Authentication.NewAuthController(authRepo)

	user.SetupUserRoutes(router, authController, userController)
	product.SetupProductRoutes(router, authController, productController)
	SIEM.SetupInsightRoutes(router, authController, SIEMController)
	// Define a route to render the SIEM rule input form
	router.GET("/siem-form", func(c *gin.Context) {
		c.HTML(http.StatusOK, "siem_form.html", nil)
	})
	// Run the app
	models.Logger.Info("Application started successfully")
	router.Run(":8080")
}
