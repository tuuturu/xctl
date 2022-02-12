package secrets

// Client defines operations available for a secrets engine
type Client interface {
	// Put knows how to store a named secret containing key/value pairs
	Put(name string, secrets map[string]string) error
	// Get knows how to retrieve a secret attribute
	Get(name string, key string) (string, error)
}
