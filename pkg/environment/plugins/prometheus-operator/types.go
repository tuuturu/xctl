package prometheus_operator

import (
	"github.com/deifyed/xctl/pkg/cloud"
)

const logFeature = "plugin/prometheus-operator"

type reconciler struct {
	cloudProvider cloud.Provider
}
