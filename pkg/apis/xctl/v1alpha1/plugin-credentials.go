package v1alpha1

// PluginCredentials contains information required to access a plugin
type PluginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
