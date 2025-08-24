package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/whoami/internal/security"
)

func SetupRoutes(router *gin.Engine, handler *HTTPHandler) {
	router.GET("/health", handler.HealthCheck)

	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/register",
			handler.rateLimiter.RateLimitMiddleware(security.RegistrationRateLimit),
			handler.Register)
		apiV1.POST("/login",
			handler.rateLimiter.RateLimitMiddleware(security.AuthRateLimit),
			handler.Login)

		passwordReset := apiV1.Group("/password-reset")
		passwordReset.Use(handler.rateLimiter.RateLimitMiddleware(security.PasswordResetRateLimit))
		{
			passwordReset.POST("/request", handler.RequestPasswordReset)
			passwordReset.POST("/reset", handler.ResetPassword)
			passwordReset.POST("/verify", handler.VerifyResetToken)
		}

		protected := apiV1.Group("/")
		protected.Use(AuthMiddleware(handler.tokenMaker))
		protected.Use(handler.rateLimiter.UserRateLimitMiddleware(security.DefaultRateLimit))
		{
			protected.GET("/me", handler.GetCurrentUser)
			protected.POST("/logout", handler.Logout)
			protected.POST("/refresh", handler.RefreshToken)

			user := protected.Group("/user")
			{
				user.POST("/deactivate", handler.DeactivateUser)
				user.POST("/activate", handler.ActivateUser)
				user.PUT("/:id", handler.UpdateUser)
				user.PUT("/:id/privacy-settings", handler.UpdateUserPrivacySettings)
				user.POST("/update-password", handler.UpdatePassword)
			}

			security := protected.Group("/security")
			{
				security.GET("/activities", handler.GetSuspiciousActivities)
				security.POST("/activities/resolve", handler.ResolveSuspiciousActivity)
				security.POST("/cleanup", handler.CleanupExpiredLockouts)
			}
		}

		email := apiV1.Group("/email")
		email.Use(handler.rateLimiter.RateLimitMiddleware(security.PasswordResetRateLimit))
		{
			email.POST("/verify", handler.VerifyEmail)
			email.POST("/resend", handler.ResendVerificationEmail)
		}
	}
}
