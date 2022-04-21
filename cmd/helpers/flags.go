package helpers

import (
	"github.com/deifyed/xctl/pkg/tools/i18n"
	"github.com/spf13/pflag"
	"os"
)

// AddEnvironmentContextFlag ensures consistency between all commands requiring an environment context flag (-c)
func AddEnvironmentContextFlag(flags *pflag.FlagSet, contextFilepath *string) {
	flags.StringVarP(
		contextFilepath,
		i18n.T("cmdFlagContextName"),
		"c",
		getEnv("XCTL_CONTEXT", "-"),
		i18n.T("cmdFlagContextUsage"),
	)
}

func getEnv(name string, defaultValue string) string {
	value := os.Getenv(name)

	if value != "" {
		return value
	}

	return defaultValue
}
