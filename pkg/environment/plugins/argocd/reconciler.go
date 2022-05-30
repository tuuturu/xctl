package argocd

import (
	"fmt"
	"strings"

	kubectlBinary "github.com/deifyed/xctl/pkg/tools/clients/kubectl/binary"

	"github.com/deifyed/xctl/pkg/config"
	helmBinary "github.com/deifyed/xctl/pkg/tools/clients/helm/binary"

	"github.com/deifyed/xctl/pkg/tools/reconciliation"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/deifyed/xctl/pkg/cloud"
)

func (r reconciler) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	log := logging.GetLogger(logFeature, "reconciliation")

	repo := repository{URL: rctx.EnvironmentManifest.Spec.Repository}

	kubeconfigPath, err := config.GetAbsoluteKubeconfigPath(rctx.EnvironmentManifest.Metadata.Name)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("acquiring kubeconfig path: %w", err)
	}

	kubectlClient, err := kubectlBinary.New(rctx.Filesystem, kubeconfigPath)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("preparing kubectl client: %w", err)
	}

	helmClient, err := helmBinary.New(rctx.Filesystem, kubeconfigPath)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("preparing helm client: %w", err)
	}

	action, err := r.determineAction(rctx, helmClient)
	if err != nil {
		return reconciliation.Result{Requeue: false}, fmt.Errorf("determining course of action: %w", err)
	}

	plugin, err := NewPlugin()
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("creating plugin: %w", err)
	}

	switch action {
	case reconciliation.ActionCreate:
		log.Debug("installing")

		err = installArgoCD(rctx, helmClient, kubectlClient, repo)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("installing: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		log.Debug("deleting")

		err = helmClient.Delete(plugin)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("uninstalling: %w", err)
		}

		err = deleteKey(rctx.Ctx, rctx.Keyring, rctx.EnvironmentManifest.Metadata.Name, repo)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("deleting deploy key: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.NoopWaitIndecisiveHandler(action)
}

func (r reconciler) String() string {
	return strings.Title(pluginName)
}

func NewReconciler(cloudProvider cloud.Provider) reconciliation.Reconciler {
	return &reconciler{
		cloudProvider: cloudProvider,
	}
}
