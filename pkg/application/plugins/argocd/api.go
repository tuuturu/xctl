package argocd

import (
	_ "embed"
	"fmt"
	"path"

	"github.com/deifyed/xctl/pkg/tools/reconciliation"
)

func (r reconciler) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	action := reconciliation.ActionCreate
	if rctx.Purge || !rctx.EnvironmentManifest.Spec.Plugins.ArgoCD {
		action = reconciliation.ActionDelete
	}

	applicationPath := path.Join(
		argoCDApplicationPath(r.absoluteEnvironmentDirectory),
		fmt.Sprintf("%s.yaml", rctx.ApplicationDeclaration.Metadata.Name),
	)

	switch action {
	case reconciliation.ActionCreate:
		argoCDApplication, err := scaffoldArgoCDApplication(rctx.EnvironmentManifest, rctx.ApplicationDeclaration)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("scaffolding: %w", err)
		}

		err = rctx.Filesystem.WriteReader(applicationPath, argoCDApplication)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("writing: %w", err)
		}

		return reconciliation.Result{}, nil
	case reconciliation.ActionDelete:
		err := rctx.Filesystem.RemoveAll(applicationPath)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("removing: %w", err)
		}

		return reconciliation.Result{}, nil
	}

	return reconciliation.Result{}, reconciliation.ErrIndecisive
}

func (r reconciler) String() string {
	return "ArgoCD application"
}

// Reconciler returns an initialized ArgoCD application reconciler
func Reconciler(absoluteEnvironmentDirectory string) reconciliation.Reconciler {
	return &reconciler{absoluteEnvironmentDirectory: absoluteEnvironmentDirectory}
}
