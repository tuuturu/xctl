package application

import (
	"fmt"
	"io"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"sigs.k8s.io/yaml"
)

func extractManifest(in io.Reader) (v1alpha1.Application, error) {
	raw, err := io.ReadAll(in)
	if err != nil {
		return v1alpha1.Application{}, fmt.Errorf("buffering: %w", err)
	}

	application := v1alpha1.NewDefaultApplication()

	err = yaml.Unmarshal(raw, &application)
	if err != nil {
		return v1alpha1.Application{}, fmt.Errorf("unmarshalling: %w", err)
	}

	return application, nil
}
