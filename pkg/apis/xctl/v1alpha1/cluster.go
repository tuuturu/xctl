package v1alpha1

const ClusterKind = "Cluster"

type Cluster struct {
	TypeMeta `json:",inline"`
	Metadata Metadata `json:"metadata"`
	URL      string   `json:"url"`
}

func NewCluster() Cluster {
	return Cluster{
		TypeMeta: TypeMeta{
			Kind:       ClusterKind,
			APIVersion: apiVersion,
		},
	}
}
