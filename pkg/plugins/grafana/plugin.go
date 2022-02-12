package grafana

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

const pluginName = "grafana"

func NewPlugin(cluster v1alpha1.Cluster) (v1alpha1.Plugin, error) {
	plugin := v1alpha1.NewPlugin(pluginName)

	plugin.Metadata.Name = pluginName
	plugin.Metadata.Namespace = "kube-system"

	plugin.Spec.Helm.Chart = "grafana"
	plugin.Spec.Helm.Version = "6.21.2"

	plugin.Spec.Helm.Repository.Name = "grafana"
	plugin.Spec.Helm.Repository.URL = "https://grafana.github.io/helm-charts"

	values, err := generateValues(cluster.Spec.RootDomain)
	if err != nil {
		return v1alpha1.Plugin{}, fmt.Errorf("generating values: %w", err)
	}

	plugin.Spec.Helm.Values = values

	return plugin, nil
}

func generateValues(host string) (string, error) {
	t, err := template.New("values").Parse(rawValues)
	if err != nil {
		return "", fmt.Errorf("parsing raw values: %w", err)
	}

	buf := bytes.Buffer{}

	err = t.Execute(&buf, struct {
		Host string
	}{
		Host: host,
	})
	if err != nil {
		return "", fmt.Errorf("injecting variables: %w", err)
	}

	return buf.String(), nil
}

//go:embed values.yaml
var rawValues string //nolint:gochecknoglobals
