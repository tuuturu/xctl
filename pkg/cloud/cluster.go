package cloud

import (
	"context"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

type Cluster struct {
	Name string
}

type ClusterService interface {
	CreateCluster(ctx context.Context, manifest v1alpha1.Cluster) error
	DeleteCluster(ctx context.Context, clusterName string) error
	GetCluster(ctx context.Context, clusterName string) (Cluster, error)
	HasCluster(ctx context.Context, clusterName string) (bool, error)
	GetKubeConfig(ctx context.Context, clusterName string) ([]byte, error)
}
