package linode

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/pkg/errors"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/cloud"
	"github.com/linode/linodego"
	"golang.org/x/oauth2"
)

func (p *provider) Authenticate() error {
	apiKey, ok := os.LookupEnv(linodego.APIEnvVar)
	if !ok {
		return fmt.Errorf(fmt.Sprintf("finding Linode API token (%s)", linodego.APIEnvVar))
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiKey})

	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{Source: tokenSource},
	}

	p.client = linodego.NewClient(oauth2Client)

	return nil
}

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
		return fmt.Errorf("creating cluster: %w", err)
	}

	_, err = p.client.WaitForLKEClusterStatus(ctx, cluster.ID, linodego.LKEClusterReady, defaultTimeoutSeconds)
	if err != nil {
		return fmt.Errorf("waiting for cluster to become ready: %w", err)
	}

	return nil
}

func (p *provider) DeleteCluster(ctx context.Context, clusterName string) error {
	lkeCluster, err := p.getCluster(ctx, clusterName)
	if err != nil {
		if errors.Is(err, config.ErrNotFound) {
			return nil
		}

		return fmt.Errorf("querying clusters: %w", err)
	}

	err = p.client.DeleteLKECluster(ctx, lkeCluster.ID)
	if err != nil {
		return fmt.Errorf("deleting cluster: %w", err)
	}

	return nil
}

func (p *provider) GetCluster(ctx context.Context, clusterName string) (cloud.Cluster, error) {
	lkeCluster, err := p.getCluster(ctx, clusterName)
	if err != nil {
		return cloud.Cluster{}, fmt.Errorf("querying clusters: %w", err)
	}

	return cloud.Cluster{
		Name: lkeCluster.Label,
	}, nil
}

func (p *provider) HasCluster(ctx context.Context, clusterName string) (bool, error) {
	_, err := p.getCluster(ctx, clusterName)
	if err != nil {
		if errors.Is(err, config.ErrNotFound) {
			return false, nil
		}

		return false, fmt.Errorf("querying clusters: %w", err)
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

func (p *provider) getCluster(ctx context.Context, clusterName string) (linodego.LKECluster, error) {
	clusters, err := p.client.ListLKEClusters(ctx, &linodego.ListOptions{})
	if err != nil {
		return linodego.LKECluster{}, fmt.Errorf("retrieving existing LKE clusters: %w", err)
	}

	for _, cluster := range clusters {
		if cluster.Label == clusterName {
			return cluster, nil
		}
	}

	return linodego.LKECluster{}, config.ErrNotFound
}

func NewLinodeProvider() cloud.Provider {
	return &provider{}
}
