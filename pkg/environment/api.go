package environment

import (
	"context"
	"fmt"
	"io"
	"strings"

	_ "embed"

	"github.com/deifyed/xctl/pkg/tools/reconciliation"

	"github.com/deifyed/xctl/pkg/plugins/grafana"

	"github.com/deifyed/xctl/pkg/plugins/prometheus"

	"github.com/deifyed/xctl/pkg/plugins/certbot"
	"github.com/deifyed/xctl/pkg/plugins/vault"

	ingress "github.com/deifyed/xctl/pkg/plugins/nginx-ingress-controller"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/deifyed/xctl/pkg/tools/spinner"
)

// Reconcile knows how to ensure reality for an environment is as declared in an environment manifest
func Reconcile(opts ReconcileOpts) error {
	log := logging.GetLogger("cmd", "cluster")

	manifest, err := extractEnvironmentManifest(opts.Manifest)
	if err != nil {
		return fmt.Errorf("extracting manifest: %w", err)
	}

	var spinnerOut io.Writer
	if opts.Debug {
		spinnerOut = io.Discard
	} else {
		spinnerOut = opts.Out
	}

	spin := spinner.NewSpinner(spinnerOut)
	spin.FinalMSG = "✅"

	schedulerOpts := reconciliation.SchedulerOpts{
		Filesystem:                      opts.Filesystem,
		Out:                             opts.Out,
		PurgeFlag:                       opts.Purge,
		ClusterDeclaration:              manifest,
		ReconciliationLoopDelayFunction: reconciliation.DefaultDelayFunction,
		QueueStepFunc: func(identifier string) {
			log.Debug(fmt.Sprintf("reconciling %s", identifier))

			spin.Suffix = fmt.Sprintf(" Reconciling %s", identifier)
		},
	}

	scheduler := reconciliation.NewScheduler(schedulerOpts,
		NewClusterReconciler(opts.Provider),
		ingress.NewReconciler(opts.Provider),
		NewDomainReconciler(opts.Provider),
		certbot.NewReconciler(opts.Provider),
		vault.NewReconciler(opts.Provider),
		prometheus.NewReconciler(opts.Provider),
		grafana.NewReconciler(opts.Provider),
	)

	spin.Start()

	_, err = scheduler.Run(context.Background())
	if err != nil {
		spin.Stop()

		return fmt.Errorf("scheduling: %w", err)
	}

	spin.Stop()

	fmt.Fprintf(opts.Out, "\n\nReconciliation complete\n")

	return nil
}

//go:embed environment-template.yaml
var clusterTemplate string //nolint:gochecknoglobals

func Scaffold() io.Reader {
	return strings.NewReader(clusterTemplate)
}
