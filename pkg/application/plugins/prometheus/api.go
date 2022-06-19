package prometheus

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/deifyed/xctl/pkg/config"
	"github.com/deifyed/xctl/pkg/tools/kustomize"

	"github.com/deifyed/xctl/pkg/tools/reconciliation"
)

type reconciler struct {
	absoluteApplicationDirectory string
}

const defaultServiceMonitorFileName = "service-monitor.yaml"

// Reconciler returns an initialized namespace reconciler
func Reconciler(absoluteApplicationDirectory string) reconciliation.Reconciler {
	return &reconciler{
		absoluteApplicationDirectory: absoluteApplicationDirectory,
	}
}

func (r reconciler) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	action := reconciliation.DetermineUserIndication(rctx, rctx.ApplicationDeclaration.Spec.Metrics != "")

	serviceMonitorPath := path.Join(
		r.absoluteApplicationDirectory,
		config.DefaultApplicationBaseDir,
		defaultServiceMonitorFileName,
	)

	switch action {
	case reconciliation.ActionCreate:
		serviceMonitor, err := scaffoldServiceMonitor(
			rctx.ApplicationDeclaration.Metadata.Name,
			rctx.ApplicationDeclaration.Spec.Metrics,
		)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("scaffolding: %w", err)
		}

		err = rctx.Filesystem.WriteReader(serviceMonitorPath, serviceMonitor)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("writing: %w", err)
		}

		err = kustomize.AddResourceToKustomization(
			rctx.Filesystem,
			path.Join(r.absoluteApplicationDirectory, config.DefaultApplicationBaseDir),
			defaultServiceMonitorFileName,
		)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("adding kustomization entry: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err := rctx.Filesystem.Remove(serviceMonitorPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return reconciliation.Result{Requeue: false}, err
			}

			return reconciliation.Result{}, fmt.Errorf("deleting service monitor: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{}, nil
}

func (r reconciler) String() string {
	return "Prometheus"
}
