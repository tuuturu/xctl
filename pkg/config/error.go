package config

import "errors"

var (
	// ErrNotFound indicates that something is missing
	ErrNotFound = errors.New("not found")
	// ErrTimeout indicates that something exceeded its deadline
	ErrTimeout = errors.New("timeout")
)
