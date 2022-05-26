package namespace

import (
	"testing"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/tools/reconciliation"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestReconciler_Reconcile(t *testing.T) {
	testCases := []struct {
		name          string
		withNamespace string
	}{
		{
			name:          "Should generate correct file with correct contents with valid inputs",
			withNamespace: "mock-namespace",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fs := &afero.Afero{Fs: afero.NewMemMapFs()}
			r := Reconciler("/infrastructure/mock-cluster")

			_, err := r.Reconcile(reconciliation.Context{
				Filesystem:          fs,
				RootDirectory:       "/",
				EnvironmentManifest: v1alpha1.Environment{},
				ApplicationDeclaration: v1alpha1.Application{
					Metadata: v1alpha1.Metadata{Namespace: tc.withNamespace},
				},
			})
			assert.NoError(t, err)

			raw, err := fs.ReadFile("/infrastructure/mock-cluster/argocd/namespaces/mock-namespace.yaml")
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.name, raw)
		})
	}
}
