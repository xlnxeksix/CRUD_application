package SIEM

import (
	"awesomeProject1/Authentication"
	"github.com/gin-gonic/gin"
)

func SetupInsightRoutes(r *gin.Engine, authController *Authentication.Controller, SIEMController *Controller) {
	r.Use(authController.BasicAuthMiddleware())
	adminGroup := r.Group("/siem")
	adminGroup.Use(authController.AdminAuthMiddleware)
	{
		adminGroup.POST("/", SIEMController.GetRuleContent)
	}
}
