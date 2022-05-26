package manifests

import (
	"fmt"
	"path"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/tools/reconciliation"
)

func (r reconciler) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	action := reconciliation.ActionCreate
	if rctx.Purge {
		action = reconciliation.ActionDelete
	}

	applicationBaseDir := path.Join(r.absoluteApplicationDir, config.DefaultApplicationBaseDir)
	applicationOverlaysDir := path.Join(
		r.absoluteApplicationDir,
		config.DefaultApplicationsOverlaysDir,
		rctx.EnvironmentManifest.Metadata.Name,
	)

	switch action {
	case reconciliation.ActionCreate:
		err := writeBaseManifests(rctx.Filesystem, applicationBaseDir, rctx.ApplicationDeclaration)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("writing base manifests: %w", err)
		}

		err = writeOverlaysPatches(rctx.Filesystem, applicationOverlaysDir, rctx.ApplicationDeclaration)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("writing overlays patches: %w", err)
		}

		return reconciliation.Result{}, nil
	case reconciliation.ActionDelete:
		return reconciliation.Result{}, nil
	}

	return reconciliation.Result{}, reconciliation.ErrIndecisive
}

func (r reconciler) String() string {
	return "manifests"
}
