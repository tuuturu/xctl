package application

import (
	_ "embed"
	"fmt"
	"io"
	"strings"

	"github.com/deifyed/xctl/pkg/application/plugins/namespace"

	"github.com/deifyed/xctl/pkg/application/plugins/argocd"

	"github.com/deifyed/xctl/pkg/application/manifests"

	"github.com/deifyed/xctl/pkg/environment"

	"github.com/deifyed/xctl/pkg/tools/logging"
	"github.com/deifyed/xctl/pkg/tools/reconciliation"
	"github.com/deifyed/xctl/pkg/tools/spinner"
)

func Reconcile(opts ReconcileOpts) error {
	environmentManifest, err := environment.ExtractManifest(opts.EnvironmentManifest)
	if err != nil {
		return fmt.Errorf("extracting environment manifest: %w", err)
	}

	applicationManifest, err := extractManifest(opts.ApplicationManifest)
	if err != nil {
		return fmt.Errorf("extracting application manifest: %w", err)
	}

	var spinnerOut io.Writer
	if logging.GetLevel() == logging.LevelInfo {
		spinnerOut = opts.Out
	} else {
		spinnerOut = io.Discard
	}

	spin := spinner.NewSpinner(spinnerOut)
	spin.FinalMSG = "✅"
	spin.Suffix = "Reconciling"

	schedulerOpts := reconciliation.SchedulerOpts{
		Filesystem:          opts.Filesystem,
		Out:                 opts.Out,
		PurgeFlag:           opts.Purge,
		RootDirectory:       opts.RepositoryRootDirectory,
		EnvironmentManifest: environmentManifest,
		ApplicationManifest: applicationManifest,
	}

	absoluteApplicationDirectory := applicationsDir(opts.RepositoryRootDirectory, applicationManifest.Metadata.Name)
	absoluteEnvironmentDirectory := environmentDir(opts.RepositoryRootDirectory, environmentManifest.Metadata.Name)

	scheduler := reconciliation.NewScheduler(schedulerOpts,
		manifests.Reconciler(absoluteApplicationDirectory),
		namespace.Reconciler(absoluteEnvironmentDirectory),
		argocd.Reconciler(absoluteEnvironmentDirectory),
	)

	spin.Start()

	_, err = scheduler.Run(opts.Context)
	if err != nil {
		spin.FinalMSG = "❌\n\n"
		spin.Stop()

		return fmt.Errorf("scheduling: %w", err)
	}

	spin.Suffix = "Reconciliation finished"

	spin.Stop()

	fmt.Fprintf(opts.Out, "\n\nReconciliation complete\n")

	return nil
}

//go:embed templates/application.yaml
var applicationTemplate string //nolint:gochecknoglobals

// Scaffold returns a stream containing an application template
func Scaffold() io.Reader {
	return strings.NewReader(applicationTemplate)
}
