package environment

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"syscall"

	"github.com/deifyed/xctl/pkg/tools/secrets"

	"github.com/deifyed/xctl/pkg/tools/github"
	"golang.org/x/term"

	"github.com/deifyed/xctl/pkg/cloud/linode"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/cloud"
	"github.com/deifyed/xctl/pkg/tools/logging"
	"github.com/deifyed/xctl/pkg/tools/secrets/keyring"
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
)

func Authenticate(manifest *v1alpha1.Environment) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		var (
			secretsClient = keyring.Client{EnvironmentName: manifest.Metadata.Name}
			log           = logging.GetLogger("environment", "authenticate")
		)

		err := handleProvider(cmd.Context(), log, secretsClient, manifest.Spec.Provider)
		if err != nil {
			return fmt.Errorf("handling cloud provider authentication: %w", err)
		}

		successPrint(cmd.OutOrStdout(), strings.Title(manifest.Spec.Provider))

		err = handleProvider(cmd.Context(), log, secretsClient, githubProvider)
		if err != nil {
			return fmt.Errorf("handling Github authentication: %w", err)
		}

		successPrint(cmd.OutOrStdout(), "Github")

		return nil
	}
}

func handleProvider(ctx context.Context, log logging.Logger, secretsClient keyring.Client, providerName string) error {
	log.Debug("Checking for existing cloud provider credentials")

	var provider cloud.AuthenticationService

	switch providerName {
	case "linode":
		provider = linode.NewLinodeProvider()
	case githubProvider:
		provider = github.Authenticator()
	default:
		return fmt.Errorf("unknown cloud provider %s", providerName)
	}

	err := provider.Authenticate(secretsClient)
	if err != nil {
		if !errors.Is(err, secrets.ErrNotFound) {
			return fmt.Errorf("authenticating with existing credentials: %w", err)
		}
	} else {
		err = provider.ValidateAuthentication(ctx)
		if err == nil {
			return nil
		}

		if err != nil && !errors.Is(err, cloud.ErrNotAuthenticated) {
			return fmt.Errorf("validating existing credentials: %w", err)
		}
	}

	err = provider.AuthenticationFlow(secretsClient, prompter)
	if err != nil {
		return fmt.Errorf("executing cloud provider authentication flow: %w", err)
	}

	err = provider.Authenticate(secretsClient)
	if err != nil {
		return fmt.Errorf("authenticating with provider: %w", err)
	}

	err = provider.ValidateAuthentication(ctx)
	if err != nil {
		return fmt.Errorf("validating credentials: %w", err)
	}

	return nil
}

func prompter(msg string, hidden bool) string {
	fmt.Print(msg)

	var result string

	if hidden {
		rawResult, _ := term.ReadPassword(syscall.Stdin)

		result = string(rawResult)
	} else {
		fmt.Scanln(&result)
	}

	fmt.Print("\n")

	return result
}

func successPrint(out io.Writer, name string) {
	fmt.Fprintf(out, "[%s] %s\n", name, aurora.Green("OK"))
}

const githubProvider = "github"
