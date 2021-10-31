package v1alpha1

const ClusterKind = "Cluster"

type ClusterSpec struct {
	RootDomain string `json:"rootDomain"`
	AdminEmail string `json:"adminEmail"`
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
	}
}
