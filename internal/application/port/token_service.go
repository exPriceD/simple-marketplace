package port

import "time"

type TokenService interface {
	Generate(userID, login string, ttl time.Duration) (string, error)
	Parse(token string) (Claims, error)
}

type Claims struct {
	UserID string
	Login  string
	Exp    time.Time
	Issued time.Time
}
