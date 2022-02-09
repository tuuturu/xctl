package prometheus

import (
	"github.com/deifyed/xctl/pkg/cloud"
)

const logFeature = "plugin/prometheus"

type reconciler struct {
	cloudProvider cloud.Provider
}
