package v1alpha1

const PluginKind = "Plugin"

type Plugin struct {
	TypeMeta `json:",inline"`
	Metadata Metadata   `json:"metadata"`
	Spec     PluginSpec `json:"spec"`
}

type PluginSpec struct {
	// Path to either a URL, tar.gz file location or directory containing chart and values
	HelmChart string `json:"helmChart"`
	// Values to apply to the chart
	Values string `json:"values"`
	// Secrets requests secrets from a path in the secret manager and populates a secret named after the plugin
	Secrets map[string]string `json:"secrets"`
}

func NewPlugin(name string) Plugin {
	return Plugin{
		TypeMeta: TypeMeta{
			Kind:       PluginKind,
			APIVersion: apiVersion,
		},
		Metadata: Metadata{
			Name: name,
		},
		Spec: PluginSpec{},
	}
}
