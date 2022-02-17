package manifests

import (
	"bytes"
	"fmt"
	"io"

	"sigs.k8s.io/yaml"
)

// ResourceAsReader converts a Kubernetes resource manifest into a readable stream
func ResourceAsReader(manifest interface{}) (io.Reader, error) {
	result, err := yaml.Marshal(manifest)
	if err != nil {
		return nil, fmt.Errorf("marshalling manifest: %w", err)
	}

	return bytes.NewReader(result), nil
}
