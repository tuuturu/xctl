package ingress

import (
	"fmt"

	helmBinary "github.com/deifyed/xctl/pkg/clients/helm/binary"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/sirupsen/logrus"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/cloud"

	"github.com/deifyed/xctl/pkg/controller/common/reconciliation"
)

func (n nginxIngressController) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	log := logging.CreateEntry(logrus.StandardLogger(), logFeature, "reconciliation")

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
	})
	if err != nil {
		return reconciliation.Result{Requeue: false}, fmt.Errorf("determining course of action: %w", err)
	}

	// action = "create"

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
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{Requeue: false}, reconciliation.ErrIndecisive
}

func (n nginxIngressController) determineAction(opts determineActionOpts) (
	reconciliation.Action, error,
) {
	indication := reconciliation.DetermineUserIndication(opts.Ctx, true)

	clusterExists, err := n.cloudProvider.HasCluster(opts.Ctx.Ctx, opts.Ctx.ClusterDeclaration.Metadata.Name)
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
			return reconciliation.ActionWait, nil
		}

		if ingressExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !ingressExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

func (n nginxIngressController) String() string {
	return "Nginx Ingress Controller"
}

func NewNginxIngressControllerReconciler(cloudProvider cloud.Provider) reconciliation.Reconciler {
	return &nginxIngressController{
		cloudProvider: cloudProvider,
	}
}
