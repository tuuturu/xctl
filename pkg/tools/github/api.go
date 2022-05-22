package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/deifyed/xctl/pkg/cloud"
	"github.com/deifyed/xctl/pkg/config"
	"github.com/deifyed/xctl/pkg/tools/secrets"
	"github.com/google/go-github/v44/github"
	"github.com/logrusorgru/aurora/v3"
	"golang.org/x/oauth2"
)

func Authenticator() cloud.AuthenticationService {
	return &authenticationService{}
}

func (a *authenticationService) AuthenticationFlow(secretsClient secrets.Client, userInputPrompter cloud.UserInputPrompter) error {
	httpClient := http.Client{}

	deviceCodeResponse, err := requestDeviceCode(httpClient, config.DefaultGithubOAuthClientID)
	if err != nil {
		return fmt.Errorf("requesting device code: %w", err)
	}

	println(fmt.Sprintf(
		"%s Enter the following code when prompted: %s",
		aurora.Yellow("Attention!"),
		aurora.Green(deviceCodeResponse.UserCode),
	))

	response := userInputPrompter("Ready? [y/N] ", false)

	if strings.ToLower(response) != "y" {
		return errors.New("user aborted")
	}

	err = openBrowser(deviceCodeResponse.VerificationURI)
	if err != nil {
		return fmt.Errorf("opening verification URI: %w", err)
	}

	accessToken, err := pollForAccessToken(httpClient, config.DefaultGithubOAuthClientID, deviceCodeResponse)
	if err != nil {
		return fmt.Errorf("polling for access token: %w", err)
	}

	err = secretsClient.Put(
		config.DefaultSecretsGithubNamespace,
		map[string]string{config.DefaultSecretsGithubAccessTokenKey: accessToken},
	)
	if err != nil {
		return fmt.Errorf("storing access token: %w", err)
	}

	return nil
}

func (a *authenticationService) Authenticate(secretsClient secrets.Client) error {
	accessToken, err := secretsClient.Get(config.DefaultSecretsGithubNamespace, config.DefaultSecretsGithubAccessTokenKey)
	if err != nil {
		return fmt.Errorf("retrieving credentials: %w", err)
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})

	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{Source: tokenSource},
	}

	a.client = github.NewClient(oauth2Client)

	return nil
}

func (a *authenticationService) ValidateAuthentication(ctx context.Context) error {
	_, response, err := a.client.Repositories.ListAll(ctx, &github.RepositoryListAllOptions{})
	if err != nil {
		return fmt.Errorf("listing repositories: %w", err)
	}

	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
		return errors.New("invalid access token")
	}

	return nil
}
