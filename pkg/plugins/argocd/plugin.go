package argocd

import (
	_ "embed"
	"fmt"

	"sigs.k8s.io/yaml"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

func NewPlugin() (v1alpha1.Plugin, error) {
	plugin := v1alpha1.NewPlugin(pluginName)

	err := yaml.Unmarshal(rawPlugin, &plugin)
	if err != nil {
		return v1alpha1.Plugin{}, fmt.Errorf("unmarshalling plugin: %w", err)
	}

	plugin.Spec.Helm.Values = rawValues

	return plugin, nil
}

const pluginName = "argocd"

var (
	//go:embed plugin.yaml
	rawPlugin []byte //nolint:gochecknoglobals
	//go:embed values.yaml
	rawValues string //nolint:gochecknoglobals
)
