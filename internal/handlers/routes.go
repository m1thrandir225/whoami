package handlers

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, handler *HTTPHandler) {
	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/register", handler.Register)
		apiV1.POST("/login", handler.Login)

		protected := apiV1.Group("/")
		protected.Use(AuthMiddleware(handler.tokenMaker))
		{
			protected.GET("/me", handler.GetCurrentUser)
			protected.POST("/logout", handler.Logout)
			protected.POST("/refresh", handler.RefreshToken)
			protected.POST("/verify-email", handler.VerifyEmail)
			protected.POST("/resend-verification-email", handler.ResendVerificationEmail)
			protected.POST("/update-password", handler.UpdatePassword)

			user := protected.Group("/user")
			{
				user.POST("/deactivate", handler.DeactivateUser)
				user.POST("/activate", handler.ActivateUser)
				user.PUT("/:id", handler.UpdateUser)
				user.PUT("/:id/privacy-settings", handler.UpdateUserPrivacySettings)
			}
		}
	}
}
