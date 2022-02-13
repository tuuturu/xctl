package hooks

import (
	"bytes"
	"io"
	"testing"

	"github.com/deifyed/xctl/pkg/apis/xctl"
	"github.com/stretchr/testify/assert"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

func TestClusterManifestInitializer(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		withInput      string
		expectManifest v1alpha1.Cluster
	}{
		{
			name:           "Should be equal to default when empty",
			withInput:      "",
			expectManifest: v1alpha1.NewDefaultCluster(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			stdin := bytes.NewReader([]byte(tc.withInput))

			path := "-"
			result := v1alpha1.Cluster{}

			err := ClusterManifestInitializer(ClusterManifestInitializerOpts{
				Io: xctl.IOStreams{
					In:  stdin,
					Out: io.Discard,
					Err: io.Discard,
				},
				Fs:              nil,
				ClusterManifest: &result,
				SourcePath:      &path,
			})(nil, nil)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectManifest, result)
		})
	}
}
