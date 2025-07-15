package handlers

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, handler *HTTPHandler) {
	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/register", handler.Register)
		apiV1.GET("/me", handler.GetCurrentUser)
		apiV1.POST("/login", handler.Login)
		user := apiV1.Group("/user")
		{
			user.POST("/deactivate", handler.DeactivateUser)
			user.POST("/activate", handler.ActivateUser)
			user.PUT("/:id", handler.UpdateUser)
			user.PUT("/:id/privacy-settings", handler.UpdateUserPrivacySettings)
		}
	}
}
