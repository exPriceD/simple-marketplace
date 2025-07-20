package shared

import "errors"

// Базовые доменные ошибки
var (
	ErrInvalidLogin       = errors.New("invalid_login_format")
	ErrInvalidPassword    = errors.New("invalid_password_format")
	ErrInvalidTitle       = errors.New("invalid_title")
	ErrInvalidDescription = errors.New("invalid_description")
	ErrInvalidPrice       = errors.New("invalid_price")
	ErrInvalidImageURL    = errors.New("invalid_image_url")
)
