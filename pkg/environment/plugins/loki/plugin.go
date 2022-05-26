package loki

import (
	_ "embed"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

const (
	pluginName      = "loki"
	pluginNamespace = config.DefaultMonitoringNamespace
)

func NewPlugin() v1alpha1.Plugin {
	plugin := v1alpha1.NewPlugin(pluginName)

	plugin.Metadata.Name = pluginName
	plugin.Metadata.Namespace = pluginNamespace

	// URL: https://artifacthub.io/packages/helm/grafana/loki
	plugin.Spec.Helm.Chart = "loki"
	plugin.Spec.Helm.Version = "2.11.1"

	plugin.Spec.Helm.Repository.Name = "grafana"
	plugin.Spec.Helm.Repository.URL = "https://grafana.github.io/helm-charts"

	plugin.Spec.Helm.Values = rawValues
	plugin.Spec.Manifests = []string{rawDatasourcesCM}

	return plugin
}

var (
	//go:embed values.yaml
	rawValues string //nolint:gochecknoglobals
	//go:embed datasource.yaml
	rawDatasourcesCM string //nolint:gochecknoglobals
)
