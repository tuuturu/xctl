package config

import "errors"

var (
	// ErrNotFound indicates that something is missing
	ErrNotFound = errors.New("errNotFound")
	// ErrTimeout indicates that something exceeded its deadline
	ErrTimeout = errors.New("errTimeout")
	// ErrNotAuthenticated indicates invalid authentication
	ErrNotAuthenticated = errors.New("errNotAuthenticated")
)
