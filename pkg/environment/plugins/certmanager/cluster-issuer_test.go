package certmanager

import (
	"io"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestClusterIssuer(t *testing.T) {
	testCases := []struct {
		name      string
		withEmail string
	}{
		{
			name:      "Should return expected output",
			withEmail: "mock@email.io",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			issuer, err := newClusterIssuers(tc.withEmail)
			assert.NoError(t, err)

			raw, err := io.ReadAll(issuer)
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.name, raw)
		})
	}
}
