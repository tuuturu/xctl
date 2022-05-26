package argocd

import (
	"errors"
	"fmt"
	"strings"

	kubectlBinary "github.com/deifyed/xctl/pkg/tools/clients/kubectl/binary"

	"github.com/deifyed/xctl/pkg/config"
	helmBinary "github.com/deifyed/xctl/pkg/tools/clients/helm/binary"

	"github.com/deifyed/xctl/pkg/tools/reconciliation"

	"github.com/deifyed/xctl/pkg/tools/clients/helm"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/deifyed/xctl/pkg/cloud"
)

func (r reconciler) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	log := logging.GetLogger(logFeature, "reconciliation")

	repo := repository{URL: rctx.EnvironmentManifest.Spec.Repository}

	kubeconfigPath, err := config.GetAbsoluteKubeconfigPath(rctx.EnvironmentManifest.Metadata.Name)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("acquiring kubeconfig path: %w", err)
	}

	kubectlClient, err := kubectlBinary.New(rctx.Filesystem, kubeconfigPath)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("preparing kubectl client: %w", err)
	}

	helmClient, err := helmBinary.New(rctx.Filesystem, kubeconfigPath)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("preparing helm client: %w", err)
	}

	action, err := r.determineAction(rctx, helmClient)
	if err != nil {
		return reconciliation.Result{Requeue: false}, fmt.Errorf("determining course of action: %w", err)
	}

	plugin, err := NewPlugin()
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("creating plugin: %w", err)
	}

	switch action {
	case reconciliation.ActionCreate:
		log.Debug("installing")

		err = helmClient.Install(plugin)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("installing: %w", err)
		}

		keys, err := generateKey()
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("generating private/public key pair: %w", err)
		}

		secret, err := generateRepositorySecret(repo, keys.PrivateKey)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("generating repository secret: %w", err)
		}

		err = installDeployKey(rctx.Ctx, rctx.Keyring, repo, keys.PublicKey)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("installing deploy key: %w", err)
		}

		err = kubectlClient.Apply(secret)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("applying secret: %w", err)
		}

		applicationsAppManifest, err := establishConfiguration(
			rctx.Filesystem,
			rctx.RootDirectory,
			rctx.EnvironmentManifest,
		)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("establishing configuration")
		}

		err = kubectlClient.Apply(applicationsAppManifest)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("applying applications app manifest: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		log.Debug("deleting")

		err = helmClient.Delete(plugin)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("uninstalling: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.NoopWaitIndecisiveHandler(action)
}

//nolint:funlen
func (r reconciler) determineAction(rctx reconciliation.Context, helm helm.Client) (reconciliation.Action, error) { //nolint:lll
	indication := reconciliation.DetermineUserIndication(rctx, rctx.EnvironmentManifest.Spec.Plugins.ArgoCD)

	var (
		clusterExists   = true
		componentExists = true
	)

	_, err := r.cloudProvider.GetCluster(rctx.Ctx, rctx.EnvironmentManifest)
	if err != nil {
		if !errors.Is(err, cloud.ErrNotFound) {
			return "", fmt.Errorf("acquiring cluster: %w", err)
		}

		clusterExists = false
	}

	plugin, err := NewPlugin()
	if err != nil {
		return "", fmt.Errorf("creating plugin: %w", err)
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
	return strings.Title(pluginName)
}

func NewReconciler(cloudProvider cloud.Provider) reconciliation.Reconciler {
	return &reconciler{
		cloudProvider: cloudProvider,
	}
}
