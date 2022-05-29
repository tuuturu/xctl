package linode

import (
	"strings"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/config"
)

func defaultLabels(cluster v1alpha1.Environment, extraLabels ...string) []string {
	result := []string{
		componentNamer(cluster, "", ""),
	}

	result = append(result, extraLabels...)

	return result
}

// N.B: max length for Linode labels are 32
func componentNamer(cluster v1alpha1.Environment, componentType string, id string) string {
	componentName := strings.Join([]string{config.ApplicationName, cluster.Metadata.Name, componentType, id}, "-")
	componentName = strings.ToLower(componentName)
	componentName = strings.TrimRight(componentName, "-")

	return componentName
}
