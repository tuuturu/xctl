package v1alpha1

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestInferKindFromManifest(t *testing.T) {
	testCases := []struct {
		name string

		withManifest []byte
		expectKind   string
	}{
		{
			name: "Should successfully identify cluster",

			withManifest: func() []byte {
				manifest := Cluster{
					TypeMeta: TypeMeta{
						Kind:       ClusterKind,
						APIVersion: apiVersion,
					},
					Metadata: Metadata{
						Name: "Test",
					},
				}

				data, _ := yaml.Marshal(manifest)

				return data
			}(),
			expectKind: ClusterKind,
		},
		{
			name: "Should successfully identify application",

			withManifest: func() []byte {
				manifest := Application{
					TypeMeta: TypeMeta{
						Kind:       ApplicationKind,
						APIVersion: apiVersion,
					},
					Metadata: Metadata{
						Name: "Test",
					},
				}

				data, _ := yaml.Marshal(manifest)

				return data
			}(),
			expectKind: ApplicationKind,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewReader(tc.withManifest)

			kind, err := InferKindFromManifest(buf)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectKind, kind)
		})
	}
}
