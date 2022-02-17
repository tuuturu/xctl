package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/deifyed/xctl/pkg/tools/i18n"

	"github.com/deifyed/xctl/cmd/hooks"
	"github.com/deifyed/xctl/pkg/tools/clients/kubectl"
	kubectlBinary "github.com/deifyed/xctl/pkg/tools/clients/kubectl/binary"
	vaultClient "github.com/deifyed/xctl/pkg/tools/clients/vault"

	"github.com/deifyed/xctl/pkg/plugins/vault"

	"github.com/deifyed/xctl/pkg/apis/xctl"
	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/cloud/linode"
	"github.com/deifyed/xctl/pkg/plugins/grafana"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type forwardOpts struct {
	io                      xctl.IOStreams
	fs                      *afero.Afero
	environmentManifestPath string
	EnvironmentManifest     v1alpha1.Environment
}

var (
	forwardCmdOpts = forwardOpts{ //nolint:gochecknoglobals
		io: xctl.IOStreams{
			In:  os.Stdin,
			Out: os.Stdout,
			Err: os.Stderr,
		},
		fs: &afero.Afero{Fs: afero.NewOsFs()},
	}
	forwardCmd = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "forward",
		Short: i18n.T("cmdForwardShortDescription"),
		Args:  cobra.ExactArgs(1),
		PreRunE: hooks.EnvironmentManifestInitializer(hooks.EnvironmentManifestInitializerOpts{
			Io:                  forwardCmdOpts.io,
			Fs:                  forwardCmdOpts.fs,
			EnvironmentManifest: &forwardCmdOpts.EnvironmentManifest,
			SourcePath:          &forwardCmdOpts.environmentManifestPath,
		}),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]

			var portForwardOpts kubectl.PortForwardOpts

			switch target {
			case "grafana":
				portForwardOpts = getGrafanaForwardOpts()
			case "vault":
				portForwardOpts = getVaultForwardOpts()
			default:
				return fmt.Errorf("service %s not found", target)
			}

			provider := linode.NewLinodeProvider()

			err := provider.Authenticate()
			if err != nil {
				return fmt.Errorf("authenticating with cloud provider: %w", err)
			}

			ctx := context.Background()

			kubeConfigPath, err := ensureKubeConfig(ensureKubeConfigOpts{
				fs:                  forwardCmdOpts.fs,
				ctx:                 ctx,
				provider:            provider,
				environmentManifest: forwardCmdOpts.EnvironmentManifest,
			})
			if err != nil {
				return fmt.Errorf("setting up kubeconfig: %w", err)
			}

			kubectlClient, err := kubectlBinary.New(forwardCmdOpts.fs, kubeConfigPath)
			if err != nil {
				return fmt.Errorf("acquiring Kubernetes client: %w", err)
			}

			stopFn, err := kubectlClient.PortForward(portForwardOpts)
			if err != nil {
				return fmt.Errorf("port forwarding: %w", err)
			}

			_, _ = fmt.Fprintf(forwardCmdOpts.io.Out,
				"Serving %s at http://localhost:%d\n",
				target,
				portForwardOpts.LocalPort,
			)

			var (
				c       = make(chan os.Signal, 1)
				running = true
			)

			signal.Notify(c, os.Interrupt)

			go func() {
				<-c

				_, _ = fmt.Fprintf(forwardCmdOpts.io.Out, "Shutting down %s forward\n", target)

				_ = stopFn()

				running = false
			}()

			for running {
				time.Sleep(1 * time.Second)
			}

			return nil
		},
	}
)

const (
	grafanaPort      = 80
	grafanaLocalPort = 8000
)

func getGrafanaForwardOpts() kubectl.PortForwardOpts {
	plugin, _ := grafana.NewPlugin(grafana.NewPluginOpts{})

	return kubectl.PortForwardOpts{
		Service: kubectl.Service{
			Name:      plugin.Metadata.Name,
			Namespace: plugin.Metadata.Namespace,
		},
		ServicePort: grafanaPort,
		LocalPort:   grafanaLocalPort,
	}
}

func getVaultForwardOpts() kubectl.PortForwardOpts {
	plugin := vault.NewVaultPlugin()

	return kubectl.PortForwardOpts{
		Service: kubectl.Service{
			Name:      plugin.Metadata.Name,
			Namespace: plugin.Metadata.Namespace,
		},
		ServicePort: vaultClient.DefaultPort,
		LocalPort:   vaultClient.DefaultPort,
	}
}

//nolint:gochecknoinits
func init() {
	flags := forwardCmd.Flags()

	flags.StringVarP(
		&forwardCmdOpts.environmentManifestPath,
		i18n.T("cmdFlagContextName"),
		"c",
		"-",
		i18n.T("cmdFlagContextUsage"),
	)

	rootCmd.AddCommand(forwardCmd)
}
