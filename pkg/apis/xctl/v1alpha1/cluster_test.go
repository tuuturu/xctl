package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponentNamer(t *testing.T) {
	testCases := []struct {
		name              string
		withClusterName   string
		withComponentType string
		withComponentID   string
		expectName        string
	}{
		{
			name:              "Should with a basic case",
			withClusterName:   "test",
			withComponentType: "loadbalancer",
			withComponentID:   "uxxv",
			expectName:        "xctl-test-loadbalancer-uxxv",
		},
		{
			name:              "Should ensure lower case",
			withClusterName:   "test2",
			withComponentType: "domain",
			withComponentID:   "TUUTURUORG",
			expectName:        "xctl-test2-domain-tuuturuorg",
		},
		{
			name:              "Works without ID",
			withClusterName:   "superprod",
			withComponentType: "cluster",
			withComponentID:   "",
			expectName:        "xctl-superprod-cluster",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			manifest := NewDefaultCluster()
			manifest.Metadata.Name = tc.withClusterName

			name := manifest.ComponentName(tc.withComponentType, tc.withComponentID)

			assert.Equal(t, tc.expectName, name)
		})
	}
}
