package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *HTTPHandler) Register(ctx *gin.Context) {}

func (h *HTTPHandler) Login(ctx *gin.Context) {}

func (h *HTTPHandler) GetCurrentUser(ctx *gin.Context) {}

func (h *HTTPHandler) DeactivateUser(ctx *gin.Context) {}

func (h *HTTPHandler) ActivateUser(ctx *gin.Context) {}

func (h *HTTPHandler) UpdateUser(ctx *gin.Context) {}

func (h *HTTPHandler) UpdateUserPrivacySettings(ctx *gin.Context) {}
