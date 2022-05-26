package argocd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_Name(t *testing.T) {
	testCases := []struct {
		name       string
		withURL    string
		expectName string
	}{
		{
			name:       "Should return the correct name",
			withURL:    "git@github.com:tuuturu/xctl.git",
			expectName: "xctl",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := repository{URL: tc.withURL}

			assert.Equal(t, tc.expectName, repo.Name())
		})
	}
}

func TestRepository_Owner(t *testing.T) {
	testCases := []struct {
		name        string
		withURL     string
		expectOwner string
	}{
		{
			name:        "Should return the correct owner",
			withURL:     "git@github.com:tuuturu/xctl.git",
			expectOwner: "tuuturu",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := repository{URL: tc.withURL}

			assert.Equal(t, tc.expectOwner, repo.Owner())
		})
	}
}
