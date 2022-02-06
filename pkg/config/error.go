package config

import "errors"

// ErrTimeout indicates that something exceeded its deadline
var ErrTimeout = errors.New("timeout")
