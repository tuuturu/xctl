package linode

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/deifyed/xctl/pkg/cloud"
	"github.com/deifyed/xctl/pkg/config"
	"github.com/deifyed/xctl/pkg/tools/secrets"
	"github.com/linode/linodego"
	"golang.org/x/oauth2"
)

const accessTokenKey = "access-token"

func (p *provider) AuthenticationFlow(secretsClient secrets.Client, userInputPrompter cloud.UserInputPrompter) error {
	accessToken := os.Getenv(linodego.APIEnvVar)

	if accessToken == "" {
		accessToken = userInputPrompter("Please enter your Linode personal access token: ", true)
	}

	err := secretsClient.Put(
		config.DefaultSecretsCloudProviderNamespace,
		map[string]string{accessTokenKey: accessToken},
	)
	if err != nil {
		return fmt.Errorf("storing credentials: %w", err)
	}

	return nil
}

func (p *provider) Authenticate(secretsClient secrets.Client) error {
	accessToken, err := secretsClient.Get(config.DefaultSecretsCloudProviderNamespace, accessTokenKey)
	if err != nil {
		return fmt.Errorf("retrieving credentials: %w", err)
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})

	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{Source: tokenSource},
	}

	p.client = linodego.NewClient(oauth2Client)
	p.client.SetDebug(false)

	return nil
}

func (p *provider) ValidateAuthentication(ctx context.Context) error {
	_, err := p.client.ListInstances(ctx, &linodego.ListOptions{})
	if err != nil {
		return fmt.Errorf("listing instances: %w", err)
	}

	return nil
}
