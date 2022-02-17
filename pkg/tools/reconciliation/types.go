package reconciliation

import (
	"context"
	"io"
	"time"

	"github.com/spf13/afero"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"

	"github.com/pkg/errors"
)

const logFeature = "reconciliationScheduler"

// Result contains information about the result of a Reconcile() call
type Result struct {
	// Requeue indicates if this Reconciliation must be run again
	Requeue bool
	// RequeueAfter sets the amount of delay before the requeued reconciliation should be done
	RequeueAfter time.Duration
}

// Action represents actions a Reconciler can take
type Action string

const (
	// ActionCreate indicates creation
	ActionCreate = "create"
	// ActionDelete indicates deletion
	ActionDelete = "delete"
	// ActionNoop indicates no necessary action
	ActionNoop = "noop"
	// ActionWait indicates the need to wait
	ActionWait = "wait"
)

// Context represents metadata required by most if not all operations on services
type Context struct {
	Ctx        context.Context
	Filesystem *afero.Afero
	Out        io.Writer

	EnvironmentManifest    v1alpha1.Environment
	ApplicationDeclaration v1alpha1.Application

	Purge bool
}

// Reconciler defines functions needed for the controller to use a reconciler
type Reconciler interface {
	// Reconcile knows how to do what is necessary to ensure the desired state is achieved
	Reconcile(ctx Context) (Result, error)
	// String returns a name that describes the Reconciler
	String() string
}

var (
	// ErrMaximumReconciliationRequeues represents the reconciler trying a single reconciler too many times
	ErrMaximumReconciliationRequeues = errors.New("max reconciliation requeues reached")
	// ErrIndecisive represents the situation where the reconciler can't figure out what to do
	ErrIndecisive = errors.New("indecisive")
)
