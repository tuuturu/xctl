package binary

import (
	"bytes"
	"testing"

	"github.com/deifyed/xctl/pkg/clients/vault"
	"github.com/stretchr/testify/assert"
)

func TestParsing(t *testing.T) {
	testCases := []struct {
		name       string
		withOutput string
		expect     vault.InitializationResponse
	}{
		{
			name:       "Should work",
			withOutput: dummyOutput,
			expect: vault.InitializationResponse{
				Token: "s.wgDHeT2gswN31rznFXURzxwq",
				Keys: []string{
					"QgmZGT7XznOgVcB8eXAF9rbD+7H+4HLEOHTMhLoKeckH",
					"Dy+CKlHP/QeK8I8tG/STLa6XPKewkU/WUwiEo5nOXX+C",
					"P521aeZ+cVLuF1paOezQ+pHY8mg/lYdfcq9c0Uv36rQw",
					"uIXO/UkipDMdp8zJ87Uj96QhY98WIPWSRgpQVCCN5tOP",
					"lGzI21QHGHtUGEx+V0oeSjm27yh5VMRoRGRO9T3dbD4b",
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
Unseal Key 1: QgmZGT7XznOgVcB8eXAF9rbD+7H+4HLEOHTMhLoKeckH
Unseal Key 2: Dy+CKlHP/QeK8I8tG/STLa6XPKewkU/WUwiEo5nOXX+C
Unseal Key 3: P521aeZ+cVLuF1paOezQ+pHY8mg/lYdfcq9c0Uv36rQw
Unseal Key 4: uIXO/UkipDMdp8zJ87Uj96QhY98WIPWSRgpQVCCN5tOP
Unseal Key 5: lGzI21QHGHtUGEx+V0oeSjm27yh5VMRoRGRO9T3dbD4b

Initial Root Token: s.wgDHeT2gswN31rznFXURzxwq

Vault initialized with 5 key shares and a key threshold of 3. Please securely
distribute the key shares printed above. When the Vault is re-sealed,
restarted, or stopped, you must supply at least 3 of these keys to unseal it
before it can start servicing requests.

Vault does not store the generated master key. Without at least 3 keys to
reconstruct the master key, Vault will remain permanently sealed!

It is possible to generate new unseal keys, provided you have a quorum of
existing unseal keys shares. See "vault operator rekey" for more information.
`
