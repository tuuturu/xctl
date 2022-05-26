package application

import (
	"path"
	"testing"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/tools/reconciliation"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestReconciler_Reconcile(t *testing.T) {
	var testCases = []struct {
		name    string
		withApp v1alpha1.Application
	}{
		{
			name: "Should generate necessary files",
			withApp: v1alpha1.Application{
				TypeMeta: v1alpha1.TypeMeta{
					Kind:       v1alpha1.ApplicationKind,
					APIVersion: "v1alpha1",
				},
				Metadata: v1alpha1.Metadata{
					Name: "mock-app",
				},
				Spec: v1alpha1.ApplicationSpec{
					Image: "xctl.tuuturu.org/mock-app:v0.0.1",
					Port:  "3000",
					Url:   "mock-app.tuuturu.org",
				},
			},
		},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fs := &afero.Afero{Fs: afero.NewMemMapFs()}
			r := reconciler{}

			appDir := path.Join(
				"/",
				config.DefaultInfrastructureDir,
				config.DefaultApplicationsDir,
				tc.withApp.Metadata.Name,
			)

			_, err := r.Reconcile(reconciliation.Context{
				Filesystem:             fs,
				RootDirectory:          "/",
				ApplicationDeclaration: tc.withApp,
			})
			assert.NoError(t, err)

			g := goldie.New(t)

			equalsGoldie(t, g, fs, path.Join(appDir, config.DefaultApplicationBaseDir, "deployment.yaml"), "deployment")
			equalsGoldie(t, g, fs, path.Join(appDir, config.DefaultApplicationBaseDir, "service.yaml"), "service")
			equalsGoldie(t, g, fs, path.Join(appDir, config.DefaultApplicationBaseDir, "ingress.yaml"), "ingress")
		})
	}
}

func equalsGoldie(t *testing.T, g goldie.Tester, fs *afero.Afero, path string, name string) {
	raw, err := fs.ReadFile(path)
	assert.NoError(t, err)

	g.Assert(t, name, raw)
}
