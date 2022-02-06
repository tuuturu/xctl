package kubectl

import "github.com/pkg/errors"

var (
	// ErrConnectionRefused indicates problems connecting to the Kubernetes Cluster
	ErrConnectionRefused = errors.New("connection refused")
	// ErrNotFound indicates something is missing
	ErrNotFound = errors.New("not found")
)
