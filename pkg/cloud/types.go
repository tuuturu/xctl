package cloud

import (
	"context"

	"github.com/deifyed/xctl/pkg/tools/secrets"
)

type UserInputPrompter func(message string, hidden bool) (input string)

type AuthenticationService interface {
	// AuthenticationFlow knows how to gather necessary information to authenticate with a cloud provider. It should
	// first check if standard environment variables are available, if not, use the userInputPrompter to prompt for
	// necessary information from the user.
	AuthenticationFlow(client secrets.Client, userInputPrompter UserInputPrompter) error
	// Authenticate knows how to retrieve credentials from a secrets client and use it to authenticate with the cloud
	// provider
	Authenticate(secrets.Client) error
	// ValidateAuthentication knows how to check if authentication is valid
	ValidateAuthentication(context.Context) error
}

// Provider defines required functionality xctl expects from a cloud provider
type Provider interface {
	AuthenticationService
	ClusterService
	DomainService
}
