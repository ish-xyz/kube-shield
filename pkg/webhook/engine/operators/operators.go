package operators

import (
	"fmt"

	"regexp"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
)

func equal(rawPayload string, check *v1.Check) *v1.CheckResult {

	var val interface{}
	checkRes := NewCheckResult()
	values := getValues(check.Field, rawPayload)

	// Empty values should be returned immediately
	if len(values) < 1 && check.Value != "" {
		return UpdateCheckResult(checkRes, false, fmt.Errorf("the empty value retrieved is not equal to %v", check.Value))
	}

	// Exit at the first non matching value of the array
	for _, v := range values {
		val = getTypedValue(v)
		if val != check.Value {
			return UpdateCheckResult(checkRes, false, fmt.Errorf("retrieved value '%v' is not equal to '%v'", val, check.Value))
		}
	}

	return UpdateCheckResult(checkRes, true, nil)
}

func notEqual(rawPayload string, check *v1.Check) *v1.CheckResult {

	var val interface{}
	checkRes := NewCheckResult()
	values := getValues(check.Field, rawPayload)

	// Matching empty values should be returned immediately
	if len(values) < 1 && check.Value == "" {
		return UpdateCheckResult(checkRes, false, fmt.Errorf("the empty value retrieved is equal to policy value '%v'", check.Value))
	}

	// Exit at the first matching value of the array
	for _, v := range values {
		val = getTypedValue(v)
		if val == check.Value {
			return UpdateCheckResult(checkRes, false, fmt.Errorf("retrieved value '%v' is equal to policy value '%v'", val, check.Value))
		}
	}

	return UpdateCheckResult(checkRes, true, nil)
}

func regex(rawPayload string, check *v1.Check) *v1.CheckResult {

	checkRes := NewCheckResult()
	values := getValues(check.Field, rawPayload)

	for _, v := range values {

		res, err := regexp.MatchString(fmt.Sprintf("%v", check.Value), fmt.Sprintf("%v", getTypedValue(v)))

		if err == nil && !res {
			err = fmt.Errorf("regex '%v'does not match for value '%v'", check.Value, v)
		}
		if !res {
			return UpdateCheckResult(checkRes, res, err)
		}
	}

	return UpdateCheckResult(checkRes, true, nil)
}

func count(rawPayload string, check *v1.Check) *v1.CheckResult {
	checkRes := NewCheckResult()
	values := getValues(check.Field, rawPayload)

	valuesNStr := fmt.Sprintf("%v", len(values))

	if valuesNStr != fmt.Sprintf("%v", check.Value) {
		return UpdateCheckResult(
			checkRes,
			false,
			fmt.Errorf("counted '%s' elements, policy has set '%v' elements", valuesNStr, check.Value),
		)
	}

	return UpdateCheckResult(checkRes, true, nil)
}

/*
- GreaterThan
- LowerThan
- Len
*/
