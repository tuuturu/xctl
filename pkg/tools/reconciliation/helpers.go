package reconciliation

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
