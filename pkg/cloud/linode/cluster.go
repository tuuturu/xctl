package linode

import (
	"context"
	"errors"
	"fmt"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/cloud"
	"github.com/deifyed/xctl/pkg/config"
	"github.com/linode/linodego"
)

func (p *provider) CreateCluster(ctx context.Context, manifest v1alpha1.Cluster) error {
	cluster, err := p.client.CreateLKECluster(ctx, linodego.LKEClusterCreateOptions{
		NodePools: []linodego.LKEClusterPoolCreateOptions{
			{
				Count: config.DefaultClusterNodeAmount,
				Type:  linodeType4GB,
			},
		},
		Label:      manifest.Metadata.Name,
		Region:     regionFrankfurt,
		K8sVersion: defaultKubernetesVersion,
	})
	if err != nil {
		if errorIsNotAuthenticated(err) {
			return config.ErrNotAuthenticated
		}

		return fmt.Errorf("creating cluster: %w", err)
	}

	err = p.awaitCreation(ctx, cluster.ID)
	if err != nil {
		return fmt.Errorf("awaiting creation of cluster: %w", err)
	}

	return nil
}

func (p *provider) DeleteCluster(ctx context.Context, clusterName string) error {
	cluster, err := p.getCluster(ctx, clusterName)
	if err != nil {
		switch {
		case errors.Is(err, config.ErrNotFound):
			return nil
		case errorIsNotAuthenticated(err):
			return config.ErrNotAuthenticated
		default:
			return fmt.Errorf("querying clusters: %w", err)
		}
	}

	err = p.client.DeleteLKECluster(ctx, cluster.ID)
	if err != nil {
		return fmt.Errorf("deleting cluster: %w", err)
	}

	err = p.awaitDeletion(ctx, clusterName)
	if err != nil {
		return fmt.Errorf("awaiting creation of cluster: %w", err)
	}

	return nil
}

func (p *provider) GetCluster(ctx context.Context, clusterName string) (cloud.Cluster, error) {
	lkeCluster, err := p.getCluster(ctx, clusterName)
	if err != nil {
		if errorIsNotAuthenticated(err) {
			return cloud.Cluster{}, config.ErrNotAuthenticated
		}

		return cloud.Cluster{}, fmt.Errorf("querying clusters: %w", err)
	}

	return cloud.Cluster{
		Name:  lkeCluster.Label,
		Ready: lkeCluster.Status == linodego.LKEClusterReady,
	}, nil
}

func (p *provider) HasCluster(ctx context.Context, clusterName string) (bool, error) {
	_, err := p.getCluster(ctx, clusterName)
	if err != nil {
		switch {
		case errors.Is(err, config.ErrNotFound):
			return false, nil
		case errorIsNotAuthenticated(err):
			return false, config.ErrNotAuthenticated
		default:
			return false, fmt.Errorf("querying clusters: %w", err)
		}
	}

	return true, nil
}

func (p *provider) GetKubeConfig(ctx context.Context, clusterName string) ([]byte, error) {
	cluster, err := p.getCluster(ctx, clusterName)
	if err != nil {
		if errors.Is(err, config.ErrNotFound) {
			return []byte{}, fmt.Errorf("could not find cluster with name %s", clusterName)
		}

		return []byte{}, fmt.Errorf("querying clusters: %w", err)
	}

	cfg, err := p.client.GetLKEClusterKubeconfig(ctx, cluster.ID)
	if err != nil {
		return []byte{}, fmt.Errorf("acquiring kube config: %w", err)
	}

	return []byte(cfg.KubeConfig), nil
}
