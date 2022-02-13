package cmd

import (
	"github.com/deifyed/xctl/pkg/tools/i18n"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "get",
	Short: i18n.T("cmdGetShortDescription"),
}

//nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(getCmd)
}
