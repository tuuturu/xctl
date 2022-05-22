package environment

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/config"
	"github.com/deifyed/xctl/pkg/tools/github"
	"github.com/deifyed/xctl/pkg/tools/logging"
	"github.com/deifyed/xctl/pkg/tools/secrets"
	"github.com/deifyed/xctl/pkg/tools/secrets/keyring"
	githubSDK "github.com/google/go-github/v44/github"
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

func Authenticate(manifest *v1alpha1.Environment) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		var (
			secretsClient = keyring.Client{ClusterName: manifest.Metadata.Name}
			log           = logging.GetLogger("environment", "authenticate")
		)

		err := handleGithub(cmd.Context(), log, secretsClient)
		if err != nil {
			return fmt.Errorf("handling Github authentication: %w", err)
		}

		return nil
	}
}

var errInvalidAccessToken = errors.New("invalid token")

func handleGithub(ctx context.Context, log logging.Logger, secretsClient keyring.Client) error {
	log.Debug("Checking for existing token")

	accessToken, err := secretsClient.Get(
		config.DefaultSecretsGithubNamespace,
		config.DefaultSecretsGithubAccessTokenKey,
	)
	if err == nil {
		err = verifyToken(ctx, accessToken)
		if err == nil {
			log.Debug("Found valid token")

			return nil
		}
	}

	if !errors.Is(err, secrets.ErrNotFound) && !errors.Is(err, errInvalidAccessToken) {
		return fmt.Errorf("querying for access token: %w", err)
	}

	log.Debug("Missing or invalid token found, proceeding with authentication")

	accessToken, err = authenticateWithGithub()
	if err != nil {
		return fmt.Errorf("authenticating with Github: %w", err)
	}

	log.Debug("Authentication success. Storing token")

	err = secretsClient.Put(
		config.DefaultSecretsGithubNamespace,
		map[string]string{config.DefaultSecretsGithubAccessTokenKey: accessToken},
	)
	if err != nil {
		return fmt.Errorf("storing Github access token: %w", err)
	}

	return nil
}

func authenticateWithGithub() (string, error) {
	client := http.Client{}

	deviceCodeResponse, err := github.RequestDeviceCode(client, config.DefaultGithubOAuthClientID)
	if err != nil {
		return "", fmt.Errorf("requesting device code: %w", err)
	}

	println(fmt.Sprintf(
		"%s Enter the following code when prompted: %s",
		aurora.Yellow("Attention!"),
		aurora.Green(deviceCodeResponse.UserCode),
	))

	if !proceedPrompt() {
		return "", errors.New("aborted by user")
	}

	err = openBrowser(deviceCodeResponse.VerificationURI)
	if err != nil {
		return "", fmt.Errorf("opening verification URI: %w", err)
	}

	accessToken, err := github.PollForAccessToken(client, config.DefaultGithubOAuthClientID, deviceCodeResponse)
	if err != nil {
		return "", fmt.Errorf("polling for access token: %w", err)
	}

	return accessToken, nil
}

func verifyToken(ctx context.Context, token string) error {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	oauth2HTTPClient := oauth2.NewClient(ctx, tokenSource)

	client := githubSDK.NewClient(oauth2HTTPClient)

	_, response, err := client.Repositories.ListAll(ctx, &githubSDK.RepositoryListAllOptions{})
	if err != nil {
		return fmt.Errorf("listing repositories: %w", err)
	}

	if response.StatusCode == http.StatusUnauthorized || response.StatusCode == http.StatusForbidden {
		return errInvalidAccessToken
	}

	return nil
}

func proceedPrompt() bool {
	print("Ready? [y/N] ")

	var confirmation string
	fmt.Scanln(&confirmation)

	return strings.ToLower(confirmation) == "y"
}

func openBrowser(url string) (err error) {
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		return fmt.Errorf("opening browser: %w", err)
	}

	return nil
}
