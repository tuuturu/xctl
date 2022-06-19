package prometheus

import (
	_ "embed"
	"fmt"
	"path"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/tools/reconciliation"
)

type reconciler struct {
	absoluteApplicationDirectory string
}

// Reconciler returns an initialized namespace reconciler
func Reconciler(absoluteApplicationDirectory string) reconciliation.Reconciler {
	return &reconciler{
		absoluteApplicationDirectory: absoluteApplicationDirectory,
	}
}

func (r reconciler) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	serviceMonitor, err := scaffoldServiceMonitor(
		rctx.ApplicationDeclaration.Metadata.Name,
		rctx.ApplicationDeclaration.Spec.Metrics,
	)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("scaffolding: %w", err)
	}

	err = rctx.Filesystem.WriteReader(
		path.Join(
			r.absoluteApplicationDirectory,
			config.DefaultApplicationBaseDir,
			"service-monitor.yaml",
		),
		serviceMonitor,
	)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("writing: %w", err)
	}

	return reconciliation.Result{}, nil
}

func (r reconciler) String() string {
	return "Prometheus"
}
