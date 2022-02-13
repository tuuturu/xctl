package cmd

import "github.com/spf13/cobra"

var getCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "get",
	Short: "query data from the xctl environment",
}

//nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(getCmd)
}
