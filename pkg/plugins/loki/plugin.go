package loki

import (
	_ "embed"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

const (
	pluginName      = "loki"
	pluginNamespace = "kube-system"
)

func NewPlugin() v1alpha1.Plugin {
	plugin := v1alpha1.NewPlugin(pluginName)

	plugin.Metadata.Name = pluginName
	plugin.Metadata.Namespace = pluginNamespace

	// URL: https://artifacthub.io/packages/helm/grafana/loki
	plugin.Spec.Helm.Chart = "loki"
	plugin.Spec.Helm.Version = "2.11.0"

	plugin.Spec.Helm.Repository.Name = "grafana"
	plugin.Spec.Helm.Repository.URL = "https://grafana.github.io/helm-charts"

	plugin.Spec.Helm.Values = rawValues

	return plugin
}

//go:embed values.yaml
var rawValues string //nolint:gochecknoglobals
