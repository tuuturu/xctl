package grafana

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

const (
	pluginName      = "grafana"
	pluginNamespace = config.DefaultMonitoringNamespace
)

func NewPlugin() v1alpha1.Plugin {
	plugin := v1alpha1.NewPlugin(pluginName)

	plugin.Metadata.Name = pluginName
	plugin.Metadata.Namespace = pluginNamespace

	// URL: https://artifacthub.io/packages/helm/grafana/grafana
	plugin.Spec.Helm.Chart = "grafana"
	plugin.Spec.Helm.Version = "6.30.2"

	plugin.Spec.Helm.Repository.Name = "grafana"
	plugin.Spec.Helm.Repository.URL = "https://grafana.github.io/helm-charts"

	plugin.Spec.Helm.Values = generateValues(valuesOpts{
		SecretName:        secretName(),
		SecretUsernameKey: adminUsernameKey,
		SecretPasswordKey: adminPasswordKey,
	})

	return plugin
}

func generateValues(opts valuesOpts) string {
	t := template.Must(template.New("values").Parse(rawValues))

	buf := bytes.Buffer{}

	err := t.Execute(&buf, opts)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

//go:embed values.yaml
var rawValues string //nolint:gochecknoglobals
