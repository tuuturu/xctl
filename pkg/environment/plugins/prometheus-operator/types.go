package prometheus_operator

import (
	"github.com/deifyed/xctl/pkg/cloud"
)

type reconciler struct {
	cloudProvider cloud.Provider
}
