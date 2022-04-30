package promtail

import (
	_ "embed"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

const (
	pluginName      = "promtail"
	pluginNamespace = "kube-system"
)

func NewPlugin() v1alpha1.Plugin {
	plugin := v1alpha1.NewPlugin(pluginName)

	plugin.Metadata.Name = pluginName
	plugin.Metadata.Namespace = pluginNamespace

	// URL: https://artifacthub.io/packages/helm/grafana/promtail
	plugin.Spec.Helm.Chart = "promtail"
	plugin.Spec.Helm.Version = "4.2.0"

	plugin.Spec.Helm.Repository.Name = "grafana"
	plugin.Spec.Helm.Repository.URL = "https://grafana.github.io/helm-charts"

	plugin.Spec.Helm.Values = rawValues

	return plugin
}

var (
	//go:embed values.yaml
	rawValues string //nolint:gochecknoglobals
)