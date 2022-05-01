package dex

import (
	"github.com/deifyed/xctl/pkg/cloud"
)

const logFeature = "plugin/dex"

type reconciler struct {
	cloudProvider cloud.Provider
}
