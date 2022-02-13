package binary

type kubeConfig struct {
	CurrentContext string              `json:"current_context"`
	Users          []kubeConfigUsers   `json:"users"`
	Contexts       []kubeConfigContext `json:"contexts"`
}

type kubeConfigUsers struct {
	Name string              `json:"name"`
	User kubeConfigUsersUser `json:"user"`
}

type kubeConfigUsersUser struct {
	Token string `json:"token"`
}

type kubeConfigContext struct {
	Name    string                   `json:"name"`
	Context kubeConfigContextContext `json:"context"`
}

type kubeConfigContextContext struct {
	User string `json:"user"`
}
