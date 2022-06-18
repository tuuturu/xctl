package reconciliation

import (
	"fmt"

	"github.com/deifyed/xctl/pkg/tools/logging"
)

func reconcile(channel chan reconciliationResult, ctx Context, reconciler Reconciler) {
	log := logging.GetLogger(logFeature, fmt.Sprintf("reconcile/%s", reconciler.String()))
	status := reconciliationResult{ID: reconciler.String()}

	log.Debug("Reconciling")
	result, err := reconciler.Reconcile(ctx)
	if err != nil {
		log.Debugf("Got error: %s", err.Error())
		status.Status = statusError
	}

	if result.Requeue {
		status.Status = statusNotReady
	} else {
		status.Status = statusReady
	}

	channel <- status
}

func generateQueue(reconcilers []Reconciler, ledger map[string]*ledgerEntry) Queue {
	relevantReconcilers := make([]Reconciler, 0)

	for _, reconciler := range reconcilers {
		if ledger[reconciler.String()].Status == statusNotReady {
			relevantReconcilers = append(relevantReconcilers, reconciler)
		}
	}

	return NewQueue(relevantReconcilers)
}

func generateLedger(reconcilers []Reconciler) map[string]*ledgerEntry {
	ledger := make(map[string]*ledgerEntry, len(reconcilers))

	for _, reconciler := range reconcilers {
		ledger[reconciler.String()] = &ledgerEntry{Status: statusNotReady}
	}

	return ledger
}

// hasThingsToDo checks if there's any reconcilers who hasn't reported themselves ready, and returns true if it finds
// one
func hasThingsToDo(ledger map[string]*ledgerEntry) bool {
	for _, entry := range ledger {
		if entry.Status == statusNotReady {
			return true
		}
	}

	return false
}

func hasError(ledger map[string]*ledgerEntry) bool {
	for _, entry := range ledger {
		if entry.Status == statusError {
			return true
		}
	}

	return false
}
