package prometheus

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

func NewPlugin() v1alpha1.Plugin {
	plugin := v1alpha1.NewPlugin(pluginName)

	plugin.Metadata.Name = pluginName
	plugin.Metadata.Namespace = config.DefaultMonitoringNamespace

	// URL: https://artifacthub.io/packages/helm/prometheus-community/prometheus
	plugin.Spec.Helm.Chart = "prometheus"
	plugin.Spec.Helm.Version = "15.10.1"
	plugin.Spec.Helm.Values = valuesTemplate

	plugin.Spec.Helm.Repository.Name = "prometheus-community"
	plugin.Spec.Helm.Repository.URL = "https://prometheus-community.github.io/helm-charts"

	plugin.Spec.Manifests = []string{datasourceConfigmap()}

	return plugin
}

func datasourceConfigmap() string {
	t := template.Must(template.New("datasource").Parse(datasourceConfigMapTemplate))

	buf := bytes.Buffer{}

	_ = t.Execute(&buf, struct {
		MonitoringNamespace string
	}{MonitoringNamespace: config.DefaultMonitoringNamespace})

	return buf.String()
}

//go:embed values.yaml
var valuesTemplate string //nolint:gochecknoglobals
//go:embed datasource.yaml
var datasourceConfigMapTemplate string

const pluginName = "prometheus"
