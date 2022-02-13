package reconciliation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestAssertDependencyExistence(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		withTests  []DependencyTestFn
		withExpect bool

		expectResult bool
		expectErr    string
	}{
		{
			name: "Should return true if expectence is true all tests are true",
			withTests: []DependencyTestFn{
				func() (bool, error) { return true, nil },
				func() (bool, error) { return true, nil },
				func() (bool, error) { return true, nil },
			},
			withExpect:   true,
			expectResult: true,
		},
		{
			name: "Should return false if expectence is true and one of the tests are false",
			withTests: []DependencyTestFn{
				func() (bool, error) { return true, nil },
				func() (bool, error) { return false, nil },
				func() (bool, error) { return true, nil },
			},
			withExpect:   true,
			expectResult: false,
		},
		{
			name: "Should return true if expectence is false all tests are false",
			withTests: []DependencyTestFn{
				func() (bool, error) { return false, nil },
				func() (bool, error) { return false, nil },
				func() (bool, error) { return false, nil },
			},
			withExpect:   false,
			expectResult: true,
		},
		{
			name: "Should return false if expectence is false and one of the tests are true",
			withTests: []DependencyTestFn{
				func() (bool, error) { return false, nil },
				func() (bool, error) { return true, nil },
				func() (bool, error) { return false, nil },
			},
			withExpect:   false,
			expectResult: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			result, err := AssertDependencyExistence(tc.withExpect, tc.withTests...)

			if tc.expectErr != "" {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectErr, err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectResult, result)
		})
	}
}
