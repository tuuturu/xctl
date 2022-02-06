package certbot

import (
	"fmt"

	"github.com/deifyed/xctl/pkg/clients/kubectl"
	ingress "github.com/deifyed/xctl/pkg/plugins/nginx-ingress-controller"

	"github.com/deifyed/xctl/pkg/clients/kubectl/binary"

	helmBinary "github.com/deifyed/xctl/pkg/clients/helm/binary"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/cloud"

	"github.com/deifyed/xctl/pkg/controller/common/reconciliation"
)

func (n certbotReconciler) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	log := logging.GetLogger(logFeature, "reconciliation")

	kubeConfigPath, err := config.GetAbsoluteKubeconfigPath(rctx.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("acquiring KubeConfig path: %w", err)
	}

	helmClient, err := helmBinary.New(rctx.Filesystem, kubeConfigPath)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("acquiring Helm client: %w", err)
	}

	kubectlClient, err := binary.New(rctx.Filesystem, kubeConfigPath)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("acquiring Kubectl client: %w", err)
	}

	plugin := newCertbotPlugin()

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

		log.Debug("configuring cluster issuer")

		err = kubectlClient.Apply(kubectl.ApplyOpts{
			Manifest: newLetsEncryptClusterIssuer(rctx.ClusterDeclaration.Spec.AdminEmail),
		})
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("creating cluster issuer: %w", err)
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

func (n certbotReconciler) determineAction(opts determineActionOpts) (reconciliation.Action, error) {
	log := opts.Logger
	indication := reconciliation.DetermineUserIndication(opts.Ctx, opts.Ctx.ClusterDeclaration.Spec.Plugins.CertBot)

	clusterExists, err := n.cloudProvider.HasCluster(opts.Ctx.Ctx, opts.Ctx.ClusterDeclaration)
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf("checking cluster existence: %w", err)
	}

	componentExists := false

	if clusterExists {
		componentExists, err = opts.Helm.Exists(opts.Plugin)
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

		if componentExists {
			log.Debug("Noop due to existing component")

			return reconciliation.ActionNoop, nil
		}

		ingressTester := func() (bool, error) {
			return opts.Helm.Exists(ingress.NewNginxIngressControllerPlugin())
		}

		hasDependencies, err := reconciliation.AssertDependencyExistence(true, ingressTester)
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("asserting dependencies existence: %w", err)
		}

		if !hasDependencies {
			return reconciliation.ActionWait, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !componentExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

func (n certbotReconciler) String() string {
	return "Certbot"
}

func NewCertbotReconciler(cloudProvider cloud.Provider) reconciliation.Reconciler {
	return &certbotReconciler{
		cloudProvider: cloudProvider,
	}
}
