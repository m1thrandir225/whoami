// Package handlers defines HTTP Handlers using Gin
package handlers

import (
	"github.com/m1thrandir225/whoami/internal/security"
	"github.com/m1thrandir225/whoami/internal/services"
	"github.com/m1thrandir225/whoami/internal/util"
)

type HTTPHandler struct {
	userService services.UserService
	tokenMaker  security.TokenMaker
	config      *util.Config
}

func NewHTTPHandler(
	userService services.UserService,
	tokenMaker security.TokenMaker,
	config util.Config,
) *HTTPHandler {
	return &HTTPHandler{
		userService: userService,
		tokenMaker:  tokenMaker,
		config:      &config,
	}
}
