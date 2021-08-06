package v1alpha1

import (
	"fmt"

	"sigs.k8s.io/yaml"
)

func InferKindFromManifest(data []byte) (string, error) {
	var parser typeParser

	err := yaml.Unmarshal(data, &parser) //nolint:typecheck
	if err != nil {
		return "", fmt.Errorf("unmarshalling data: %w", err)
	}

	if parser.Kind == "" {
		return "", fmt.Errorf("determining kind")
	}

	return parser.Kind, nil
}
