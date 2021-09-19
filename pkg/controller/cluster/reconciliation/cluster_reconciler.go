package reconciliation

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path"

	"github.com/deifyed/xctl/pkg/cloud"
	"github.com/deifyed/xctl/pkg/config"
	"github.com/deifyed/xctl/pkg/controller/common/reconciliation"
	"github.com/spf13/afero"
)

type clusterReconciler struct {
	clusterService cloud.ClusterService
}

func (c *clusterReconciler) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	action := reconciliation.DetermineUserIndication(rctx, true)

	clusterExists, err := c.clusterService.HasCluster(rctx.Ctx, rctx.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("checking cluster existence: %w", err)
	}

	switch action {
	case reconciliation.ActionCreate:
		if !clusterExists {
			err = c.clusterService.CreateCluster(rctx.Ctx, rctx.ClusterDeclaration)
			if err != nil {
				return reconciliation.Result{}, fmt.Errorf("creating cluster: %w", err)
			}
		}

		err = generateKubeconfig(rctx.Ctx, rctx.Filesystem, c.clusterService, rctx.ClusterDeclaration.Metadata.Name)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("generating kubeconfig: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		if !clusterExists {
			return reconciliation.Result{Requeue: false}, nil
		}

		err := c.clusterService.DeleteCluster(rctx.Ctx, rctx.ClusterDeclaration.Metadata.Name)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("deleting cluster: %w", err)
		}

		clusterDir, err := config.GetAbsoluteXCTLClusterDir(rctx.ClusterDeclaration.Metadata.Name)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("acquiring cluster directory: %w", err)
		}

		err = rctx.Filesystem.RemoveAll(clusterDir)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("deleting cluster directory: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{}, reconciliation.ErrIndecisive
}

func generateKubeconfig(ctx context.Context, fs *afero.Afero, provider cloud.ClusterService, clusterName string) error {
	rawConfig, err := provider.GetKubeConfig(ctx, clusterName)
	if err != nil {
		return fmt.Errorf("getting kubeconfig: %w", err)
	}

	decodedConfig, err := base64.StdEncoding.DecodeString(string(rawConfig))
	if err != nil {
		return fmt.Errorf("decoding kubeconfig: %w", err)
	}

	kubeConfigPath, err := config.GetAbsoluteKubeconfigPath(clusterName)
	if err != nil {
		return fmt.Errorf("acquiring KubeConfigPath: %w", err)
	}

	err = os.MkdirAll(path.Dir(kubeConfigPath), 0o700)
	if err != nil {
		return fmt.Errorf("preparing folder structure: %w", err)
	}

	err = fs.WriteFile(kubeConfigPath, decodedConfig, 0o600)
	if err != nil {
		return fmt.Errorf("writing kubeconfig: %w", err)
	}

	return nil
}

func (c *clusterReconciler) String() string {
	return "Kubernetes Cluster"
}

func NewClusterReconciler(clusterService cloud.ClusterService) reconciliation.Reconciler {
	return &clusterReconciler{clusterService: clusterService}
}
