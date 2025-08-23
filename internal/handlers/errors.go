package handlers

import "errors"

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrNotFound            = errors.New("not found")
	ErrInternalServer      = errors.New("internal server error")
	ErrBadRequest          = errors.New("bad request")
	ErrUnprocessableEntity = errors.New("unprocessable entity")
	ErrConflict            = errors.New("conflict")
	ErrTooManyRequests     = errors.New("too many requests")
	ErrNotImplemented      = errors.New("not implemented")
	ErrInvalidToken        = errors.New("invalid token")
	ErrExpiredToken        = errors.New("expired token")
)
