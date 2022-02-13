package binary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsConnectionRefused(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		withMessage  string
		expectResult bool
	}{
		{
			name: "Should return true on a connection refused error",

			withMessage:  `Error initializing: Put \"http://127.0.0.1:8200/v1/sys/init\": dial tcp 127.0.0.1:8200: connect: connection refused\n`, //nolint:lll
			expectResult: true,
		},
		{
			name: "Should return false on a random message",

			withMessage:  "somethingSuperimpressive and more than Something CRAZY",
			expectResult: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectResult, isConnectionRefused(tc.withMessage))
		})
	}
}
