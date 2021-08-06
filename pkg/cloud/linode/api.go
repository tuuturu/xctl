package linode

import (
	"context"
	"fmt"
	"net/http"
	"os"

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

func (p *provider) DeleteCluster(ctx context.Context, manifest v1alpha1.Cluster) error {
	clusters, err := p.client.ListLKEClusters(ctx, &linodego.ListOptions{})
	if err != nil {
		return fmt.Errorf("retrieving existing LKE clusters: %w", err)
	}

	id := -1

	for _, cluster := range clusters {
		if cluster.Label == manifest.Metadata.Name {
			id = cluster.ID

			break
		}
	}

	if id == -1 {
		return fmt.Errorf("finding cluster with name: %s", manifest.Metadata.Name)
	}

	err = p.client.DeleteLKECluster(ctx, id)
	if err != nil {
		return fmt.Errorf("deleting cluster: %w", err)
	}

	return nil
}

func NewLinodeProvider() cloud.Provider {
	return &provider{}
}
