package v1alpha1

const apiVersion = "v1alpha1"

type TypeMeta struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
}

type Metadata struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}
