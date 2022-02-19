package helm

import "errors"

// ErrUnreachable indicates something being unreachable. The cluster f.eks
var ErrUnreachable = errors.New("unreachable")

// ErrTimeout indicates something timed out.
var ErrTimeout = errors.New("timeout")

// ErrAlreadyExists indicates failure when installing something that already exists
var ErrAlreadyExists = errors.New("already exists")
