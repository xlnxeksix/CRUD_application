package routers

import (
	"awesomeProject1/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userController *controllers.UserController, productController *controllers.ProductController, userProductController *controllers.UserProductController) {
	users := r.Group("/users")
	{
		users.POST("/", userController.CreateUserHandler)
		users.GET("/:id", userController.GetSpecificUserHandler)
		users.GET("/", userController.GetAllUsersHandler)
		users.PUT("/:id", userController.UpdateUserHandler)
		users.DELETE("/:id", userController.DeleteUserHandler)
		users.GET("/get_first_user", userController.GetUsersWithUsername)
	}

	products := r.Group("/products")
	{
		products.POST("/", productController.CreateproductHandler)
		products.GET("/:id", productController.GetSpecificproductHandler)
		products.GET("/", productController.GetAllproductsHandler)
		products.PUT("/:id", productController.UpdateproductHandler)
		products.DELETE("/:id", productController.DeleteproductHandler)
	}

	userproducts := r.Group("/userProducts")
	{
		userproducts.POST("/", userProductController.CreateUserProductHandler)
	}
}
