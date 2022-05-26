package namespace

import (
	_ "embed"
	"fmt"
	"path"

	"github.com/deifyed/xctl/pkg/tools/reconciliation"
)

type reconciler struct {
	absoluteEnvironmentDirectory string
}

// Reconciler returns an initialized namespace reconciler
func Reconciler(absoluteEnvironmentDirectory string) reconciliation.Reconciler {
	return &reconciler{
		absoluteEnvironmentDirectory: absoluteEnvironmentDirectory,
	}
}

func (r reconciler) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	namespace, err := scaffoldNamespace(rctx.ApplicationDeclaration.Metadata.Namespace)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("scaffolding: %w", err)
	}

	err = rctx.Filesystem.WriteReader(
		path.Join(
			r.absoluteEnvironmentDirectory,
			"argocd",
			"namespaces",
			fmt.Sprintf("%s.yaml", rctx.ApplicationDeclaration.Metadata.Namespace),
		),
		namespace,
	)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("writing: %w", err)
	}

	return reconciliation.Result{}, nil
}

func (r reconciler) String() string {
	return "Namespace"
}
