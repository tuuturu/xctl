package v1alpha1

const ApplicationKind = "Application"

func NewApplication() Application {
	return Application{
		TypeMeta: TypeMeta{
			Kind:       ApplicationKind,
			APIVersion: apiVersion,
		},
	}
}
