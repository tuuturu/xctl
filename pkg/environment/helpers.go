package environment

import (
	"fmt"
	"io"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"

	"sigs.k8s.io/yaml"
)

// extractEnvironmentManifest knows how to produce an environment manifest from a reader source
func extractEnvironmentManifest(source io.Reader) (v1alpha1.Environment, error) {
	manifest := v1alpha1.NewDefaultEnvironment()

	rawManifest, err := io.ReadAll(source)
	if err != nil {
		return v1alpha1.Environment{}, fmt.Errorf("reading: %w", err)
	}

	err = yaml.Unmarshal(rawManifest, &manifest)
	if err != nil {
		return v1alpha1.Environment{}, fmt.Errorf("unmarshalling: %w", err)
	}

	return manifest, nil
}
