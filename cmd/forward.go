package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/deifyed/xctl/pkg/plugins/argocd"

	"github.com/logrusorgru/aurora/v3"

	"github.com/deifyed/xctl/cmd/helpers"

	"github.com/deifyed/xctl/pkg/plugins/grafana"

	"github.com/deifyed/xctl/pkg/tools/i18n"

	"github.com/deifyed/xctl/cmd/hooks"
	"github.com/deifyed/xctl/pkg/tools/clients/kubectl"
	kubectlBinary "github.com/deifyed/xctl/pkg/tools/clients/kubectl/binary"

	"github.com/deifyed/xctl/pkg/apis/xctl"
	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/cloud/linode"
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
				portForwardOpts = grafana.PortForwardOpts()
			case "argocd":
				portForwardOpts = argocd.PortForwardOpts()
			default:
				return fmt.Errorf("service %s not found", target)
			}

			provider := linode.NewLinodeProvider()

			err := provider.Authenticate()
			if err != nil {
				return fmt.Errorf("authenticating with cloud provider: %w", err)
			}

			kubeConfigPath, err := ensureKubeConfig(ensureKubeConfigOpts{
				fs:                  forwardCmdOpts.fs,
				ctx:                 cmd.Context(),
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
				"Serving %s at %s\n",
				aurora.Green(strings.Title(target)),
				aurora.Green(fmt.Sprintf("http://localhost:%d", portForwardOpts.LocalPort)),
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

//nolint:gochecknoinits
func init() {
	helpers.AddEnvironmentContextFlag(forwardCmd.Flags(), &forwardCmdOpts.environmentManifestPath)

	rootCmd.AddCommand(forwardCmd)
}
