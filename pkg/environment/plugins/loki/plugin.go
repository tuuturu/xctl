package loki

import (
	"bytes"
	_ "embed"
	"text/template"
	"time"

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

	plugin.Spec.Helm.Values = values()
	plugin.Spec.Manifests = []string{monitoringNamespaceCM()}

	return plugin
}

func monitoringNamespaceCM() string {
	t := template.Must(template.New("datasource").Parse(rawDatasourcesCM))

	buf := bytes.Buffer{}

	_ = t.Execute(&buf, struct {
		MonitoringNamespace string
	}{MonitoringNamespace: config.DefaultMonitoringNamespace})

	return buf.String()
}

func values() string {
	t := template.Must(template.New("values").Parse(rawValues))

	buf := bytes.Buffer{}

	_ = t.Execute(&buf, struct {
		Date string
	}{time.Now().Format("2006-01-02")})

	return buf.String()
}

var (
	//go:embed values.yaml
	rawValues string //nolint:gochecknoglobals
	//go:embed datasource.yaml
	rawDatasourcesCM string //nolint:gochecknoglobals
)
