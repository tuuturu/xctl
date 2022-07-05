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
	log := logging.GetLogger("plugin", n.String())

	kubeConfigPath, err := config.GetAbsoluteKubeconfigPath(rctx.EnvironmentManifest.Metadata.Name)
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

	log.Debugf("Action: %s", action)

	switch action {
	case reconciliation.ActionCreate:
		err = helmClient.Install(plugin)
		if err != nil {
			return reconciliation.Result{Requeue: false}, fmt.Errorf("installing helm chart: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err = helmClient.Delete(plugin)
		if err != nil {
			return reconciliation.Result{Requeue: false}, fmt.Errorf("uninstalling helm chart: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.NoopWaitIndecisiveHandler(action)
}

func (n nginxIngressController) determineAction(opts determineActionOpts) (reconciliation.Action, error) {
	indication := reconciliation.DetermineUserIndication(
		opts.Ctx,
		opts.Ctx.EnvironmentManifest.Spec.Plugins.NginxIngressController,
	)

	clusterExists, err := n.cloudProvider.HasCluster(opts.Ctx.Ctx, opts.Ctx.EnvironmentManifest)
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
