package reconciliation

import (
	"context"
	"errors"
	"io"

	"github.com/deifyed/xctl/pkg/config"
	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/deifyed/xctl/pkg/tools/secrets"

	"github.com/spf13/afero"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

type reconciliationResult struct {
	ID     string
	Status string
}

const (
	statusReady    = "ready"
	statusNotReady = "notready"
	statusError    = "error"
)

type ledgerEntry struct {
	Status string
	Count  int
}

// Run initiates scheduling of reconcilers
func (c *Scheduler) Run(ctx context.Context) (Result, error) {
	log := logging.GetLogger(logFeature, "run")
	ledger := generateLedger(c.reconcilers)
	reconciliationContext := c.metadata(ctx)

	channel := make(chan reconciliationResult, 1)
	defer func() {
		close(channel)
	}()

	for hasThingsToDo(ledger) {
		log.Debug("Iteration")
		queue := generateQueue(c.reconcilers, ledger)
		maxUpdates := len(queue.reconcilers)
		updates := 0

		log.Debug("Starting reconcilers")
		for reconciler := queue.Pop(); reconciler != nil; reconciler = queue.Pop() {
			ledger[reconciler.String()].Count++

			go reconcile(channel, reconciliationContext, reconciler)
		}

		for updates < maxUpdates {
			log.Debug("Awaiting result")
			status := <-channel
			log.Debugf("%s: %s", status.ID, status.Status)

			if status.Status == statusReady {
				ledger[status.ID].Status = statusReady
			}

			updates++
		}

		if hasError(ledger) {
			return Result{}, errors.New("reconciling")
		}

		if hasTooManyRequeues(ledger) {
			return Result{}, ErrMaximumReconciliationRequeues
		}
	}

	return Result{}, nil
}

func hasTooManyRequeues(ledger map[string]*ledgerEntry) bool {
	for _, entry := range ledger {
		if entry.Count > config.DefaultMaxReconciliationRequeues {
			return true
		}
	}

	return false
}

func (c *Scheduler) metadata(ctx context.Context) Context {
	return Context{
		Ctx:                    ctx,
		Filesystem:             c.fs,
		Out:                    c.out,
		Keyring:                c.keyring,
		RootDirectory:          c.rootDirectory,
		EnvironmentManifest:    c.environmentManifest,
		ApplicationDeclaration: c.applicationManifest,
		Purge:                  c.purgeFlag,
	}
}

// NewScheduler initializes a Scheduler
func NewScheduler(opts SchedulerOpts, reconcilers ...Reconciler) Scheduler {
	return Scheduler{
		fs:            opts.Filesystem,
		out:           opts.Out,
		keyring:       opts.Keyring,
		rootDirectory: opts.RootDirectory,

		purgeFlag:           opts.PurgeFlag,
		environmentManifest: opts.EnvironmentManifest,
		applicationManifest: opts.ApplicationManifest,

		reconcilers: reconcilers,
	}
}

// SchedulerOpts contains required data for scheduling reconciliations
type SchedulerOpts struct {
	Filesystem *afero.Afero
	// Out provides reconcilers a way to express data
	Out io.Writer
	// Keyring provides reconcilers access to the keyring
	Keyring secrets.Client
	// RootDirectory defines the working directory, usually the IAC repository root
	RootDirectory string

	// Context of the scheduling. Signifies the intent of the user
	// PurgeFlag indicates if everything should be deleted
	PurgeFlag           bool
	EnvironmentManifest v1alpha1.Environment
	ApplicationManifest v1alpha1.Application
}

// Scheduler knows how to run reconcilers in a reasonable way
type Scheduler struct {
	fs            *afero.Afero
	out           io.Writer
	keyring       secrets.Client
	rootDirectory string

	purgeFlag           bool
	environmentManifest v1alpha1.Environment
	applicationManifest v1alpha1.Application

	reconcilers []Reconciler
}
