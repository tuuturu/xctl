package argocd

import (
	"github.com/deifyed/xctl/pkg/cloud"
)

const logFeature = "plugin/argocd"

type reconciler struct {
	cloudProvider cloud.Provider
}
