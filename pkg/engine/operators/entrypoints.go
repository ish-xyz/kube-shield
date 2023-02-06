package operators

import v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"

// entrypoint function for Equal operator
func equal(payload string, check *v1.Check) *v1.CheckResult {
	res := compareValues(payload, check)
	res.Match = res.Match == RUN_EQUAL // this is explicit on purpose
	return res
}

// entrypoint function for NotEqual operator
func notEqual(payload string, check *v1.Check) *v1.CheckResult {
	res := compareValues(payload, check)
	res.Match = res.Match == RUN_NOT_EQUAL // this is explicit on purpose

	return res
}

func greaterThan(payload string, check *v1.Check) *v1.CheckResult {

	res := compareNumbers(payload, check)
	res.Match = res.Match == RUN_GREATER_THAN // this is explicit on purpose

	return res
}

// entrypoint function for IterateEqual operator
func equalIterate(payload string, check *v1.Check) *v1.CheckResult {

	// TODO: not implemented yet
	// get values from Array()
	// iterate through each item and compare with static value

	return nil
}
