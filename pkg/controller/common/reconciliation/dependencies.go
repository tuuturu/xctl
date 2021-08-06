package reconciliation

import (
	"fmt"
)

// DependencyTestFn defines a function which tests if a dependency is met
type DependencyTestFn func() (bool, error)

// AssertDependencyExistence asserts that the existence of all the provided tests is as expected
func AssertDependencyExistence(expectExistence bool, tests ...DependencyTestFn) (bool, error) {
	for _, test := range tests {
		actualExistence, err := test()
		if err != nil {
			return true, fmt.Errorf("checking dependency: %w", err)
		}

		if expectExistence != actualExistence {
			return false, nil
		}
	}

	return true, nil
}
