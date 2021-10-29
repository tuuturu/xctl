package vault

const DefaultPort = 8200

type InitializationResponse struct {
	Token string
	Keys  []string
}

type Client interface {
	Initialize() (InitializationResponse, error)
	SetToken(token string)
	Unseal(key string) error
}
