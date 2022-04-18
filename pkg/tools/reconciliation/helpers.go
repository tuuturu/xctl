package reconciliation

import (
	"errors"
	"time"

	"github.com/deifyed/xctl/pkg/config"
	"github.com/deifyed/xctl/pkg/tools/clients/helm"
	"github.com/deifyed/xctl/pkg/tools/clients/kubectl"
)

// DetermineUserIndication knows how to interpret what operation the user wants for the certain reconciler
func DetermineUserIndication(metadata Context, componentFlag bool) Action {
	if metadata.Purge || !componentFlag {
		return ActionDelete
	}

	return ActionCreate
}

// NoopWaitIndecisiveHandler handles NOOP, Wait and indecisiveness in a streamlined way
func NoopWaitIndecisiveHandler(action Action) (Result, error) {
	switch action {
	case ActionWait:
		return Result{Requeue: true}, nil
	case ActionNoop:
		return Result{Requeue: false}, nil
	default:
		return Result{}, ErrIndecisive
	}
}

// DefaultDelayFunction defines a sane default reconciliation loop delay function
func DefaultDelayFunction() {
	time.Sleep(config.DefaultReconciliationLoopDelayDuration)
}

func isQueueableError(err error) bool {
	queueableErrors := []error{
		helm.ErrUnreachable, helm.ErrTimeout,
		kubectl.ErrConnectionRefused,
	}

	for _, potentialQueueableError := range queueableErrors {
		if errors.Is(err, potentialQueueableError) {
			return true
		}
	}

	return false
}
