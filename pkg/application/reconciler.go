package application

import (
	"fmt"
	"path"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/tools/reconciliation"
)

type reconciler struct{}

func (r reconciler) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	action := reconciliation.ActionCreate
	if rctx.Purge {
		action = reconciliation.ActionDelete
	}

	applicationBaseDir := path.Join(
		applicationsDir(rctx.RootDirectory, rctx.ApplicationDeclaration.Metadata.Name),
		config.DefaultApplicationBaseDir,
	)

	switch action {
	case reconciliation.ActionCreate:
		err := writeBaseManifests(rctx.Filesystem, applicationBaseDir, rctx.ApplicationDeclaration)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("writing base manifests: %w", err)
		}

		return reconciliation.Result{}, nil
	case reconciliation.ActionDelete:
		return reconciliation.Result{}, nil
	}

	return reconciliation.Result{}, reconciliation.ErrIndecisive
}

func applicationsDir(rootDir string, appName string) string {
	return path.Join(
		rootDir,
		config.DefaultInfrastructureDir,
		config.DefaultApplicationsDir,
		appName,
	)
}

func writeBaseManifests(fs readerWriter, targetDir string, application v1alpha1.Application) error {
	deployment, err := scaffoldDeployment(application)
	if err != nil {
		return fmt.Errorf("scaffolding deployment: %w", err)
	}

	err = fs.WriteReader(path.Join(targetDir, "deployment.yaml"), deployment)
	if err != nil {
		return fmt.Errorf("writing deployment: %w", err)
	}

	service, err := scaffoldService(application)
	if err != nil {
		return fmt.Errorf("scaffolding service: %w", err)
	}

	err = fs.WriteReader(path.Join(targetDir, "service.yaml"), service)
	if err != nil {
		return fmt.Errorf("writing service: %w", err)
	}

	ingress, err := scaffoldIngress(application)
	if err != nil {
		return fmt.Errorf("scaffolding ingress: %w", err)
	}

	err = fs.WriteReader(path.Join(targetDir, "ingress.yaml"), ingress)
	if err != nil {
		return fmt.Errorf("writing ingress: %w", err)
	}

	return nil
}

func (r reconciler) String() string {
	return "manifests"
}
