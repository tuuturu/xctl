package vault

import "net/url"

const DefaultPort = 8200

type InitializationResponse struct {
	RootToken     string   `json:"root_token"`
	UnsealKeysB64 []string `json:"unseal_keys_b64"`
}

type ConfigureKubernetesAuthenticationOpts struct {
	Host             url.URL
	TokenReviewerJWT string
	CACert           string
	Issuer           url.URL
}

type Operator interface {
	Initialize() (InitializationResponse, error)
	Unseal(key string) error
}

type Auth interface {
	EnableKubernetesAuthentication() error
	ConfigureKubernetesAuthentication(ConfigureKubernetesAuthenticationOpts) error
}

type Client interface {
	Auth
	Operator
	SetToken(token string)
}
