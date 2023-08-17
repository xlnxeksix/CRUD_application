package user

import (
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, userController *Controller) {

	users := r.Group("/users")
	{
		users.POST("/", userController.CreateUserHandler)
		users.GET("/:id", userController.GetSpecificUserHandler)
		users.GET("/", userController.GetAllUsersHandler)
		users.PUT("/:id", userController.UpdateUserHandler)
		users.DELETE("/:id", userController.DeleteUserHandler)
	}
}
