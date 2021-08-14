package v1alpha1

import (
	"fmt"
	"io"

	"sigs.k8s.io/yaml"
)

// InferKindFromManifest knows how to determine the kind of a manifest
func InferKindFromManifest(reader io.Reader) (string, error) {
	var parser struct {
		TypeMeta `json:",inline"`
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("reading data: %w", err)
	}

	err = yaml.Unmarshal(content, &parser)
	if err != nil {
		return "", fmt.Errorf("unmarshalling data: %w", err)
	}

	if parser.Kind == "" {
		return "", fmt.Errorf("determining kind")
	}

	return parser.Kind, nil
}
