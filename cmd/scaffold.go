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

			switch {
			case contains([]string{v1alpha1.EnvironmentKind, "env"}, resource):
				output = environment.Scaffold()
			case contains([]string{v1alpha1.ApplicationKind, "app"}, resource):
				output = application.Scaffold()
			default:
				return &i18n.HumanReadableError{
					Content: "resource not found",
					Key:     "cmd/scaffold/resourceNotFound",
				}
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

func contains(haystack []string, needle string) bool {
	for _, item := range haystack {
		if strings.EqualFold(item, needle) {
			return true
		}
	}

	return false
}

//nolint:gochecknoinits
func init() {
	scaffoldCmd.Flags().BoolVarP(&scaffoldCmdOpts.Raw, "raw", "r", false, "scaffold only required fields")

	rootCmd.AddCommand(scaffoldCmd)
}

type scaffoldCmdOptsContainer struct {
	Raw bool
}
