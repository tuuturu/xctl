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

type ClusterNode struct {
	Name string
	IPv4 string
}

type Cluster struct {
	// Name represents a way to identify a Cluster
	Name string
	// Ready represents whether the cluster is ready to be operated
	Ready bool
	// PublicIPv6 represents the IP of which the cluster is available for public requests
	PublicIPv6 string
	// Nodes contains details about the cluster's nodes
	Nodes []ClusterNode
}

type ClusterService interface {
	// CreateCluster knows how to create a cluster
	CreateCluster(ctx context.Context, manifest v1alpha1.Cluster) error
	// DeleteCluster knows how to delete a cluster
	DeleteCluster(ctx context.Context, manifest v1alpha1.Cluster) error
	// GetCluster knows how to retrieve information regarding a Cluster
	GetCluster(ctx context.Context, manifest v1alpha1.Cluster) (Cluster, error)
	// HasCluster knows if a cluster exists
	HasCluster(ctx context.Context, manifest v1alpha1.Cluster) (bool, error)
	// GetKubeConfig knows how to retrieve a KubeConfig
	GetKubeConfig(ctx context.Context, manifest v1alpha1.Cluster) ([]byte, error)
}
