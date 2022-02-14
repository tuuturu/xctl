package reconciliation

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/afero"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/tools/logging"
)

// Run initiates scheduling of reconcilers
func (c *Scheduler) Run(ctx context.Context) (Result, error) {
	log := logging.GetLogger(logFeature, "run")
	queue := NewQueue(c.reconcilers)
	reconciliationContext := c.metadata(ctx)

	for reconciler := queue.Pop(); reconciler != nil; reconciler = queue.Pop() {
		c.queueStepFunc(reconciler.String())

		result, err := reconciler.Reconcile(reconciliationContext)
		if err != nil {
			if !isQueueableError(err) {
				return Result{}, fmt.Errorf("reconciling %s: %w", reconciler.String(), err)
			}

			result.Requeue = true
		}

		if result.Requeue {
			log.Debug(fmt.Sprintf("requeueing %s", reconciler.String()))

			err = queue.Push(reconciler)
			if err != nil {
				return Result{}, fmt.Errorf("passing requeue check for %s: %w", reconciler.String(), err)
			}
		}

		c.reconciliationLoopDelayFunction()
	}

	return Result{}, nil
}

func (c *Scheduler) metadata(ctx context.Context) Context {
	return Context{
		Ctx:                    ctx,
		Filesystem:             c.fs,
		Out:                    c.out,
		ClusterDeclaration:     c.clusterDeclaration,
		ApplicationDeclaration: c.applicationDeclaration,
		Purge:                  c.purgeFlag,
	}
}

// NewScheduler initializes a Scheduler
func NewScheduler(opts SchedulerOpts, reconcilers ...Reconciler) Scheduler {
	return Scheduler{
		fs:  opts.Filesystem,
		out: opts.Out,

		purgeFlag:              opts.PurgeFlag,
		clusterDeclaration:     opts.ClusterDeclaration,
		applicationDeclaration: opts.ApplicationDeclaration,

		reconciliationLoopDelayFunction: opts.ReconciliationLoopDelayFunction,
		queueStepFunc:                   opts.QueueStepFunc,
		reconcilers:                     reconcilers,
	}
}

// SchedulerOpts contains required data for scheduling reconciliations
type SchedulerOpts struct {
	Filesystem *afero.Afero
	// Out provides reconcilers a way to express data
	Out io.Writer

	// Context of the scheduling. Signifies the intent of the user
	// PurgeFlag indicates if everything should be deleted
	PurgeFlag bool
	// ReconciliationLoopDelayFunction introduces delay to the reconciliation process
	ReconciliationLoopDelayFunction func()
	ClusterDeclaration              v1alpha1.Cluster
	ApplicationDeclaration          v1alpha1.Application
	QueueStepFunc                   func(identifier string)
}

// Scheduler knows how to run reconcilers in a reasonable way
type Scheduler struct {
	fs  *afero.Afero
	out io.Writer

	purgeFlag              bool
	clusterDeclaration     v1alpha1.Cluster
	applicationDeclaration v1alpha1.Application

	reconciliationLoopDelayFunction func()
	queueStepFunc                   func(string)
	reconcilers                     []Reconciler
}
