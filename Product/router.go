package product

import (
	"awesomeProject1/Authentication"
	"github.com/gin-gonic/gin"
)

func SetupProductRoutes(r *gin.Engine, authController *Authentication.Controller, productController *Controller) {
	r.Use(authController.BasicAuthMiddleware())
	adminGroup := r.Group("/products")
	adminGroup.Use(authController.AdminAuthMiddleware)
	{
		adminGroup.POST("/", productController.CreateProductHandler)
		adminGroup.GET("/", productController.GetAllProductsHandler)
		adminGroup.DELETE("/:id", productController.DeleteProductHandler)
	}
	userGroup := r.Group("/products")
	userGroup.Use(authController.UserAuthMiddleware)
	{
		userGroup.GET("/:id", productController.GetSpecificProductHandler)
		userGroup.PUT("/:id", productController.UpdateProductHandler)
	}
}
