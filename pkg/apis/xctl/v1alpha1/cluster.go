package v1alpha1

const ClusterKind = "Cluster"

type Plugins struct {
	CertBot                bool `json:"certBot"`
	NginxIngressController bool `json:"nginxIngressController"`
	Vault                  bool `json:"vault"`
}

type ClusterSpec struct {
	RootDomain string  `json:"rootDomain"`
	AdminEmail string  `json:"adminEmail"`
	Plugins    Plugins `json:"plugins"`
}

type Cluster struct {
	TypeMeta `json:",inline"`
	Metadata Metadata    `json:"metadata"`
	Spec     ClusterSpec `json:"spec"`
}

func NewCluster() Cluster {
	return Cluster{
		TypeMeta: TypeMeta{
			Kind:       ClusterKind,
			APIVersion: apiVersion,
		},
		Spec: ClusterSpec{
			Plugins: Plugins{
				CertBot:                true,
				NginxIngressController: true,
				Vault:                  true,
			},
		},
	}
}
