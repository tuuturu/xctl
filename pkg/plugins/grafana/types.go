package grafana

import (
	"github.com/deifyed/xctl/pkg/cloud"
)

const logFeature = "plugin/grafana"

type reconciler struct {
	cloudProvider cloud.Provider
}
