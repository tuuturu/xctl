package kubectl

import "github.com/pkg/errors"

// ErrConnectionRefused indicates problems connecting to the Kubernetes Cluster
var ErrConnectionRefused = errors.New("connection refused")
