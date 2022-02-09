package v1alpha1

const ClusterKind = "Cluster"

type ClusterSpecPlugins struct {
	CertBot                bool `json:"certBot"`
	NginxIngressController bool `json:"nginxIngressController"`
	Vault                  bool `json:"vault"`
	Prometheus             bool `json:"prometheus"`
}

type ClusterSpec struct {
	RootDomain string             `json:"rootDomain"`
	AdminEmail string             `json:"adminEmail"`
	Plugins    ClusterSpecPlugins `json:"plugins"`
}

type Cluster struct {
	TypeMeta `json:",inline"`
	Metadata Metadata    `json:"metadata"`
	Spec     ClusterSpec `json:"spec"`
}

func NewDefaultCluster() Cluster {
	return Cluster{
		TypeMeta: TypeMeta{
			Kind:       ClusterKind,
			APIVersion: apiVersion,
		},
		Spec: ClusterSpec{
			Plugins: ClusterSpecPlugins{
				CertBot:                true,
				NginxIngressController: true,
				Vault:                  true,
			},
		},
	}
}
