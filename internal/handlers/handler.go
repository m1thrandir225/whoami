// Package handlers defines HTTP Handlers using Gin
package handlers

import (
	"github.com/m1thrandir225/whoami/internal/security"
	"github.com/m1thrandir225/whoami/internal/services"
	"github.com/m1thrandir225/whoami/internal/util"
)

type HTTPHandler struct {
	userService             services.UserService
	securityService         services.SecurityService
	passwordSecurityService services.PasswordSecurityService
	passwordResetService    services.PasswordResetService
	emailService            services.EmailService
	tokenMaker              security.TokenMaker
	tokenBlacklist          security.TokenBlacklist
	sessionService          services.SessionService
	config                  *util.Config
	rateLimiter             *security.RateLimiter
}

func NewHTTPHandler(
	userService services.UserService,
	securityService services.SecurityService,
	passwordSecurityService services.PasswordSecurityService,
	passwordResetService services.PasswordResetService,
	emailService services.EmailService,
	tokenMaker security.TokenMaker,
	tokenBlacklist security.TokenBlacklist,
	sessionService services.SessionService,
	rateLimiter *security.RateLimiter,
	config util.Config,
) *HTTPHandler {
	return &HTTPHandler{
		userService:             userService,
		securityService:         securityService,
		passwordSecurityService: passwordSecurityService,
		passwordResetService:    passwordResetService,
		emailService:            emailService,
		tokenMaker:              tokenMaker,
		tokenBlacklist:          tokenBlacklist,
		sessionService:          sessionService,
		config:                  &config,
		rateLimiter:             rateLimiter,
	}
}
