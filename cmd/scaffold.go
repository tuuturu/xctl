package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/deifyed/xctl/pkg/application"

	"github.com/deifyed/xctl/pkg/tools/yaml"

	"github.com/deifyed/xctl/pkg/environment"

	"github.com/deifyed/xctl/pkg/tools/i18n"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/spf13/cobra"
)

var (
	scaffoldCmdOpts = scaffoldCmdOptsContainer{
		Raw: false,
	}
	scaffoldCmd = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "scaffold",
		Short: i18n.T("cmdScaffoldShortDescription"),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var output io.Reader
			resource := strings.ToLower(args[0])

			switch resource {
			case strings.ToLower(v1alpha1.EnvironmentKind):
				output = environment.Scaffold()
			case strings.ToLower(v1alpha1.ApplicationKind):
				output = application.Scaffold()
			default:
				return fmt.Errorf("unable to recognize resource type \"%s\". Valid resource types are: %+v",
					resource,
					[]string{strings.ToLower(v1alpha1.EnvironmentKind)},
				)
			}

			if scaffoldCmdOpts.Raw {
				output = yaml.RemoveComments(output)
			}

			_, err := io.Copy(cmd.OutOrStdout(), output)
			if err != nil {
				return fmt.Errorf("scaffolding %s: %w", resource, err)
			}

			return nil
		},
	}
)

//nolint:gochecknoinits
func init() {
	scaffoldCmd.Flags().BoolVarP(&scaffoldCmdOpts.Raw, "raw", "r", false, "scaffold only required fields")

	rootCmd.AddCommand(scaffoldCmd)
}

type scaffoldCmdOptsContainer struct {
	Raw bool
}
