package v1alpha1

const EnvironmentKind = "Environment"

type EnvironmentSpecPlugins struct {
	CertBot                bool `json:"certBot"`
	NginxIngressController bool `json:"nginxIngressController"`
	Prometheus             bool `json:"prometheus"`
	Grafana                bool `json:"grafana"`
	Loki                   bool `json:"loki"`
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
				CertBot:                true,
				NginxIngressController: true,
				Prometheus:             true,
				Grafana:                true,
				Loki:                   true,
			},
		},
	}
}
