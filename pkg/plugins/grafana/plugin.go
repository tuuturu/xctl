package grafana

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

const pluginName = "grafana"

func NewPlugin(opts NewPluginOpts) (v1alpha1.Plugin, error) {
	plugin := v1alpha1.NewPlugin(pluginName)

	plugin.Metadata.Name = pluginName
	plugin.Metadata.Namespace = "kube-system"

	// URL: https://artifacthub.io/packages/helm/grafana/grafana
	plugin.Spec.Helm.Chart = "grafana"
	plugin.Spec.Helm.Version = "6.26.3"

	plugin.Spec.Helm.Repository.Name = "grafana"
	plugin.Spec.Helm.Repository.URL = "https://grafana.github.io/helm-charts"

	values, err := generateValues(opts)
	if err != nil {
		return v1alpha1.Plugin{}, fmt.Errorf("generating values: %w", err)
	}

	plugin.Spec.Helm.Values = values

	return plugin, nil
}

func generateValues(opts NewPluginOpts) (string, error) {
	t, err := template.New("values").Parse(rawValues)
	if err != nil {
		return "", fmt.Errorf("parsing raw values: %w", err)
	}

	buf := bytes.Buffer{}

	err = t.Execute(&buf, opts)
	if err != nil {
		return "", fmt.Errorf("injecting variables: %w", err)
	}

	return buf.String(), nil
}

//go:embed values.yaml
var rawValues string //nolint:gochecknoglobals
