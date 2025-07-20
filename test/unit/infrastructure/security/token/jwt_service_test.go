package token_test

import (
	"testing"
	"time"

	"github.com/exPriceD/simple-marketplace/internal/infrastructure/security/token"
)

func TestJWTService_GenerateParse(t *testing.T) {
	s := token.NewJWTService("secret", time.Minute)
	tok, err := s.Generate("uid123", "loginX", 0)
	if err != nil {
		t.Fatalf("generate err: %v", err)
	}
	claims, err := s.Parse(tok)
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	if claims.UserID != "uid123" || claims.Login != "loginX" {
		t.Fatalf("claims mismatch %+v", claims)
	}
}

func TestJWTService_Expired(t *testing.T) {
	s := token.NewJWTService("secret", time.Millisecond*10)
	tok, _ := s.Generate("u", "l", 0)
	time.Sleep(time.Millisecond * 20)
	_, err := s.Parse(tok)
	if err == nil {
		t.Fatalf("expected expired error")
	}
}
