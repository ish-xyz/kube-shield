package operators

import (
	"fmt"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
)

func equal(rawPayload string, check *v1.Check) *v1.CheckResult {

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

	// convert Strings

	val, err := getNumber[Number](check.Value)
	if err != nil {
		return CreateCheckResult(false, fmt.Sprintf("can't convert value '%s' to number", check.Value))
	}

	fmt.Println(val)

	return CreateCheckResult(true, "")
}
