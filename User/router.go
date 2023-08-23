package user

import (
	"awesomeProject1/Authentication"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, authController *Authentication.Controller, userController *Controller) {
	r.Use(authController.BasicAuthMiddleware())
	adminGroup := r.Group("/user")
	adminGroup.Use(authController.AdminAuthMiddleware)
	{
		adminGroup.POST("/", userController.CreateUserHandler)
		adminGroup.GET("/", userController.GetAllUsersHandler)
		adminGroup.DELETE("/:id", userController.DeleteUserHandler)
	}
	userGroup := r.Group("/user")
	userGroup.Use(authController.UserAuthMiddleware)
	{
		userGroup.GET("/:id", userController.GetSpecificUserHandler)
		userGroup.PUT("/:id", userController.UpdateUserHandler)
	}
	sudo := r.Group("/create_admin")
	{
		sudo.POST("/", userController.CreateUserHandler)
	}
}
