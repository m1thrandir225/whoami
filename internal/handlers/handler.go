// Package handlers defines HTTP Handlers using Gin
package handlers

import (
	"github.com/m1thrandir225/whoami/internal/oauth"
	"github.com/m1thrandir225/whoami/internal/security"
	"github.com/m1thrandir225/whoami/internal/services"
	"github.com/m1thrandir225/whoami/internal/util"
)

type OAuthProviders struct {
	Google *oauth.GoogleProvider
	GitHub *oauth.GitHubProvider
}

type HTTPHandler struct {
	userService             services.UserService
	securityService         services.SecurityService
	passwordSecurityService services.PasswordSecurityService
	passwordResetService    services.PasswordResetService
	emailService            services.EmailService
	auditService            services.AuditService
	userDevicesService      services.UserDevicesService
	dataExportsService      services.DataExportsService
	oauthService            services.OAuthService
	oauthProviders          OAuthProviders
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
	auditService services.AuditService,
	userDevicesService services.UserDevicesService,
	dataExportsService services.DataExportsService,
	oauthService services.OAuthService,
	oauthProviders OAuthProviders,
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
		auditService:            auditService,
		userDevicesService:      userDevicesService,
		dataExportsService:      dataExportsService,
		oauthService:            oauthService,
		oauthProviders:          oauthProviders,
		tokenMaker:              tokenMaker,
		tokenBlacklist:          tokenBlacklist,
		sessionService:          sessionService,
		config:                  &config,
		rateLimiter:             rateLimiter,
	}
}
