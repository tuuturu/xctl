package ingress

import (
	"errors"
	"fmt"

	reconciliation2 "github.com/deifyed/xctl/pkg/tools/reconciliation"

	"github.com/deifyed/xctl/pkg/tools/clients/helm"
	helmBinary "github.com/deifyed/xctl/pkg/tools/clients/helm/binary"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/cloud"
)

//nolint:funlen
func (n nginxIngressController) Reconcile(rctx reconciliation2.Context) (reconciliation2.Result, error) {
	log := logging.GetLogger(logFeature, "reconciliation")

	kubeConfigPath, err := config.GetAbsoluteKubeconfigPath(rctx.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation2.Result{}, fmt.Errorf("acquiring KubeConfig path: %w", err)
	}

	helmClient, err := helmBinary.New(rctx.Filesystem, kubeConfigPath)
	if err != nil {
		return reconciliation2.Result{}, fmt.Errorf("acquiring Helm client: %w", err)
	}

	plugin := NewNginxIngressControllerPlugin()

	action, err := n.determineAction(determineActionOpts{
		Ctx:    rctx,
		Helm:   helmClient,
		Plugin: plugin,
		Logger: log,
	})
	if err != nil {
		return reconciliation2.Result{Requeue: false}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation2.ActionCreate:
		log.Debug("installing")

		err = helmClient.Install(plugin)
		if err != nil {
			switch {
			case errors.Is(err, helm.ErrTimeout):
				log.Info("Requeuing due to timeout")

				return reconciliation2.Result{Requeue: true}, nil
			case errors.Is(err, helm.ErrUnreachable):
				log.Info("Requeuing due to cluster being unreachable")

				return reconciliation2.Result{Requeue: true}, nil
			default:
				return reconciliation2.Result{Requeue: false}, fmt.Errorf("installing helm chart: %w", err)
			}
		}

		return reconciliation2.Result{Requeue: false}, nil
	case reconciliation2.ActionDelete:
		log.Debug("deleting")

		err = helmClient.Delete(plugin)
		if err != nil {
			return reconciliation2.Result{Requeue: false}, fmt.Errorf("uninstalling helm chart: %w", err)
		}

		return reconciliation2.Result{Requeue: false}, nil
	}

	return reconciliation2.NoopWaitIndecisiveHandler(action)
}

func (n nginxIngressController) determineAction(opts determineActionOpts) (reconciliation2.Action, error) {
	log := opts.Logger

	indication := reconciliation2.DetermineUserIndication(
		opts.Ctx,
		opts.Ctx.ClusterDeclaration.Spec.Plugins.NginxIngressController,
	)

	clusterExists, err := n.cloudProvider.HasCluster(opts.Ctx.Ctx, opts.Ctx.ClusterDeclaration)
	if err != nil {
		return reconciliation2.ActionNoop, fmt.Errorf("checking cluster existence: %w", err)
	}

	ingressExists := false

	if clusterExists {
		ingressExists, err = opts.Helm.Exists(opts.Plugin)
		if err != nil {
			return reconciliation2.ActionNoop, fmt.Errorf("checking component existence: %w", err)
		}
	}

	switch indication {
	case reconciliation2.ActionCreate:
		if !clusterExists {
			log.Debug("Waiting due to cluster not ready")

			return reconciliation2.ActionWait, nil
		}

		if ingressExists {
			log.Debug("NOOP since component already exists")

			return reconciliation2.ActionNoop, nil
		}

		return reconciliation2.ActionCreate, nil
	case reconciliation2.ActionDelete:
		if !ingressExists {
			log.Debug("NOOP since cluster is already taken down")

			return reconciliation2.ActionNoop, nil
		}

		return reconciliation2.ActionDelete, nil
	}

	return reconciliation2.ActionNoop, reconciliation2.ErrIndecisive
}

func (n nginxIngressController) String() string {
	return "Nginx Ingress Controller"
}

func NewReconciler(cloudProvider cloud.Provider) reconciliation2.Reconciler {
	return &nginxIngressController{
		cloudProvider: cloudProvider,
	}
}
