package venv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceAsMap(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		withSlice []string
		expectMap map[string]string
	}{
		{
			name: "Should work",

			withSlice: []string{
				"VAR_A=valueA",
				"VAR_B=valueB",
				"VAR_C=valueC",
			},
			expectMap: map[string]string{
				"VAR_A": "valueA",
				"VAR_B": "valueB",
				"VAR_C": "valueC",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			m := SliceAsMap(tc.withSlice)

			assert.Equal(t, len(tc.expectMap), len(m))

			for key, value := range m {
				assert.Equal(t, tc.expectMap[key], value)
			}

			for key, value := range tc.expectMap {
				assert.Equal(t, value, m[key])
			}
		})
	}
}

func TestMapAsSlice(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		withMap     map[string]string
		expectSlice []string
	}{
		{
			name: "Should work",

			withMap: map[string]string{
				"VAR_A": "valueA",
				"VAR_B": "valueB",
			},
			expectSlice: []string{
				"VAR_A=valueA",
				"VAR_B=valueB",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			s := MapAsSlice(tc.withMap)

			assert.Equal(t, len(tc.expectSlice), len(s))

			for _, item := range tc.expectSlice {
				assert.Contains(t, s, item)
			}

			for _, item := range s {
				assert.Contains(t, tc.expectSlice, item)
			}
		})
	}
}

//nolint:funlen
func TestMergeVariables(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		withSlices  [][]string
		expectSlice []string
	}{
		{
			name: "Should merge multiple slices",

			withSlices: [][]string{
				{
					"A=B",
					"B=C",
					"C=D",
				},
				{
					"D=E",
					"E=F",
					"G=H",
				},
			},
			expectSlice: []string{
				"A=B",
				"B=C",
				"C=D",
				"D=E",
				"E=F",
				"G=H",
			},
		},
		{
			name: "Should prioritize the right most slice",

			withSlices: [][]string{
				{
					"A=B",
					"B=C",
				},
				{
					"A=important",
					"C=D",
				},
			},
			expectSlice: []string{
				"A=important",
				"B=C",
				"C=D",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			merged := MergeVariables(tc.withSlices...)

			assert.Equal(t, len(tc.expectSlice), len(merged))

			for _, item := range tc.expectSlice {
				assert.Contains(t, merged, item)
			}
		})
	}
}
