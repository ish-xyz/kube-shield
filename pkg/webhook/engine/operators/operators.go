package operators

import (
	"fmt"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
)

const (
	GREATER  = "GreaterThan"
	LOWER    = "LowerThan"
	EQUAL    = "Equal"
	NOTEQUAL = "NotEqual"
)

func compareStrings(payload string, check *v1.Check) *v1.CheckResult {

	var msg string
	payloadValues := getPayloadValues(check.Field, payload)

	if len(payloadValues) < 1 {

		msg = fmt.Sprintf("%s: field: %s returned an empty value, policy has value: %s", check.Operator, check.Field, check.Value)

		if check.Value != "" && check.Operator == EQUAL {
			return CreateCheckResult(false, msg)
		}
		if check.Value == "" && check.Operator == NOTEQUAL {
			return CreateCheckResult(false, msg)
		}

	}

	for _, v := range payloadValues {

		val := getTypedValue(v)
		msg := fmt.Sprintf("%s: retrieved value '%s' policy defined value: '%s'", check.Operator, val, check.Value)

		if check.Operator == EQUAL {
			if val != check.Value {
				return CreateCheckResult(false, msg)
			}
		}

		if check.Operator == NOTEQUAL {
			if val == check.Value {
				return CreateCheckResult(false, msg)
			}
		}

	}

	return CreateCheckResult(true, "")
}

func compareNumbers(payload string, check *v1.Check) *v1.CheckResult {

	checkVal, err := getFloat(check.Value)
	if err != nil {
		return CreateCheckResult(
			false,
			fmt.Sprintf("%s: invalid check.Value '%v' is not a number", check.Operator, check.Value),
		)
	}

	values := getPayloadValues(check.Field, payload)
	for _, v := range values {

		payloadVal, err := getFloat(getTypedValue(v))
		if err != nil {
			return CreateCheckResult(
				false,
				fmt.Sprintf("%s: invalid value retrieved '%v' is not a number", check.Operator, v),
			)
		}

		if check.Operator == GREATER {
			if payloadVal <= checkVal {
				return CreateCheckResult(
					false,
					fmt.Sprintf("%s: retrieved value '%v' is lower than policy value '%v'", check.Operator, payloadVal, checkVal),
				)
			}
		}

		if check.Operator == LOWER {
			if payloadVal >= checkVal {
				return CreateCheckResult(
					false,
					fmt.Sprintf("%s: retrieved value '%v' is greater than policy value '%v'", check.Operator, payloadVal, checkVal),
				)
			}
		}
	}

	return CreateCheckResult(true, "")
}
