package cmd

import (
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/deifyed/xctl/cmd/preruns"
	"github.com/deifyed/xctl/pkg/apis/xctl"
	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/clients/kubectl"
	kubectlBinary "github.com/deifyed/xctl/pkg/clients/kubectl/binary"
	"github.com/deifyed/xctl/pkg/clients/vault"
	vaultBinary "github.com/deifyed/xctl/pkg/clients/vault/binary"
	"github.com/deifyed/xctl/pkg/cloud/linode"
	vaultPlugin "github.com/deifyed/xctl/pkg/plugins/vault"
	"github.com/deifyed/xctl/pkg/secrets"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type getCredentialsOpts struct {
	io                  xctl.IOStreams
	fs                  *afero.Afero
	clusterManifestPath string
	ClusterManifest     v1alpha1.Cluster
}

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
		Short: "retrieves credentials",
		Args:  cobra.ExactArgs(1),
		PreRunE: preruns.ClusterManifestIniter(preruns.ClusterManifestIniterOpts{
			Io:              getCredentialsCmdOpts.io,
			Fs:              getCredentialsCmdOpts.fs,
			ClusterManifest: &getCredentialsCmdOpts.ClusterManifest,
			SourcePath:      &getCredentialsCmdOpts.clusterManifestPath,
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
				fs:              getCredentialsCmdOpts.fs,
				ctx:             ctx,
				provider:        provider,
				clusterManifest: getCredentialsCmdOpts.ClusterManifest,
			})
			if err != nil {
				return fmt.Errorf("setting up kubeconfig: %w", err)
			}

			kubectlClient, err := kubectlBinary.New(forwardCmdOpts.fs, kubeConfigPath)
			if err != nil {
				return fmt.Errorf("acquiring Kubernetes client: %w", err)
			}

			vaultPlugin := vaultPlugin.NewVaultPlugin()

			stopFn, err := kubectlClient.PortForward(kubectl.PortForwardOpts{
				Service: kubectl.Service{
					Name:      vaultPlugin.Metadata.Name,
					Namespace: vaultPlugin.Metadata.Namespace,
				},
				ServicePort: vault.DefaultPort,
				LocalPort:   vault.DefaultPort,
			})
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
	flags := getCredentialsCmd.Flags()

	flags.StringVarP(
		&getCredentialsCmdOpts.clusterManifestPath,
		"context",
		"c",
		"-",
		"cluster manifest representing context of the virtual environment",
	)

	getCmd.AddCommand(getCredentialsCmd)
}
