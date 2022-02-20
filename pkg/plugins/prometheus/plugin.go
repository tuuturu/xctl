package prometheus

import (
	_ "embed"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

func NewPlugin() v1alpha1.Plugin {
	plugin := v1alpha1.NewPlugin(pluginName)

	plugin.Metadata.Name = pluginName
	plugin.Metadata.Namespace = "kube-system"

	// URL: https://artifacthub.io/packages/helm/prometheus-community/prometheus/15.1.3
	plugin.Spec.Helm.Chart = "prometheus"
	plugin.Spec.Helm.Version = "15.1.3"
	plugin.Spec.Helm.Values = template

	plugin.Spec.Helm.Repository.Name = "prometheus-community"
	plugin.Spec.Helm.Repository.URL = "https://prometheus-community.github.io/helm-charts"

	return plugin
}

//go:embed values.yaml
var template string //nolint:gochecknoglobals

const pluginName = "prometheus"
