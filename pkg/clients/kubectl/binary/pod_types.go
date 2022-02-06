package binary

type getPodResultStatusContainerStatus struct {
	Ready bool `json:"ready"`
}

type getPodResultStatus struct {
	ContainerStatuses []getPodResultStatusContainerStatus `json:"containerStatuses"`
}

type getPodResult struct {
	Status getPodResultStatus `json:"status"`
}
