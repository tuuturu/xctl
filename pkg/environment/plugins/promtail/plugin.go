package promtail

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

const (
	pluginName      = "promtail"
	pluginNamespace = config.DefaultMonitoringNamespace
)

func NewPlugin() v1alpha1.Plugin {
	plugin := v1alpha1.NewPlugin(pluginName)

	plugin.Metadata.Name = pluginName
	plugin.Metadata.Namespace = pluginNamespace

	// URL: https://artifacthub.io/packages/helm/grafana/promtail
	plugin.Spec.Helm.Chart = "promtail"
	plugin.Spec.Helm.Version = "5.1.0"

	plugin.Spec.Helm.Repository.Name = "grafana"
	plugin.Spec.Helm.Repository.URL = "https://grafana.github.io/helm-charts"

	plugin.Spec.Helm.Values = values()

	return plugin
}

func values() string {
	t := template.New("values").Delims("{{{", "}}}")
	t = template.Must(t.Parse(rawValues))

	buf := bytes.Buffer{}

	_ = t.Execute(&buf, struct {
		MonitoringNamespace string
	}{MonitoringNamespace: config.DefaultMonitoringNamespace})

	return buf.String()
}

var (
	//go:embed values.yaml
	rawValues string //nolint:gochecknoglobals
)
