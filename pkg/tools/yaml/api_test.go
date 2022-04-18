package yaml

import (
	"bytes"
	_ "embed"
	"io"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/environment-example.yaml
var environmentExample []byte

func TestRemoveComments(t *testing.T) {
	testCases := []struct {
		name        string
		withContent []byte
	}{
		{
			name:        "Should print out a clean and pretty env example",
			withContent: environmentExample,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cleanStream := RemoveComments(bytes.NewReader(tc.withContent))

			result, err := io.ReadAll(cleanStream)
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.name, result)
		})
	}
}
