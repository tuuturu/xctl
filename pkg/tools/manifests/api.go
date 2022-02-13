package manifests

import (
	"bytes"
	"fmt"
	"io"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"

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

// ExtractEnvironmentManifest knows how to produce an environment manifest from a reader source
func ExtractEnvironmentManifest(source io.Reader) (v1alpha1.Cluster, error) {
	manifest := v1alpha1.NewDefaultCluster()

	rawManifest, err := io.ReadAll(source)
	if err != nil {
		return v1alpha1.Cluster{}, fmt.Errorf("reading: %w", err)
	}

	err = yaml.Unmarshal(rawManifest, &manifest)
	if err != nil {
		return v1alpha1.Cluster{}, fmt.Errorf("unmarshalling: %w", err)
	}

	return manifest, nil
}
