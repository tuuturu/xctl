package promtail

import (
	"github.com/deifyed/xctl/pkg/cloud"
)

const logFeature = "plugin/promtail"

type reconciler struct {
	cloudProvider cloud.Provider
}
