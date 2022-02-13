package certbot

import (
	"fmt"

	reconciliation2 "github.com/deifyed/xctl/pkg/tools/reconciliation"

	"github.com/deifyed/xctl/pkg/tools/clients/helm"
	helmBinary "github.com/deifyed/xctl/pkg/tools/clients/helm/binary"
	"github.com/deifyed/xctl/pkg/tools/clients/kubectl/binary"

	"github.com/deifyed/xctl/pkg/tools/manifests"

	"github.com/pkg/errors"

	ingress "github.com/deifyed/xctl/pkg/plugins/nginx-ingress-controller"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/cloud"
)

func (n certbotReconciler) Reconcile(rctx reconciliation2.Context) (reconciliation2.Result, error) {
	log := logging.GetLogger(logFeature, "reconciliation")

	kubeConfigPath, err := config.GetAbsoluteKubeconfigPath(rctx.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation2.Result{}, fmt.Errorf("acquiring KubeConfig path: %w", err)
	}

	helmClient, err := helmBinary.New(rctx.Filesystem, kubeConfigPath)
	if err != nil {
		return reconciliation2.Result{}, fmt.Errorf("acquiring Helm client: %w", err)
	}

	kubectlClient, err := binary.New(rctx.Filesystem, kubeConfigPath)
	if err != nil {
		return reconciliation2.Result{}, fmt.Errorf("acquiring Kubectl client: %w", err)
	}

	plugin := newCertbotPlugin()

	action, err := n.determineAction(determineActionOpts{
		Ctx:    rctx,
		Helm:   helmClient,
		Plugin: plugin,
		Logger: log,
	})
	if err != nil {
		if errors.Is(err, helm.ErrUnreachable) {
			log.Debug("requeuing due to helm: cluster unreachable")

			return reconciliation2.Result{Requeue: true}, nil
		}

		return reconciliation2.Result{Requeue: false}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation2.ActionCreate:
		log.Debug("installing")

		err = helmClient.Install(plugin)
		if err != nil {
			return reconciliation2.Result{Requeue: false}, fmt.Errorf("installing helm chart: %w", err)
		}

		log.Debug("configuring cluster issuer")

		manifest, err := manifests.ResourceAsReader(newLetsEncryptClusterIssuer(rctx.ClusterDeclaration.Spec.AdminEmail))
		if err != nil {
			return reconciliation2.Result{}, fmt.Errorf("creating cluster issuer: %w", err)
		}

		err = kubectlClient.Apply(manifest)
		if err != nil {
			return reconciliation2.Result{}, fmt.Errorf("creating cluster issuer: %w", err)
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

func (n certbotReconciler) determineAction(opts determineActionOpts) (reconciliation2.Action, error) {
	log := opts.Logger
	indication := reconciliation2.DetermineUserIndication(opts.Ctx, opts.Ctx.ClusterDeclaration.Spec.Plugins.CertBot)

	clusterExists, err := n.cloudProvider.HasCluster(opts.Ctx.Ctx, opts.Ctx.ClusterDeclaration)
	if err != nil {
		return reconciliation2.ActionNoop, fmt.Errorf("checking cluster existence: %w", err)
	}

	componentExists := false

	if clusterExists {
		componentExists, err = opts.Helm.Exists(opts.Plugin)
		if err != nil {
			if errors.Is(err, helm.ErrUnreachable) {
				return "", fmt.Errorf("checking component existence: %w", err)
			}

			return reconciliation2.ActionNoop, fmt.Errorf("checking component existence: %w", err)
		}
	}

	switch indication {
	case reconciliation2.ActionCreate:
		if !clusterExists {
			log.Debug("Waiting due to cluster not ready")

			return reconciliation2.ActionWait, nil
		}

		if componentExists {
			log.Debug("Noop due to existing component")

			return reconciliation2.ActionNoop, nil
		}

		ingressTester := func() (bool, error) {
			return opts.Helm.Exists(ingress.NewNginxIngressControllerPlugin())
		}

		hasDependencies, err := reconciliation2.AssertDependencyExistence(true, ingressTester)
		if err != nil {
			return reconciliation2.ActionNoop, fmt.Errorf("asserting dependencies existence: %w", err)
		}

		if !hasDependencies {
			return reconciliation2.ActionWait, nil
		}

		return reconciliation2.ActionCreate, nil
	case reconciliation2.ActionDelete:
		if !componentExists {
			return reconciliation2.ActionNoop, nil
		}

		return reconciliation2.ActionDelete, nil
	}

	return reconciliation2.ActionNoop, reconciliation2.ErrIndecisive
}

func (n certbotReconciler) String() string {
	return "Certbot"
}

func NewReconciler(cloudProvider cloud.Provider) reconciliation2.Reconciler {
	return &certbotReconciler{
		cloudProvider: cloudProvider,
	}
}
