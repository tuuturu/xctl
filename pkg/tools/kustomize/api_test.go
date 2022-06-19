package kustomize

import (
	"bytes"
	"path"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestAddResourceToKustomization(t *testing.T) {
	testCases := []struct {
		name                string
		withAbsoluteWorkDir string
		withResourcePaths   []string
	}{
		{
			name:                "Should work when kustomization does not exist",
			withAbsoluteWorkDir: "/infrastructure/applications/mock-app/base",
			withResourcePaths:   []string{"service-monitor.yaml"},
		},
		{
			name:                "Should add only one upon duplicate adds",
			withAbsoluteWorkDir: "/infrastructure/applications/mock-app/base",
			withResourcePaths:   []string{"service-monitor.yaml", "service-monitor.yaml"},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fs := &afero.Afero{Fs: afero.NewMemMapFs()}

			for _, resource := range tc.withResourcePaths {
				err := AddResourceToKustomization(fs, tc.withAbsoluteWorkDir, resource)
				assert.NoError(t, err)
			}

			g := goldie.New(t)

			rawFile, err := fs.ReadFile(path.Join(tc.withAbsoluteWorkDir, defaultKustomizationFilename))
			assert.NoError(t, err)

			g.Assert(t, tc.name, rawFile)
		})
	}
}

func TestTwoCallsToAdd(t *testing.T) {
	fs := &afero.Afero{Fs: afero.NewMemMapFs()}

	absoluteWorkDir := "/"
	resource := "service-monitor.yaml"

	err := AddResourceToKustomization(fs, absoluteWorkDir, resource)
	assert.NoError(t, err)

	firstRead, err := fs.ReadFile(path.Join(absoluteWorkDir, defaultKustomizationFilename))
	assert.NoError(t, err)

	err = AddResourceToKustomization(fs, absoluteWorkDir, resource)
	assert.NoError(t, err)

	secondRead, err := fs.ReadFile(path.Join(absoluteWorkDir, defaultKustomizationFilename))
	assert.NoError(t, err)

	assert.True(t, bytes.Equal(firstRead, secondRead))
}
