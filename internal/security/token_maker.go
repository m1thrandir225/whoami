// Package security
package security

import (
	"time"

	"github.com/google/uuid"
)

type TokenMaker interface {
	CreateToken(userID uuid.UUID, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
