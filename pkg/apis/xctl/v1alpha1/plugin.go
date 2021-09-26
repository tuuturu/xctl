package v1alpha1

const PluginKind = "Plugin"

// Plugin defines an installable plugin in xctl
type Plugin struct {
	TypeMeta `json:",inline"`
	Metadata Metadata   `json:"metadata"`
	Spec     PluginSpec `json:"spec"`
}

// PluginSpec contains the different plugin capabilities
type PluginSpec struct {
	// Helm contains a Helm chart to install
	Helm PluginSpecHelm `json:"helm"`
	// Secrets requests secrets from a path in the secret manager and populates a secret named after the plugin
	Secrets map[string]string `json:"secrets"`
	// PostInstallScript defines a script to be run post install
	Hooks PluginSpecHooks `json:"hooks"`
}

// PluginSpecHelm contains necessary information for installing a Helm chart
type PluginSpecHelm struct {
	// Chart defines the URL where the chart can be found
	Chart string `json:"chart"`
	// Values defines the values to apply to the chart
	Values string `json:"values"`
}

// PluginSpecHooks contains scripts ran at certain plugin life cycle events
type PluginSpecHooks struct {
	// PostInstall will trigger after all components in the plugin is successfully installed
	PostInstall string `json:"postInstall"`
	// PostUninstall will trigger after all components in the plugin is successfully uninstalled
	PostUninstall string `json:"postUninstall"`
}

// NewPlugin initializes a plugin
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
