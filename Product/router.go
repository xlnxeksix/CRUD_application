package product

import (
	"github.com/gin-gonic/gin"
)

func SetupProductRoutes(r *gin.Engine, productController *Controller) {

	products := r.Group("/products")
	{
		products.POST("/", productController.CreateProductHandler)
		products.GET("/:id", productController.GetSpecificProductHandler)
		products.GET("/", productController.GetAllProductsHandler)
		products.PUT("/:id", productController.UpdateProductHandler)
		products.DELETE("/:id", productController.DeleteProductHandler)
	}
}
