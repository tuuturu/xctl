package v1alpha1

const ClusterKind = "Cluster"

func NewCluster() Cluster {
	return Cluster{
		TypeMeta: TypeMeta{
			Kind:       ClusterKind,
			APIVersion: apiVersion,
		},
	}
}
