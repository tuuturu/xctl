package manifests

import "github.com/deifyed/xctl/pkg/tools/reconciliation"

// Reconciler returns an initialized manifests reconciler
func Reconciler(absoluteApplicationDir string) reconciliation.Reconciler {
	return &reconciler{
		absoluteApplicationDir: absoluteApplicationDir,
	}
}
