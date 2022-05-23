package argocd

import (
	"fmt"
	"path"
	"testing"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestEstablishConfig(t *testing.T) {
	testCases := []struct {
		name            string
		withRootDir     string
		withEnvironment v1alpha1.Environment
	}{
		{
			name:        "Should create the correct files in the correct folders",
			withRootDir: "/",
			withEnvironment: v1alpha1.Environment{
				Metadata: v1alpha1.Metadata{Name: "mock-env"},
				Spec:     v1alpha1.EnvironmentSpec{Repository: "git@github.com:mockorg/mock.git"},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fs := &afero.Afero{Fs: afero.NewMemMapFs()}

			_, err := establishConfiguration(fs, "/", tc.withEnvironment)
			assert.NoError(t, err)

			cfgDir := path.Join(tc.withRootDir, "infrastructure", tc.withEnvironment.Metadata.Name, "argocd")

			applicationsApplicationsRaw, err := fs.ReadFile(path.Join(cfgDir, "applications.yaml"))
			assert.NoError(t, err)

			applicationsReadmeRaw, err := fs.ReadFile(path.Join(cfgDir, "applications", "README.md"))
			assert.NoError(t, err)

			g := goldie.New(t)

			g.Assert(t, fmt.Sprintf("%s-appsapp.yaml", tc.name), applicationsApplicationsRaw)
			g.Assert(t, fmt.Sprintf("%s-appsreadme.txt", tc.name), applicationsReadmeRaw)
		})
	}
}
