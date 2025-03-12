// internal/domain/errors.go
package domain

import "errors"

// Predefined errors for domain-related issues.
var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("resource not found")
	ErrDBFailure    = errors.New("database error")
)
