package v1alpha1

const ApplicationKind = "Application"

type Application struct {
	TypeMeta `json:",inline"`
	Metadata Metadata `json:"metadata"`
}

func NewApplication() Application {
	return Application{
		TypeMeta: TypeMeta{
			Kind:       ApplicationKind,
			APIVersion: apiVersion,
		},
	}
}
