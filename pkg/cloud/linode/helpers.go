package linode

import (
	"strings"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/config"
)

func defaultLabels(cluster v1alpha1.Cluster, extraLabels ...string) []string {
	result := []string{
		componentNamer(cluster, "", ""),
	}

	result = append(result, extraLabels...)

	return result
}

func componentNamer(cluster v1alpha1.Cluster, componentType string, id string) string {
	componentName := strings.Join([]string{config.ApplicationName, cluster.Metadata.Name, componentType, id}, "-")
	componentName = strings.ToLower(componentName)
	componentName = strings.TrimRight(componentName, "-")

	return componentName
}
