package reconciliation

import (
	"time"

	"github.com/deifyed/xctl/pkg/config"
)

// DetermineUserIndication knows how to interpret what operation the user wants for the certain reconciler
func DetermineUserIndication(metadata Context, componentFlag bool) Action {
	if metadata.Purge || !componentFlag {
		return ActionDelete
	}

	return ActionCreate
}

// DefaultDelayFunction defines a sane default reconciliation loop delay function
func DefaultDelayFunction() {
	time.Sleep(config.DefaultReconciliationLoopDelayDuration)
}
