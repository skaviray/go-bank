package token

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKey = 32

type JWTMaker struct {
	secretKey string
}

func NewJwtMaker(sectetKey string) (Maker, error) {
	if len(sectetKey) < minSecretKey {
		return nil, fmt.Errorf("invalid key size, must be atleast %d charecters", minSecretKey)
	}
	maker := JWTMaker{
		secretKey: sectetKey,
	}
	return &maker, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	jwtClaim := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
		ExpiresAt: jwt.NewNumericDate(payload.ExpiresAt),
		Audience:  []string{payload.Username},
		ID:        payload.Id.String(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim)
	return jwtToken.SignedString([]byte(maker.secretKey))

}
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}
	jwtClaim, ok := jwtToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, ErrInvalidToken
	}
	id, err := uuid.Parse(jwtClaim.ID)
	if err != nil {
		return nil, ErrInvalidToken
	}
	payload := &Payload{
		Id:        id,
		ExpiresAt: jwtClaim.ExpiresAt.Time,
		IssuedAt:  jwtClaim.IssuedAt.Time,
		Username:  strings.Join(jwtClaim.Audience, ","),
	}
	return payload, nil
}
