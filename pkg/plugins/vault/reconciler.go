package vault

import (
	"errors"
	"fmt"

	"github.com/deifyed/xctl/pkg/tools/clients/helm"
	kubectl2 "github.com/deifyed/xctl/pkg/tools/clients/kubectl"
	"github.com/deifyed/xctl/pkg/tools/clients/vault"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/cloud"

	"github.com/deifyed/xctl/pkg/controller/common/reconciliation"
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
		kubectl:    clients.kubectl,
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
		clusterExists    = true
		vaultExists      = true
		vaultInitialized = true
	)

	cluster, err := v.cloudProvider.GetCluster(opts.rctx.Ctx, opts.rctx.ClusterDeclaration)
	if err != nil {
		if !errors.Is(err, cloud.ErrNotFound) {
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

		vaultInitialized, err = opts.kubectl.PodReady(kubectl2.Pod{
			Name:      "vault-0",
			Namespace: opts.plugin.Metadata.Namespace,
		})
		if err != nil {
			if !errors.Is(err, kubectl2.ErrNotFound) {
				return "", fmt.Errorf("checking pod ready status: %w", err)
			}

			vaultInitialized = false
		}
	}

	switch indication {
	case reconciliation.ActionCreate:
		if !clusterExists || !cluster.Ready {
			return reconciliation.ActionWait, nil
		}

		if vaultExists && vaultInitialized {
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

func NewReconciler(cloudProvider cloud.Provider) reconciliation.Reconciler {
	return &vaultReconciler{
		cloudProvider: cloudProvider,
	}
}
