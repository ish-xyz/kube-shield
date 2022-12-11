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

	var err string
	payloadValues := getPayloadValues(check.Field, payload)

	if len(payloadValues) < 1 && check.Value != "" {
		err = fmt.Sprintf("field: %s returned an empty value, policy has value: %s", check.Field, check.Value)
		return CreateCheckResult(false, err)
	}

	for _, v := range payloadValues {
		val := getStringValue(v)
		if check.Operator == EQUAL {
			if val != check.Value {
				err := fmt.Sprintf("%s: value '%s' is not equal to policy defined value: '%s'", check.Operator, val, check.Value)
				return CreateCheckResult(false, err)
			}
		}

		if check.Operator == NOTEQUAL {
			if val == check.Value {
				err := fmt.Sprintf("%s: value '%s' is equal to policy defined value: '%s'", check.Operator, val, check.Value)
				return CreateCheckResult(false, err)
			}
		}

	}

	return CreateCheckResult(true, err)
}

func compareNumbers(payload string, check *v1.Check) *v1.CheckResult {

	valN, err := getNumber(check.Value)
	if err != nil {
		return CreateCheckResult(
			false,
			fmt.Sprintf("failed to convert number '%s': '%v'", check.Value, err),
		)
	}

	values := getPayloadValues(check.Field, payload)
	for _, v := range values {
		payloadN, err := getNumber(getStringValue(v))
		if err != nil {
			return CreateCheckResult(
				false,
				fmt.Sprintf("failed to convert number '%s': '%v'", getStringValue(v), err),
			)
		}
		if check.Operator == GREATER {
			if payloadN <= valN {
				return CreateCheckResult(
					false,
					fmt.Sprintf("%s: retrieved value '%v' is lower than policy value '%v'", check.Operator, payloadN, check.Value),
				)
			}
		}

		if check.Operator == LOWER {
			if payloadN >= valN {
				return CreateCheckResult(
					false,
					fmt.Sprintf("%s: retrieved value '%v' is greater than policy value '%v'", check.Operator, payloadN, check.Value),
				)
			}
		}
	}

	return CreateCheckResult(true, "")
}
