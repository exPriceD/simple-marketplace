package token

import (
	"errors"
	"time"

	"github.com/exPriceD/simple-marketplace/internal/application/port"
	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTService(secret string, ttl time.Duration) port.TokenService {
	return &jwtService{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

type claims struct {
	UID   string `json:"uid"`
	Login string `json:"login"`
	jwt.RegisteredClaims
}

func (s *jwtService) Generate(userID, login string, ttl time.Duration) (string, error) {
	now := time.Now().UTC()
	d := ttl
	if d <= 0 {
		d = s.ttl
	}
	c := claims{
		UID:   userID,
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(d)),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return tok.SignedString(s.secret)
}

func (s *jwtService) Parse(tokenStr string) (port.Claims, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected_signing_method")
		}
		return s.secret, nil
	})
	if err != nil {
		return port.Claims{}, err
	}
	c, ok := tok.Claims.(*claims)
	if !ok || !tok.Valid {
		return port.Claims{}, errors.New("invalid_token")
	}
	return port.Claims{
		UserID: c.UID,
		Login:  c.Login,
		Exp:    c.ExpiresAt.Time,
		Issued: c.IssuedAt.Time,
	}, nil
}
