package routers

import (
	"awesomeProject1/Product"
	"awesomeProject1/User"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userController *user.Controller, productController *product.Controller) {
	users := r.Group("/users")
	{
		users.POST("/", userController.CreateUserHandler)
		users.POST("/prudcts", userController.CreateUserHandler)

		users.GET("/:id", userController.GetSpecificUserHandler)
		users.GET("/", userController.GetAllUsersHandler)
		users.PUT("/:id", userController.UpdateUserHandler)
		users.DELETE("/:id", userController.DeleteUserHandler)
	}

	products := r.Group("/products")
	{
		products.POST("/", productController.CreateProductHandler)
		products.GET("/:id", productController.GetSpecificProductHandler)
		products.GET("/", productController.GetAllProductsHandler)
		products.PUT("/:id", productController.UpdateProductHandler)
		products.DELETE("/:id", productController.DeleteProductHandler)
	}
}
