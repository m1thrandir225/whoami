// Package handlers defines HTTP Handlers using Gin
package handlers

import "github.com/m1thrandir225/whoami/internal/services"

type HTTPHandler struct {
	userService *services.UserService
}
