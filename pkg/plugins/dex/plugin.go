package dex

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"sigs.k8s.io/yaml"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

func NewPlugin(url string) (v1alpha1.Plugin, error) {
	plugin := v1alpha1.NewPlugin(pluginName)

	err := yaml.Unmarshal(rawPlugin, &plugin)
	if err != nil {
		return v1alpha1.Plugin{}, fmt.Errorf("unmarshalling plugin: %w", err)
	}

	t, err := template.New("values").Parse(rawValues)
	if err != nil {
		return v1alpha1.Plugin{}, fmt.Errorf("parsing template: %w", err)
	}

	buf := bytes.Buffer{}

	err = t.Execute(&buf, pluginOpts{URL: url})
	if err != nil {
		return v1alpha1.Plugin{}, fmt.Errorf("executing template: %w", err)
	}

	plugin.Spec.Helm.Values = buf.String()

	return plugin, nil
}

const pluginName = "Dex"

type pluginOpts struct {
	URL               string `json:"url"`
	ArgoCDRedirectURI string `json:"argoCDRedirectURI"`
}

var (
	//go:embed plugin.yaml
	rawPlugin []byte //nolint:gochecknoglobals
	//go:embed values.yaml
	rawValues string //nolint:gochecknoglobals
)
