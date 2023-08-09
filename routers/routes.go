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
	}

	products := r.Group("/products")
	{
		products.POST("/", productController.CreateProductHandler)
		products.GET("/:id", productController.GetSpecificproductHandler)
		products.GET("/", productController.GetAllProductsHandler)
		products.PUT("/:id", productController.UpdateProductHandler)
		products.DELETE("/:id", productController.DeleteProductHandler)
	}

	userproducts := r.Group("/userProducts")
	{
		userproducts.POST("/allocate", userProductController.AllocateProductToUser)
		userproducts.POST("/deallocate", userProductController.DeallocateProductFromUser)
	}
}
