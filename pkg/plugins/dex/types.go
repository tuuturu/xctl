package dex

import (
	"github.com/deifyed/xctl/pkg/cloud"
)

const (
	logFeature = "plugin/dex"
	localURL   = "http://dex.operations:5556/dex"
)

type reconciler struct {
	cloudProvider cloud.Provider
}
