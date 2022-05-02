package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/deifyed/xctl/pkg/plugins/argocd"

	"github.com/deifyed/xctl/pkg/config"
	"github.com/deifyed/xctl/pkg/plugins/grafana"
	kubectlBinary "github.com/deifyed/xctl/pkg/tools/clients/kubectl/binary"

	"github.com/deifyed/xctl/cmd/helpers"

	"github.com/deifyed/xctl/pkg/tools/i18n"

	"github.com/deifyed/xctl/cmd/hooks"

	"github.com/deifyed/xctl/pkg/apis/xctl"
	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/cloud/linode"
	"github.com/logrusorgru/aurora/v3"
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
		Use:    "credentials",
		Short:  i18n.T("cmdGetCredentialsShortDescription"),
		Args:   cobra.ExactArgs(1),
		Hidden: true,
		PreRunE: hooks.EnvironmentManifestInitializer(hooks.EnvironmentManifestInitializerOpts{
			Io:                  getCredentialsCmdOpts.io,
			Fs:                  getCredentialsCmdOpts.fs,
			EnvironmentManifest: &getCredentialsCmdOpts.environmentManifest,
			SourcePath:          &getCredentialsCmdOpts.environmentManifestPath,
		}),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]

			provider := linode.NewLinodeProvider()

			err := provider.Authenticate()
			if err != nil {
				return fmt.Errorf("authenticating with cloud provider: %w", err)
			}

			kubeconfigPath, err := config.GetAbsoluteKubeconfigPath(getCredentialsCmdOpts.environmentManifest.Metadata.Name)
			if err != nil {
				return fmt.Errorf("acquiring kubeconfig path: %w", err)
			}

			kubectlClient, err := kubectlBinary.New(getCredentialsCmdOpts.fs, kubeconfigPath)
			if err != nil {
				return fmt.Errorf("acquiring kubectl client: %w", err)
			}

			var credentials v1alpha1.PluginCredentials

			switch target {
			case "grafana":
				credentials, err = grafana.Credentials(kubectlClient)
			case "argocd":
				credentials, err = argocd.Credentials(kubectlClient)
			default:
				return fmt.Errorf("credentials for %s not found", target)
			}

			if err != nil {
				return fmt.Errorf("acquiring Grafana credentials: %w", err)
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
				Username:    fmt.Sprint(aurora.Green(credentials.Username)),
				Password:    fmt.Sprint(aurora.Green(credentials.Password)),
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
