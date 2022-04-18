package loki

import (
	"github.com/deifyed/xctl/pkg/cloud"
)

const logFeature = "plugin/loki"

type reconciler struct {
	cloudProvider cloud.Provider
}
