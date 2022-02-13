package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/deifyed/xctl/pkg/plugins/vault"

	"github.com/deifyed/xctl/cmd/preruns"
	"github.com/deifyed/xctl/pkg/apis/xctl"
	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/clients/kubectl"
	kubectlBinary "github.com/deifyed/xctl/pkg/clients/kubectl/binary"
	vaultClient "github.com/deifyed/xctl/pkg/clients/vault"
	"github.com/deifyed/xctl/pkg/cloud/linode"
	"github.com/deifyed/xctl/pkg/plugins/grafana"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type forwardOpts struct {
	io                  xctl.IOStreams
	fs                  *afero.Afero
	clusterManifestPath string
	ClusterManifest     v1alpha1.Cluster
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
		Short: "Sets up a connection to a service",
		Args:  cobra.ExactArgs(1),
		PreRunE: preruns.ClusterManifestIniter(preruns.ClusterManifestIniterOpts{
			Io:              forwardCmdOpts.io,
			Fs:              forwardCmdOpts.fs,
			ClusterManifest: &forwardCmdOpts.ClusterManifest,
			SourcePath:      &forwardCmdOpts.clusterManifestPath,
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
				fs:              forwardCmdOpts.fs,
				ctx:             ctx,
				provider:        provider,
				clusterManifest: forwardCmdOpts.ClusterManifest,
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

func init() {
	flags := forwardCmd.Flags()

	flags.StringVarP(
		&forwardCmdOpts.clusterManifestPath,
		"cluster-declaration",
		"c",
		"-",
		"cluster declaration representing context of the virtual environment",
	)

	rootCmd.AddCommand(forwardCmd)
}
