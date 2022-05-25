package application

import (
	"fmt"
	"path"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/config"
	"github.com/deifyed/xctl/pkg/tools/reconciliation"
	"github.com/spf13/afero"
)

type reconciler struct{}

func (r reconciler) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	action := reconciliation.ActionCreate
	if rctx.Purge {
		action = reconciliation.ActionDelete
	}

	switch action {
	case reconciliation.ActionCreate:
		err := writeBaseManifests(rctx.Filesystem, rctx.RootDirectory, rctx.ApplicationDeclaration)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("writing base manifests: %w", err)
		}

		return reconciliation.Result{}, nil
	case reconciliation.ActionDelete:
		return reconciliation.Result{}, nil
	}

	return reconciliation.Result{}, reconciliation.ErrIndecisive
}

func writeBaseManifests(fs *afero.Afero, rootDir string, application v1alpha1.Application) error {
	applicationDir := path.Join(rootDir, config.DefaultApplicationsDir, application.Metadata.Name)
	baseDir := path.Join(applicationDir, config.DefaultApplicationBaseDir)

	deployment, err := scaffoldDeployment(application)
	if err != nil {
		return fmt.Errorf("scaffolding deployment: %w", err)
	}

	err = fs.WriteReader(path.Join(baseDir, "deployment.yaml"), deployment)
	if err != nil {
		return fmt.Errorf("writing deployment: %w", err)
	}

	service, err := scaffoldService(application)
	if err != nil {
		return fmt.Errorf("scaffolding service: %w", err)
	}

	err = fs.WriteReader(path.Join(baseDir, "service.yaml"), service)
	if err != nil {
		return fmt.Errorf("writing service: %w", err)
	}

	ingress, err := scaffoldIngress(application)
	if err != nil {
		return fmt.Errorf("scaffolding ingress: %w", err)
	}

	err = fs.WriteReader(path.Join(baseDir, "ingress.yaml"), ingress)
	if err != nil {
		return fmt.Errorf("writing ingress: %w", err)
	}

	return nil
}

func (r reconciler) String() string {
	return "manifests"
}
