package prometheus_operator

import (
	"fmt"
	"strings"

	"github.com/deifyed/xctl/pkg/environment/plugins/prometheus"

	"github.com/deifyed/xctl/pkg/config"
	helmBinary "github.com/deifyed/xctl/pkg/tools/clients/helm/binary"

	"github.com/deifyed/xctl/pkg/tools/reconciliation"

	"github.com/deifyed/xctl/pkg/tools/clients/helm"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/deifyed/xctl/pkg/cloud"
)

func (r reconciler) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	log := logging.GetLogger(logFeature, "reconciliation")

	kubeconfigPath, err := config.GetAbsoluteKubeconfigPath(rctx.EnvironmentManifest.Metadata.Name)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("acquiring kubeconfig path: %w", err)
	}

	helmClient, err := helmBinary.New(rctx.Filesystem, kubeconfigPath)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("preparing helm client: %w", err)
	}

	action, err := r.determineAction(rctx, helmClient)
	if err != nil {
		return reconciliation.Result{Requeue: false}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation.ActionCreate:
		log.Debug("installing")

		err = helmClient.Install(plugin())
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("installing: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		log.Debug("deleting")

		err = helmClient.Delete(plugin())
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("uninstalling: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.NoopWaitIndecisiveHandler(action)
}

//nolint:funlen
func (r reconciler) determineAction(rctx reconciliation.Context, helm helm.Client) (reconciliation.Action, error) { //nolint:lll
	indication := reconciliation.DetermineUserIndication(rctx, rctx.EnvironmentManifest.Spec.Plugins.Prometheus)

	var (
		err              error
		clusterExists    bool
		componentExists  bool
		prometheusExists bool
	)

	clusterExists, err = r.cloudProvider.HasCluster(rctx.Ctx, rctx.EnvironmentManifest)
	if err != nil {
		return "", fmt.Errorf("acquiring cluster: %w", err)
	}

	if clusterExists {
		prometheusExists, err = helm.Exists(prometheus.NewPlugin())
		if err != nil {
			return "", fmt.Errorf("checking Prometheus existence: %w", err)
		}

		componentExists, err = helm.Exists(plugin())
		if err != nil {
			return "", fmt.Errorf("checking Prometheus Operator existence: %w", err)
		}
	}

	switch indication {
	case reconciliation.ActionCreate:
		if !clusterExists {
			return reconciliation.ActionWait, nil
		}

		if !prometheusExists {
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
	return strings.Title(plugin().Metadata.Name)
}

func Reconciler(cloudProvider cloud.Provider) reconciliation.Reconciler {
	return &reconciler{
		cloudProvider: cloudProvider,
	}
}
