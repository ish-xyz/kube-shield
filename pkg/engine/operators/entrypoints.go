package operators

import v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"

// entrypoint function for Equal operator
func equal(payload string, check *v1.Check) *v1.CheckResult {
	return compareValues(payload, check)
}

// entrypoint function for NotEqual operator
func notEqual(payload string, check *v1.Check) *v1.CheckResult {
	res := compareValues(payload, check)
	res.Match = !res.Match

	return res
}

func greaterThan(payload string, check *v1.Check) *v1.CheckResult {
	return compareNumbers(payload, check, false)
}

func lowerThan(payload string, check *v1.Check) *v1.CheckResult {

	res := compareNumbers(payload, check, false)
	res.Match = !res.Match

	return res
}

func greaterOrEqual(payload string, check *v1.Check) *v1.CheckResult {
	return compareNumbers(payload, check, true)
}

func lowerOrEqual(payload string, check *v1.Check) *v1.CheckResult {

	res := compareNumbers(payload, check, true)
	res.Match = !res.Match

	return res
}

// // entrypoint function for IterateEqual operator
// func equalIterate(payload string, check *v1.Check) *v1.CheckResult {

// 	// TODO: not implemented yet
// 	// get values from Array()
// 	// iterate through each item and compare with static value

// 	return nil
// }
