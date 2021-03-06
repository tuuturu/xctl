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
	// Manifests contains manifests to apply to the cluster
	Manifests PluginSpecManifests `json:"manifests"`
}

// PluginSpecHelmRepository contains information about a Helm repository
type PluginSpecHelmRepository struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// PluginSpecHelm contains necessary information for installing a Helm chart
type PluginSpecHelm struct {
	// Repository defines the repository that contains the chart
	Repository PluginSpecHelmRepository `json:"repository"`
	// Chart defines the URL where the chart can be found
	Chart string `json:"chart"`
	// Values defines the values to apply to the chart
	Values string `json:"values"`
	// Version defines the chart version
	Version string
}

type PluginSpecManifests []string

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
