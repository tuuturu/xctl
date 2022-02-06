package vault

import (
	"errors"
	"fmt"

	"github.com/deifyed/xctl/pkg/clients/vault"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/cloud"

	"github.com/deifyed/xctl/pkg/controller/common/reconciliation"

	"github.com/deifyed/xctl/pkg/clients/helm"
)

func (v vaultReconciler) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	log := logging.GetLogger(logFeature, "reconciliation")

	kubeConfigPath, err := config.GetAbsoluteKubeconfigPath(rctx.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("acquiring KubeConfig path: %w", err)
	}

	clients, err := prepareClients(rctx.Filesystem, kubeConfigPath)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("preparing clients: %w", err)
	}

	plugin := NewVaultPlugin()

	action, err := v.determineAction(determineActionOpts{
		rctx:       rctx,
		helmClient: clients.helm,
		plugin:     plugin,
	})
	if err != nil {
		return reconciliation.Result{Requeue: false}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation.ActionCreate:
		log.Debug("installing")

		err = installVault(clients)
		if err != nil {
			switch {
			case errors.Is(err, helm.ErrUnreachable):
				return reconciliation.Result{Requeue: true}, nil
			case errors.Is(err, vault.ErrConnectionRefused):
				return reconciliation.Result{Requeue: true}, nil
			default:
				return reconciliation.Result{}, fmt.Errorf("installing: %w", err)
			}
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		log.Debug("deleting")

		err = clients.helm.Delete(plugin)
		if err != nil {
			return reconciliation.Result{Requeue: false}, fmt.Errorf("uninstalling vault: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.NoopWaitIndecisiveHandler(action)
}

func (v vaultReconciler) determineAction(opts determineActionOpts) (reconciliation.Action, error) {
	indication := reconciliation.DetermineUserIndication(opts.rctx, opts.rctx.ClusterDeclaration.Spec.Plugins.Vault)

	var (
		clusterExists = true
		vaultExists   = true
	)

	cluster, err := v.cloudProvider.GetCluster(opts.rctx.Ctx, opts.rctx.ClusterDeclaration.Metadata.Name)
	if err != nil {
		if !errors.Is(err, config.ErrNotFound) {
			return "", fmt.Errorf("checking cluster existence: %w", err)
		}

		clusterExists = false
	}

	if clusterExists && cluster.Ready {
		vaultExists, err = opts.helmClient.Exists(opts.plugin)
		if err != nil {
			if errors.Is(err, helm.ErrUnreachable) {
				return reconciliation.ActionWait, nil
			}

			return reconciliation.ActionNoop, fmt.Errorf("checking vault existence: %w", err)
		}
	}

	switch indication {
	case reconciliation.ActionCreate:
		if !clusterExists || !cluster.Ready {
			return reconciliation.ActionWait, nil
		}

		if vaultExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !clusterExists || !vaultExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
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
