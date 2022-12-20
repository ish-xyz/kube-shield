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

func Run(payload string, check *v1.Check) (*v1.CheckResult, error) {
	switch check.Operator {
	case EQUAL:
		return compare(payload, check), nil
	case NOTEQUAL:
		return compare(payload, check), nil
	}
	return nil, fmt.Errorf("unknown operator '%s'", check.Operator)
}

func compare(payload string, check *v1.Check) *v1.CheckResult {

	payloadValues := getPayloadValues(check.Field, payload)
	checkValues := getPolicyValue(check.Value, payload)

	msg := fmt.Sprintf("operator is '%s', fetched value is '%v', defined value is '%s', ", check.Operator, checkValues, payloadValues)

	// if there are no retrieved values and then the check is '== $any' or '!= nil', then fail
	if (len(payloadValues) == 0) && (checkValues != "" && check.Operator == EQUAL) || (checkValues == "" && check.Operator == NOTEQUAL) {
		return CreateCheckResult(false, msg)
	}

	for _, v := range payloadValues {

		val := getTypedPayloadValue(v)
		if check.Operator == EQUAL && val != checkValues {
			return CreateCheckResult(false, msg)
		}

		if check.Operator == NOTEQUAL && val == checkValues {
			return CreateCheckResult(false, msg)
		}

	}

	return CreateCheckResult(true, msg)
}
