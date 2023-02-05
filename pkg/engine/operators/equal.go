package operators

import (
	"fmt"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
)

// Running Equal operators
func equal(payload string, check *v1.Check) *v1.CheckResult {

	// Initialisation
	res := &v1.CheckResult{Status: CHECK_INIT_ERROR, Match: false, Error: fmt.Errorf("init error")}
	expectedValue, _ := getValue(fmt.Sprintf(`{"input": %v}`, check.Value), "$_.input")
	retrievedValue, err := getValue(payload, check.Field)
	if err != nil {
		res.Error = err
		return res
	}

	if expectedValue.Type != retrievedValue.Type {
		res.Error = fmt.Errorf(
			"type mismatch - expected: '%s', actual: '%s'",
			expectedValue.Type,
			retrievedValue.Type,
		)
		return res
	}

	if expectedValue.Raw != retrievedValue.Raw {
		res.Error = fmt.Errorf(
			"different values - expected: '%s', actual: '%s'",
			expectedValue.Raw,
			retrievedValue.Raw,
		)
		return res
	}

	// expectedValue and retrievedValue are a match
	res.Match = true
	res.Error = fmt.Errorf("equal values - expected: '%s', actual: '%s'", expectedValue.Raw, retrievedValue.Raw)

	return res
}
