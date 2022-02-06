package cloud

import (
	"context"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

const (
	// DefaultAutoscalerMinimumNodes defines the minimum amount of nodes that should be available
	DefaultAutoscalerMinimumNodes = 2
	// DefaultAutoscalerMaximumNodes defines the maximum amount of nodes that should be available
	DefaultAutoscalerMaximumNodes = 10
)

type Cluster struct {
	// Name represents a way to identify a Cluster
	Name string
	// Ready represents whether the cluster is ready to be operated
	Ready bool
}

type ClusterService interface {
	// CreateCluster knows how to create a cluster
	CreateCluster(ctx context.Context, manifest v1alpha1.Cluster) error
	// DeleteCluster knows how to delete a cluster
	DeleteCluster(ctx context.Context, clusterName string) error
	// GetCluster knows how to retrieve information regarding a Cluster
	GetCluster(ctx context.Context, clusterName string) (Cluster, error)
	// HasCluster knows if a cluster exists
	HasCluster(ctx context.Context, clusterName string) (bool, error)
	// GetKubeConfig knows how to retrieve a KubeConfig
	GetKubeConfig(ctx context.Context, clusterName string) ([]byte, error)
}
