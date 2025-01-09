package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token has invalid claims: token is expired")
)

type Payload struct {
	Id        uuid.UUID `json:"id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Username  string    `json:"username"`
}

// type Payload struct {
// 	Id uuid.UUID `json:"id"`
// 	jwt.RegisteredClaims
// }

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Payload{
		Id:        uuid,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
		Username:  username,
	}, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrTokenExpired
	}
	return nil
}
