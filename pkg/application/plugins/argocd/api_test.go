package argocd

import (
	"context"
	"fmt"
	"path"
	"testing"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/config"
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
			name: "Should produce a valid application manifest with reasonable input",
			withApp: v1alpha1.Application{
				Metadata: v1alpha1.Metadata{
					Name:      "mock-app",
					Namespace: "mock-namespace",
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fs := &afero.Afero{Fs: afero.NewMemMapFs()}

			environment := v1alpha1.NewDefaultEnvironment()
			environment.Metadata.Name = "mock-cluster"
			environment.Spec.Repository = "git@github.com:mock-org/mock-repo.git"

			absoluteEnvironmentDirectory := path.Join(
				"/",
				config.DefaultInfrastructureDir,
				environment.Metadata.Name,
			)

			r := Reconciler("/infrastructure/mock-cluster")

			_, err := r.Reconcile(reconciliation.Context{
				Ctx:                    context.Background(),
				Filesystem:             fs,
				RootDirectory:          "/",
				EnvironmentManifest:    environment,
				ApplicationDeclaration: tc.withApp,
			})
			assert.NoError(t, err)

			g := goldie.New(t)

			raw, err := fs.ReadFile(path.Join(
				absoluteEnvironmentDirectory,
				"argocd",
				"applications",
				fmt.Sprintf("%s.yaml", tc.withApp.Metadata.Name),
			))
			assert.NoError(t, err)

			g.Assert(t, tc.name, raw)
		})
	}
}
