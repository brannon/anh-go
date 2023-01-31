package anh

import "github.com/pkg/errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotFound           = errors.New("resource not found")
)
