package v1alpha1

const ApplicationKind = "Application"

type Application struct {
	TypeMeta `json:",inline"`
	Metadata Metadata        `json:"metadata"`
	Spec     ApplicationSpec `json:"spec"`
}

type ApplicationSpec struct {
	Image string `json:"image"`
	Port  string `json:"port"`
	Url   string `json:"url"`
}

func NewDefaultApplication() Application {
	return Application{
		TypeMeta: TypeMeta{
			Kind:       ApplicationKind,
			APIVersion: apiVersion,
		},
		Metadata: Metadata{
			Name:      "",
			Namespace: "default",
		},
	}
}
