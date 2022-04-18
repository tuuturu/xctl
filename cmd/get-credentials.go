package cmd

import (
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/deifyed/xctl/cmd/helpers"

	"github.com/deifyed/xctl/pkg/tools/i18n"

	"github.com/deifyed/xctl/pkg/tools/secrets"

	"github.com/deifyed/xctl/cmd/hooks"
	kubectlBinary "github.com/deifyed/xctl/pkg/tools/clients/kubectl/binary"
	vaultBinary "github.com/deifyed/xctl/pkg/tools/clients/vault/binary"

	"github.com/deifyed/xctl/pkg/apis/xctl"
	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/cloud/linode"
	vaultPlugin "github.com/deifyed/xctl/pkg/plugins/vault"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	getCredentialsCmdOpts = getCredentialsOpts{ //nolint:gochecknoglobals
		io: xctl.IOStreams{
			In:  os.Stdin,
			Out: os.Stdout,
			Err: os.Stderr,
		},
		fs: &afero.Afero{Fs: afero.NewOsFs()},
	}
	getCredentialsCmd = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "credentials",
		Short: i18n.T("cmdGetCredentialsShortDescription"),
		Args:  cobra.ExactArgs(1),
		PreRunE: hooks.EnvironmentManifestInitializer(hooks.EnvironmentManifestInitializerOpts{
			Io:                  getCredentialsCmdOpts.io,
			Fs:                  getCredentialsCmdOpts.fs,
			EnvironmentManifest: &getCredentialsCmdOpts.environmentManifest,
			SourcePath:          &getCredentialsCmdOpts.environmentManifestPath,
		}),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]

			if target != "grafana" {
				return fmt.Errorf("credentials for %s not found", target)
			}

			provider := linode.NewLinodeProvider()

			err := provider.Authenticate()
			if err != nil {
				return fmt.Errorf("authenticating with cloud provider: %w", err)
			}

			ctx := context.Background()

			kubeConfigPath, err := ensureKubeConfig(ensureKubeConfigOpts{
				fs:                  getCredentialsCmdOpts.fs,
				ctx:                 ctx,
				provider:            provider,
				environmentManifest: getCredentialsCmdOpts.environmentManifest,
			})
			if err != nil {
				return fmt.Errorf("setting up kubeconfig: %w", err)
			}

			kubectlClient, err := kubectlBinary.New(forwardCmdOpts.fs, kubeConfigPath)
			if err != nil {
				return fmt.Errorf("acquiring Kubernetes client: %w", err)
			}

			stopFn, err := kubectlClient.PortForward(vaultPlugin.PortForwardOpts())
			if err != nil {
				return fmt.Errorf("forwarding vault: %w", err)
			}

			defer func() {
				_ = stopFn()
			}()

			var secretsClient secrets.Client

			secretsClient, err = vaultBinary.New(getCredentialsCmdOpts.fs)
			if err != nil {
				return fmt.Errorf("preparing vault binary: %w", err)
			}

			username, err := secretsClient.Get("grafana", "adminUsername")
			if err != nil {
				return fmt.Errorf("retrieving username: %w", err)
			}

			if username == "" {
				username = "N/A"
			}

			password, err := secretsClient.Get("grafana", "adminPassword")
			if err != nil {
				return fmt.Errorf("retrieving password: %w", err)
			}

			if password == "" {
				password = "N/A"
			}

			t, err := template.New("credentials").Parse(getCredentialsTemplate)
			if err != nil {
				return fmt.Errorf("parsing template: %w", err)
			}

			err = t.Execute(getCredentialsCmdOpts.io.Out, struct {
				ServiceName string
				Username    string
				Password    string
			}{
				ServiceName: target,
				Username:    username,
				Password:    password,
			})
			if err != nil {
				return fmt.Errorf("printing credentials: %w", err)
			}

			return nil
		},
	}
)

const getCredentialsTemplate = `
{{ .ServiceName }} credentials
	Username: {{ .Username }}
	Password: {{ .Password }}
`

//nolint:gochecknoinits
func init() {
	helpers.AddEnvironmentContextFlag(getCredentialsCmd.Flags(), &getCredentialsCmdOpts.environmentManifestPath)

	getCmd.AddCommand(getCredentialsCmd)
}

type getCredentialsOpts struct {
	io                      xctl.IOStreams
	fs                      *afero.Afero
	environmentManifestPath string
	environmentManifest     v1alpha1.Environment
}
