package vault

const DefaultPort = 8200

type InitializationResponse struct {
	RootToken     string   `json:"root_token"`
	UnsealKeysB64 []string `json:"unseal_keys_b64"`
}

type Client interface {
	Initialize() (InitializationResponse, error)
	SetToken(token string)
	Unseal(key string) error
}
