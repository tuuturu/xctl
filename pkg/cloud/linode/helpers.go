package linode

import (
	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

func defaultLabels(cluster v1alpha1.Cluster, extraLabels ...string) []string {
	result := []string{
		cluster.Metadata.Name,
	}

	result = append(result, extraLabels...)

	return result
}
