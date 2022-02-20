package binary

import (
	"bytes"
	"path"
	"testing"

	"github.com/deifyed/xctl/pkg/config"
	"github.com/spf13/afero"

	"github.com/deifyed/xctl/pkg/tools/clients/vault"

	"github.com/stretchr/testify/assert"
)

func TestDownloadBinary(t *testing.T) {
	t.Skipf("skipping due to actual download. should be ran after bumping and/or in CI")
	t.Parallel()

	fs := &afero.Afero{Fs: afero.NewMemMapFs()}

	actualPath, err := getVaultPath(fs)
	assert.NoError(t, err)

	binariesDir, err := config.GetAbsoluteBinariesDir()
	assert.NoError(t, err)

	expectedPath := path.Join(binariesDir, "vault", version, "vault")

	assert.Equal(t, expectedPath, actualPath)
}

func TestParsing(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		withOutput string
		expect     vault.InitializationResponse
	}{
		{
			name:       "Should work",
			withOutput: dummyOutput,
			expect: vault.InitializationResponse{
				RootToken: "s.CCo5VpuY1S9qHOZbN2n9eHs2",
				UnsealKeysB64: []string{
					"GiT/H48inyWI1Y7kyW7fq7nX37iAkwo/iAolQNUWkdnx",
					"J3fOIecFDXtbGLbe93/oN5vLpvjXuZj703Pmt5Q5MPXN",
					"5/SgmglMsoQAtEThFvrx9CqXVs/IfZ/lJCeyT5cCpLVm",
					"C9p2QR/oOpP5xtswWeuNf1V0dJNRMPfmFuCRNwReNd5H",
					"pF5blqYQQe45idTsbAJKTK3jkdNqFQlhzwI6tW7WRi6W",
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			response, err := parseInitializationResponse(bytes.NewBufferString(tc.withOutput))
			assert.NoError(t, err)

			assert.Equal(t, tc.expect, response)
		})
	}
}

const dummyOutput = `
{
  "unseal_keys_b64": [
    "GiT/H48inyWI1Y7kyW7fq7nX37iAkwo/iAolQNUWkdnx",
    "J3fOIecFDXtbGLbe93/oN5vLpvjXuZj703Pmt5Q5MPXN",
    "5/SgmglMsoQAtEThFvrx9CqXVs/IfZ/lJCeyT5cCpLVm",
    "C9p2QR/oOpP5xtswWeuNf1V0dJNRMPfmFuCRNwReNd5H",
    "pF5blqYQQe45idTsbAJKTK3jkdNqFQlhzwI6tW7WRi6W"
  ],
  "unseal_keys_hex": [
    "1a24ff1f8f229f2588d58ee4c96edfabb9d7dfb880930a3f880a2540d51691d9f1",
    "2777ce21e7050d7b5b18b6def77fe8379bcba6f8d7b998fbd373e6b7943930f5cd",
    "e7f4a09a094cb28400b444e116faf1f42a9756cfc87d9fe52427b24f9702a4b566",
    "0bda76411fe83a93f9c6db3059eb8d7f557474935130f7e616e09137045e35de47",
    "a45e5b96a61041ee3989d4ec6c024a4cade391d36a150961cf023ab56ed6462e96"
  ],
  "unseal_shares": 5,
  "unseal_threshold": 3,
  "recovery_keys_b64": [],
  "recovery_keys_hex": [],
  "recovery_keys_shares": 5,
  "recovery_keys_threshold": 3,
  "root_token": "s.CCo5VpuY1S9qHOZbN2n9eHs2"
}
`
