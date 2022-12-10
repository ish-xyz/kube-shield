package operators

import (
	"fmt"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
)

const (
	GREATER  = "GreaterThan"
	LOWER    = "LowerThan"
	EQUAL    = "Equal"
	NOTEQUAL = "NOTEQUAL"
)

func compareStrings(rawPayload string, check *v1.Check) *v1.CheckResult {

	var err string
	values := getValues(check.Field, rawPayload)

	if len(values) < 1 && check.Value != "" {
		err = fmt.Sprintf("field: %s returned an empty value, policy has value: %s", check.Field, check.Value)
		return CreateCheckResult(false, err)
	}

	for _, v := range values {
		val := getStringValue(v)
		if val != check.Value {
			err := fmt.Sprintf("value '%s' is not equal to policy defined value: '%s'", val, check.Value)
			return CreateCheckResult(false, err)
		}
	}

	return CreateCheckResult(true, err)
}

func compareNumbers(rawPayload string, check *v1.Check) *v1.CheckResult {

	valN, err := getNumber(check.Value)
	if err != nil {
		return CreateCheckResult(
			true,
			fmt.Sprintf("failed to convert number '%s': '%v'", check.Value, err),
		)
	}

	values := getValues(check.Field, rawPayload)
	for _, v := range values {
		payloadN, err := getNumber(getStringValue(v))
		if err != nil {
			return CreateCheckResult(
				true,
				fmt.Sprintf("failed to convert number '%s': '%v'", getStringValue(v), err),
			)
		}
		if check.Operator == GREATER {
			if payloadN <= valN {
				return CreateCheckResult(
					false,
					"retrieved value '%d' is lower than policy value '%d'",
				)
			}
		}

		if check.Operator == LOWER {
			if payloadN >= valN {
				return CreateCheckResult(
					false,
					"retrieved value '%d' is lower than policy value '%d'",
				)
			}
		}
	}

	return CreateCheckResult(true, "")
}
