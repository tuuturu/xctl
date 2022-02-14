package prometheus

import (
	"errors"
	"fmt"

	"github.com/deifyed/xctl/pkg/tools/reconciliation"

	"github.com/deifyed/xctl/pkg/tools/clients/helm"
	"github.com/deifyed/xctl/pkg/tools/clients/helm/binary"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/deifyed/xctl/pkg/cloud"
)

func (r reconciler) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	log := logging.GetLogger(logFeature, "reconciliation")

	kubeConfigPath, err := config.GetAbsoluteKubeconfigPath(rctx.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("acquiring kube config path: %w", err)
	}

	helmClient, err := binary.New(rctx.Filesystem, kubeConfigPath)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("acquiring Helm client: %w", err)
	}

	plugin := NewPlugin()

	action, err := r.determineAction(rctx, helmClient, plugin)
	if err != nil {
		return reconciliation.Result{Requeue: false}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation.ActionCreate:
		log.Debug("installing")

		err = helmClient.Install(plugin)
		if err != nil {
			if errors.Is(err, helm.ErrUnreachable) {
				return reconciliation.Result{Requeue: true}, nil
			}

			return reconciliation.Result{}, fmt.Errorf("running helm install: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		log.Debug("deleting")

		err = helmClient.Delete(plugin)
		if err != nil {
			if errors.Is(err, helm.ErrUnreachable) {
				return reconciliation.Result{Requeue: true}, nil
			}

			return reconciliation.Result{}, fmt.Errorf("running helm delete: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.NoopWaitIndecisiveHandler(action)
}

func (r reconciler) determineAction(rctx reconciliation.Context, helm helm.Client, plugin v1alpha1.Plugin) (reconciliation.Action, error) { //nolint:lll
	indication := reconciliation.DetermineUserIndication(rctx, rctx.ClusterDeclaration.Spec.Plugins.Prometheus)

	var (
		clusterExists   = true
		componentExists = true
	)

	_, err := r.cloudProvider.GetCluster(rctx.Ctx, rctx.ClusterDeclaration)
	if err != nil {
		if !errors.Is(err, cloud.ErrNotFound) {
			return "", fmt.Errorf("acquiring cluster: %w", err)
		}

		clusterExists = false
	}

	if clusterExists {
		componentExists, err = helm.Exists(plugin)
		if err != nil {
			return "", fmt.Errorf("checking component existence: %w", err)
		}
	}

	switch indication {
	case reconciliation.ActionCreate:
		if !clusterExists {
			return reconciliation.ActionWait, nil
		}

		if componentExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !clusterExists {
			return reconciliation.ActionNoop, nil
		}

		if !componentExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

func (r reconciler) String() string {
	return "Prometheus"
}

func NewReconciler(cloudProvider cloud.Provider) reconciliation.Reconciler {
	return &reconciler{
		cloudProvider: cloudProvider,
	}
}
