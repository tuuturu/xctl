package cloud

import (
	"context"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

type ClusterService interface {
	CreateCluster(ctx context.Context, manifest v1alpha1.Cluster) error
	DeleteCluster(ctx context.Context, manifest v1alpha1.Cluster) error
}
