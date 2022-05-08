package v1alpha1

const EnvironmentKind = "Environment"

type EnvironmentSpecPlugins struct {
	CertManager            bool `json:"certManager"`
	NginxIngressController bool `json:"nginxIngressController"`
	Prometheus             bool `json:"prometheus"`
	Grafana                bool `json:"grafana"`
	Loki                   bool `json:"loki"`
	Promtail               bool `json:"promtail"`
	ArgoCD                 bool `json:"argoCD"`
}

type EnvironmentSpec struct {
	Domain     string                 `json:"domain"`
	AdminEmail string                 `json:"adminEmail"`
	Plugins    EnvironmentSpecPlugins `json:"plugins"`
}

type Environment struct {
	TypeMeta `json:",inline"`
	Metadata Metadata        `json:"metadata"`
	Spec     EnvironmentSpec `json:"spec"`
}

func NewDefaultEnvironment() Environment {
	return Environment{
		TypeMeta: TypeMeta{
			Kind:       EnvironmentKind,
			APIVersion: apiVersion,
		},
		Spec: EnvironmentSpec{
			Plugins: EnvironmentSpecPlugins{
				CertManager:            true,
				NginxIngressController: true,
				Prometheus:             true,
				Grafana:                true,
				Loki:                   true,
				Promtail:               true,
				ArgoCD:                 false,
			},
		},
	}
}
