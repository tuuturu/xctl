package prometheus

import (
	"errors"
	"fmt"

	reconciliation2 "github.com/deifyed/xctl/pkg/tools/reconciliation"

	helm2 "github.com/deifyed/xctl/pkg/tools/clients/helm"
	"github.com/deifyed/xctl/pkg/tools/clients/helm/binary"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/deifyed/xctl/pkg/cloud"
)

func (r reconciler) Reconcile(rctx reconciliation2.Context) (reconciliation2.Result, error) {
	log := logging.GetLogger(logFeature, "reconciliation")

	kubeConfigPath, err := config.GetAbsoluteKubeconfigPath(rctx.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation2.Result{}, fmt.Errorf("acquiring kube config path: %w", err)
	}

	helmClient, err := binary.New(rctx.Filesystem, kubeConfigPath)
	if err != nil {
		return reconciliation2.Result{}, fmt.Errorf("acquiring Helm client: %w", err)
	}

	plugin := NewPlugin()

	action, err := r.determineAction(rctx, helmClient, plugin)
	if err != nil {
		return reconciliation2.Result{Requeue: false}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation2.ActionCreate:
		log.Debug("installing")

		err = helmClient.Install(plugin)
		if err != nil {
			if errors.Is(err, helm2.ErrUnreachable) {
				return reconciliation2.Result{Requeue: true}, nil
			}

			return reconciliation2.Result{}, fmt.Errorf("running helm install: %w", err)
		}

		return reconciliation2.Result{Requeue: false}, nil
	case reconciliation2.ActionDelete:
		log.Debug("deleting")

		err = helmClient.Delete(plugin)
		if err != nil {
			if errors.Is(err, helm2.ErrUnreachable) {
				return reconciliation2.Result{Requeue: true}, nil
			}

			return reconciliation2.Result{}, fmt.Errorf("running helm delete: %w", err)
		}

		return reconciliation2.Result{Requeue: false}, nil
	}

	return reconciliation2.NoopWaitIndecisiveHandler(action)
}

func (r reconciler) determineAction(rctx reconciliation2.Context, helm helm2.Client, plugin v1alpha1.Plugin) (reconciliation2.Action, error) { //nolint:lll
	indication := reconciliation2.DetermineUserIndication(rctx, rctx.ClusterDeclaration.Spec.Plugins.Prometheus)

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
	case reconciliation2.ActionCreate:
		if !clusterExists {
			return reconciliation2.ActionWait, nil
		}

		if componentExists {
			return reconciliation2.ActionNoop, nil
		}

		return reconciliation2.ActionCreate, nil
	case reconciliation2.ActionDelete:
		if !clusterExists {
			return reconciliation2.ActionNoop, nil
		}

		if !componentExists {
			return reconciliation2.ActionNoop, nil
		}

		return reconciliation2.ActionDelete, nil
	}

	return reconciliation2.ActionNoop, reconciliation2.ErrIndecisive
}

func (r reconciler) String() string {
	return "Prometheus"
}

func NewReconciler(cloudProvider cloud.Provider) reconciliation2.Reconciler {
	return &reconciler{
		cloudProvider: cloudProvider,
	}
}
