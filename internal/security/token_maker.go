// Package security
package security

import (
	"time"
)

type TokenMaker interface {
	CreateToken(userID int64, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
