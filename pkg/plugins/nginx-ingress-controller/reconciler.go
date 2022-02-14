package ingress

import (
	"fmt"

	"github.com/deifyed/xctl/pkg/tools/reconciliation"

	helmBinary "github.com/deifyed/xctl/pkg/tools/clients/helm/binary"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/cloud"
)

//nolint:funlen
func (n nginxIngressController) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	log := logging.GetLogger(logFeature, "reconciliation")

	kubeConfigPath, err := config.GetAbsoluteKubeconfigPath(rctx.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("acquiring KubeConfig path: %w", err)
	}

	helmClient, err := helmBinary.New(rctx.Filesystem, kubeConfigPath)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("acquiring Helm client: %w", err)
	}

	plugin := NewNginxIngressControllerPlugin()

	action, err := n.determineAction(determineActionOpts{
		Ctx:    rctx,
		Helm:   helmClient,
		Plugin: plugin,
		Logger: log,
	})
	if err != nil {
		return reconciliation.Result{Requeue: false}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation.ActionCreate:
		log.Debug("installing")

		err = helmClient.Install(plugin)
		if err != nil {
			return reconciliation.Result{Requeue: false}, fmt.Errorf("installing helm chart: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		log.Debug("deleting")

		err = helmClient.Delete(plugin)
		if err != nil {
			return reconciliation.Result{Requeue: false}, fmt.Errorf("uninstalling helm chart: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.NoopWaitIndecisiveHandler(action)
}

func (n nginxIngressController) determineAction(opts determineActionOpts) (reconciliation.Action, error) {
	log := opts.Logger

	indication := reconciliation.DetermineUserIndication(
		opts.Ctx,
		opts.Ctx.ClusterDeclaration.Spec.Plugins.NginxIngressController,
	)

	clusterExists, err := n.cloudProvider.HasCluster(opts.Ctx.Ctx, opts.Ctx.ClusterDeclaration)
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf("checking cluster existence: %w", err)
	}

	ingressExists := false

	if clusterExists {
		ingressExists, err = opts.Helm.Exists(opts.Plugin)
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("checking component existence: %w", err)
		}
	}

	switch indication {
	case reconciliation.ActionCreate:
		if !clusterExists {
			log.Debug("Waiting due to cluster not ready")

			return reconciliation.ActionWait, nil
		}

		if ingressExists {
			log.Debug("NOOP since component already exists")

			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !ingressExists {
			log.Debug("NOOP since cluster is already taken down")

			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

func (n nginxIngressController) String() string {
	return "Nginx Ingress Controller"
}

func NewReconciler(cloudProvider cloud.Provider) reconciliation.Reconciler {
	return &nginxIngressController{
		cloudProvider: cloudProvider,
	}
}
