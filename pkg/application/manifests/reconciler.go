package manifests

import (
	"errors"
	"fmt"
	"path"

	"github.com/spf13/afero"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/tools/reconciliation"
)

func (r reconciler) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	action := reconciliation.DetermineUserIndication(rctx, true)

	applicationBaseDir := path.Join(r.absoluteApplicationDir, config.DefaultApplicationBaseDir)
	applicationOverlaysDir := path.Join(r.absoluteApplicationDir, config.DefaultApplicationsOverlaysDir)
	environmentSpecificOverlaysDir := path.Join(applicationOverlaysDir, rctx.EnvironmentManifest.Metadata.Name)

	switch action {
	case reconciliation.ActionCreate:
		err := writeBaseManifests(rctx.Filesystem, applicationBaseDir, rctx.ApplicationDeclaration)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("writing base manifests: %w", err)
		}

		err = writeOverlaysPatches(rctx.Filesystem, environmentSpecificOverlaysDir, rctx.ApplicationDeclaration)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("writing overlays patches: %w", err)
		}

		return reconciliation.Result{}, nil
	case reconciliation.ActionDelete:
		err := rctx.Filesystem.RemoveAll(environmentSpecificOverlaysDir)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("removing overlays directory: %w", err)
		}

		exists, err := rctx.Filesystem.Exists(applicationOverlaysDir)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("checking overlays directory existence: %w", err)
		}

		if !exists {
			return reconciliation.Result{}, nil
		}

		isEmpty, err := rctx.Filesystem.IsEmpty(applicationOverlaysDir)
		if err != nil && !errors.Is(err, afero.ErrFileNotFound) {
			return reconciliation.Result{}, fmt.Errorf("checking if overlays directory is empty: %w", err)
		}

		if !isEmpty {
			return reconciliation.Result{}, nil
		}

		err = rctx.Filesystem.RemoveAll(r.absoluteApplicationDir)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("removing application directory: %w", err)
		}

		return reconciliation.Result{}, nil
	}

	return reconciliation.Result{}, reconciliation.ErrIndecisive
}

func (r reconciler) String() string {
	return "manifests"
}
