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
	emailService            services.EmailService
	tokenMaker              security.TokenMaker
	config                  *util.Config
	rateLimiter             *security.RateLimiter
}

func NewHTTPHandler(
	userService services.UserService,
	securityService services.SecurityService,
	passwordSecurityService services.PasswordSecurityService,
	emailService services.EmailService,
	tokenMaker security.TokenMaker,
	rateLimiter *security.RateLimiter,
	config util.Config,
) *HTTPHandler {
	return &HTTPHandler{
		userService:             userService,
		securityService:         securityService,
		passwordSecurityService: passwordSecurityService,
		emailService:            emailService,
		tokenMaker:              tokenMaker,
		config:                  &config,
		rateLimiter:             rateLimiter,
	}
}
