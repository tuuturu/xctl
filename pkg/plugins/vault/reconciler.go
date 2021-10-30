package vault

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/clients/helm"
	"github.com/deifyed/xctl/pkg/cloud"

	"github.com/deifyed/xctl/pkg/controller/common/reconciliation"
)

func (v vaultReconciler) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	kubeConfigPath, err := config.GetAbsoluteKubeconfigPath(rctx.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("acquiring KubeConfig path: %w", err)
	}

	clients, err := prepareClients(rctx.Filesystem, kubeConfigPath)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("preparing clients: %w", err)
	}

	plugin := NewVaultPlugin()

	action, err := v.determineAction(rctx, clients.helm, plugin)
	if err != nil {
		return reconciliation.Result{Requeue: false}, fmt.Errorf("determining course of action: %w", err)
	}

	// action = "delete"

	switch action {
	case reconciliation.ActionCreate:
		logrus.Debug("installing Vault")

		err = clients.helm.Install(plugin)
		if err != nil {
			return reconciliation.Result{Requeue: false}, fmt.Errorf("installing vault: %w", err)
		}

		logrus.Debug("Initializing Vault")

		err = initializeVault(clients.kubectl, clients.vault)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("initializing vault: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err = clients.helm.Delete(plugin)
		if err != nil {
			return reconciliation.Result{Requeue: false}, fmt.Errorf("uninstalling vault: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{Requeue: false}, reconciliation.ErrIndecisive
}

func (v vaultReconciler) determineAction(rctx reconciliation.Context, helmClient helm.Client, plugin v1alpha1.Plugin) (
	reconciliation.Action, error,
) {
	indication := reconciliation.DetermineUserIndication(rctx, true)

	clusterExists, err := v.cloudProvider.HasCluster(rctx.Ctx, rctx.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf("checking cluster existence: %w", err)
	}

	vaultExists := false

	if clusterExists {
		vaultExists, err = helmClient.Exists(plugin)
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("checking vault existence: %w", err)
		}
	}

	switch indication {
	case reconciliation.ActionCreate:
		if !clusterExists {
			return reconciliation.ActionWait, nil
		}

		if vaultExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !vaultExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

func (v vaultReconciler) String() string {
	return "Vault"
}

func NewVaultReconciler(cloudProvider cloud.Provider) reconciliation.Reconciler {
	return &vaultReconciler{
		cloudProvider: cloudProvider,
	}
}
