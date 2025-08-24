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
			passwordReset.POST("/verify", handler.VerifyResetToken)
			passwordReset.POST("/reset", handler.ResetPassword)
		}

		oauth := apiV1.Group("/oauth")
		{
			oauth.GET("/login/:provider", handler.OAuthLogin)
			oauth.GET("/callback/:provider", handler.OAuthCallback)
		}

		protected := apiV1.Group("/")
		protected.Use(AuthMiddleware(handler.tokenMaker, handler.tokenBlacklist))
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

			sessions := protected.Group("/sessions")
			{
				sessions.GET("", handler.GetUserSessions)
				sessions.DELETE("/:token", handler.RevokeSession)
				sessions.DELETE("", handler.RevokeAllSessions)
			}

			security := protected.Group("/security")
			{
				security.GET("/activities", handler.GetSuspiciousActivities)
				security.POST("/activities/resolve", handler.ResolveSuspiciousActivity)
				security.POST("/cleanup", handler.CleanupExpiredLockouts)
			}

			audit := protected.Group("/audit")
			{
				audit.GET("/user/:user_id", handler.GetAuditLogsByUserID)
				audit.GET("/action/:action", handler.GetAuditLogsByAction)
				audit.GET("/resource/:resource_type", handler.GetAuditLogsByResourceType)
				audit.GET("/resource/:resource_type/:resource_id", handler.GetAuditLogsByResourceID)
				audit.GET("/ip/:ip", handler.GetAuditLogsByIP)
				audit.GET("/date-range", handler.GetAuditLogsByDateRange)
				audit.GET("/recent", handler.GetRecentAuditLogs)
				audit.POST("/cleanup", handler.CleanupOldAuditLogs)
			}

			// User Devices routes
			devices := protected.Group("/devices")
			{
				devices.GET("", handler.GetUserDevices)
				devices.GET("/:id", handler.GetUserDevice)
				devices.PUT("/:id", handler.UpdateUserDevice)
				devices.DELETE("/:id", handler.DeleteUserDevice)
				devices.DELETE("", handler.DeleteAllUserDevices)
				devices.PATCH("/:id/trust", handler.MarkDeviceAsTrusted)
			}

			// Data Exports routes
			exports := protected.Group("/exports")
			{
				exports.POST("", handler.RequestDataExport)
				exports.GET("", handler.GetDataExports)
				exports.GET("/:id", handler.GetDataExport)
				exports.GET("/:id/download", handler.DownloadDataExport)
				exports.DELETE("/:id", handler.DeleteDataExport)
			}

			oauth := protected.Group("/oauth")
			{
				oauth.POST("/link", handler.LinkOAuthAccount)
				oauth.GET("/accounts", handler.GetOAuthAccounts)
				oauth.DELETE("/unlink/:provider", handler.UnlinkOAuthAccount)
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
