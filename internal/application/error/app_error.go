package apperror

import (
	"errors"
	"fmt"
)

type Kind string

const (
	KindValidation Kind = "validation"
	KindConflict   Kind = "conflict"
	KindAuth       Kind = "auth"
	KindNotFound   Kind = "not_found"
	KindInternal   Kind = "internal"
)

type AppError struct {
	Op   string
	Code string
	Kind Kind
	Err  error
	Meta map[string]any
}

func (e *AppError) Error() string {
	if e.Op != "" {
		return fmt.Sprintf("%s: %s", e.Op, e.Code)
	}
	return e.Code
}

func (e *AppError) Unwrap() error { return e.Err }

func New(op, code string, kind Kind, err error) *AppError {
	return &AppError{Op: op, Code: code, Kind: kind, Err: err}
}

func IsKind(err error, k Kind) bool {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae.Kind == k
	}
	return false
}
