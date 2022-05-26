package environment

import (
	"context"
	_ "embed"
	"encoding/base64"
	"fmt"
	"path"
	"strings"

	ingress "github.com/deifyed/xctl/pkg/environment/plugins/nginx-ingress-controller"

	kubectlBinary "github.com/deifyed/xctl/pkg/tools/clients/kubectl/binary"

	"github.com/deifyed/xctl/pkg/tools/reconciliation"

	"github.com/deifyed/xctl/pkg/tools/clients/helm"
	helmBinary "github.com/deifyed/xctl/pkg/tools/clients/helm/binary"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/deifyed/xctl/pkg/cloud"
	"github.com/deifyed/xctl/pkg/config"
	"github.com/spf13/afero"
)

//nolint:funlen
// Reconcile knows how to ensure reality for a cluster is as declared in an environment manifest
func (c *clusterReconciler) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	log := logging.GetLogger("cluster", "reconciliation")
	action := reconciliation.DetermineUserIndication(rctx, true)

	clusterExists, err := c.clusterService.HasCluster(rctx.Ctx, rctx.EnvironmentManifest)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("checking cluster existence: %w", err)
	}

	kubeconfigPath, err := config.GetAbsoluteKubeconfigPath(rctx.EnvironmentManifest.Metadata.Name)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("acquiring kubeconfig path: %w", err)
	}

	kubectlClient, err := kubectlBinary.New(rctx.Filesystem, kubeconfigPath)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("preparing kubectl client: %w", err)
	}

	switch action {
	case reconciliation.ActionCreate:
		log.Debug("creating")

		if !clusterExists {
			err = c.clusterService.CreateCluster(rctx.Ctx, rctx.EnvironmentManifest)
			if err != nil {
				return reconciliation.Result{}, fmt.Errorf("creating cluster: %w", err)
			}
		}

		err = generateKubeconfig(rctx.Ctx, rctx.Filesystem, c.clusterService, rctx.EnvironmentManifest)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("generating kubeconfig: %w", err)
		}

		err = kubectlClient.Apply(strings.NewReader(namespacesTemplate))
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("applying namespaces: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		log.Debug("deleting")

		if !clusterExists {
			return reconciliation.Result{Requeue: false}, nil
		}

		kubeConfigPath, err := config.GetAbsoluteKubeconfigPath(rctx.EnvironmentManifest.Metadata.Name)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("acquiring KubeConfig path: %w", err)
		}

		helmClient, err := helmBinary.New(rctx.Filesystem, kubeConfigPath)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("acquiring helm client: %w", err)
		}

		ok, err := ensureDependencies(helmClient)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("checking dependencies: %w", err)
		}

		if !ok {
			return reconciliation.Result{Requeue: true}, nil
		}

		err = c.clusterService.DeleteCluster(rctx.Ctx, rctx.EnvironmentManifest)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("deleting cluster: %w", err)
		}

		clusterDir, err := config.GetAbsoluteXCTLClusterDir(rctx.EnvironmentManifest.Metadata.Name)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("acquiring cluster directory: %w", err)
		}

		err = rctx.Filesystem.RemoveAll(clusterDir)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("deleting cluster directory: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.NoopWaitIndecisiveHandler(action)
}

func ensureDependencies(helmClient helm.Client) (bool, error) {
	ingressControllerExists, err := helmClient.Exists(ingress.NewNginxIngressControllerPlugin())
	if err != nil {
		return false, fmt.Errorf("acquiring ingress controller existence: %w", err)
	}

	if ingressControllerExists {
		return false, nil
	}

	return true, nil
}

func generateKubeconfig(ctx context.Context, fs *afero.Afero, provider cloud.ClusterService, manifest v1alpha1.Environment) error { //nolint:lll
	rawConfig, err := provider.GetKubeConfig(ctx, manifest)
	if err != nil {
		return fmt.Errorf("getting kubeconfig: %w", err)
	}

	decodedConfig, err := base64.StdEncoding.DecodeString(string(rawConfig))
	if err != nil {
		return fmt.Errorf("decoding kubeconfig: %w", err)
	}

	kubeConfigPath, err := config.GetAbsoluteKubeconfigPath(manifest.Metadata.Name)
	if err != nil {
		return fmt.Errorf("acquiring KubeConfigPath: %w", err)
	}

	err = fs.MkdirAll(path.Dir(kubeConfigPath), 0o700)
	if err != nil {
		return fmt.Errorf("preparing folder structure: %w", err)
	}

	err = fs.WriteFile(kubeConfigPath, decodedConfig, 0o600)
	if err != nil {
		return fmt.Errorf("writing kubeconfig: %w", err)
	}

	return nil
}

// String returns a string representing the reconciler
func (c *clusterReconciler) String() string {
	return "Kubernetes Environment"
}

// NewClusterReconciler returns an initialized cluster reconciler
func NewClusterReconciler(clusterService cloud.ClusterService) reconciliation.Reconciler {
	return &clusterReconciler{clusterService: clusterService}
}

type clusterReconciler struct {
	clusterService cloud.ClusterService
}

//go:embed namespaces.yaml
var namespacesTemplate string
