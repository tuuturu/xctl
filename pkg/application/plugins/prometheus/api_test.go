package prometheus

import (
	"fmt"
	"testing"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/tools/reconciliation"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestReconciler_Reconcile(t *testing.T) {
	testCases := []struct {
		name    string
		withApp v1alpha1.Application
	}{
		{
			name: "Should generate correct file with correct contents with valid inputs",
			withApp: v1alpha1.Application{
				Metadata: v1alpha1.Metadata{Name: "mock-app", Namespace: "mock-namespace"},
				Spec:     v1alpha1.ApplicationSpec{Metrics: "/metrics"},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fs := &afero.Afero{Fs: afero.NewMemMapFs()}
			r := Reconciler(fmt.Sprintf("/infrastructure/applications/%s", tc.withApp.Metadata.Name))

			_, err := r.Reconcile(reconciliation.Context{
				Filesystem:             fs,
				RootDirectory:          "/",
				ApplicationDeclaration: tc.withApp,
			})
			assert.NoError(t, err)

			raw, err := fs.ReadFile(fmt.Sprintf(
				"/infrastructure/applications/%s/base/service-monitor.yaml",
				tc.withApp.Metadata.Name,
			))
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.name, raw)
		})
	}
}
