package operators

import (
	"fmt"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
)

// Running equal operator
func equal(payload, key string, value interface{}) *v1.CheckResult {

	expectedValue, err := getPayloadValues(fmt.Sprintf(`{"input": %v}`, value), "$_.input")
	if err != nil {
		return &v1.CheckResult{
			Match:   false,
			Message: "failed to retrieved expected value",
		}
	}

	retrievedValue, _ := getPayloadValues(payload, key)
	if err != nil {
		return &v1.CheckResult{
			Match:   false,
			Message: "failed to retrieved expected value",
		}
	}

	if expectedValue.Type != retrievedValue.Type {
		return &v1.CheckResult{
			Match: false,
			Message: fmt.Sprintf(
				"type mismatch retrieved value has type '%s', expected value has type '%s'",
				expectedValue.Type,
				retrievedValue.Type,
			),
		}
	}

	if expectedValue.Raw != retrievedValue.Raw {
		return &v1.CheckResult{
			Match: false,
			Message: fmt.Sprintf(
				"different values expected: '%s', actual: '%s'",
				expectedValue.Raw,
				retrievedValue.Raw,
			),
		}
	}

	return &v1.CheckResult{
		Match:   true,
		Message: "",
	}
}
