package vault

import (
	"errors"
	"fmt"

	reconciliation2 "github.com/deifyed/xctl/pkg/tools/reconciliation"

	"github.com/deifyed/xctl/pkg/tools/clients/helm"
	kubectl2 "github.com/deifyed/xctl/pkg/tools/clients/kubectl"
	"github.com/deifyed/xctl/pkg/tools/clients/vault"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/cloud"
)

func (v vaultReconciler) Reconcile(rctx reconciliation2.Context) (reconciliation2.Result, error) {
	log := logging.GetLogger(logFeature, "reconciliation")

	kubeConfigPath, err := config.GetAbsoluteKubeconfigPath(rctx.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation2.Result{}, fmt.Errorf("acquiring KubeConfig path: %w", err)
	}

	clients, err := prepareClients(rctx.Filesystem, kubeConfigPath)
	if err != nil {
		return reconciliation2.Result{}, fmt.Errorf("preparing clients: %w", err)
	}

	plugin := NewVaultPlugin()

	action, err := v.determineAction(determineActionOpts{
		rctx:       rctx,
		helmClient: clients.helm,
		kubectl:    clients.kubectl,
		plugin:     plugin,
	})
	if err != nil {
		return reconciliation2.Result{Requeue: false}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation2.ActionCreate:
		log.Debug("installing")

		err = installVault(clients)
		if err != nil {
			switch {
			case errors.Is(err, helm.ErrUnreachable):
				return reconciliation2.Result{Requeue: true}, nil
			case errors.Is(err, vault.ErrConnectionRefused):
				return reconciliation2.Result{Requeue: true}, nil
			default:
				return reconciliation2.Result{}, fmt.Errorf("installing: %w", err)
			}
		}

		return reconciliation2.Result{Requeue: false}, nil
	case reconciliation2.ActionDelete:
		log.Debug("deleting")

		err = clients.helm.Delete(plugin)
		if err != nil {
			return reconciliation2.Result{Requeue: false}, fmt.Errorf("uninstalling vault: %w", err)
		}

		return reconciliation2.Result{Requeue: false}, nil
	}

	return reconciliation2.NoopWaitIndecisiveHandler(action)
}

func (v vaultReconciler) determineAction(opts determineActionOpts) (reconciliation2.Action, error) {
	indication := reconciliation2.DetermineUserIndication(opts.rctx, opts.rctx.ClusterDeclaration.Spec.Plugins.Vault)

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
				return reconciliation2.ActionWait, nil
			}

			return reconciliation2.ActionNoop, fmt.Errorf("checking vault existence: %w", err)
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
	case reconciliation2.ActionCreate:
		if !clusterExists || !cluster.Ready {
			return reconciliation2.ActionWait, nil
		}

		if vaultExists && vaultInitialized {
			return reconciliation2.ActionNoop, nil
		}

		return reconciliation2.ActionCreate, nil
	case reconciliation2.ActionDelete:
		if !clusterExists || !vaultExists {
			return reconciliation2.ActionNoop, nil
		}

		return reconciliation2.ActionDelete, nil
	}

	return reconciliation2.ActionNoop, reconciliation2.ErrIndecisive
}

func (v vaultReconciler) String() string {
	return "Vault"
}

func NewReconciler(cloudProvider cloud.Provider) reconciliation2.Reconciler {
	return &vaultReconciler{
		cloudProvider: cloudProvider,
	}
}
