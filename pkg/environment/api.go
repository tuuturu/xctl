package environment

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/deifyed/xctl/pkg/plugins/promtail"

	"github.com/deifyed/xctl/pkg/plugins/loki"

	_ "embed"

	"github.com/deifyed/xctl/pkg/tools/reconciliation"

	"github.com/deifyed/xctl/pkg/plugins/grafana"

	"github.com/deifyed/xctl/pkg/plugins/prometheus"

	"github.com/deifyed/xctl/pkg/plugins/certbot"
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
	if logging.GetLevel() == logging.LevelInfo {
		spinnerOut = opts.Out
	} else {
		spinnerOut = io.Discard
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
		prometheus.NewReconciler(opts.Provider),
		grafana.NewReconciler(opts.Provider),
		loki.NewReconciler(opts.Provider),
		promtail.NewReconciler(opts.Provider),
	)

	spin.Start()

	_, err = scheduler.Run(context.Background())
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

//go:embed environment-template.yaml
var clusterTemplate string //nolint:gochecknoglobals

func Scaffold() io.Reader {
	return strings.NewReader(clusterTemplate)
}
