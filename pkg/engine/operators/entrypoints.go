package operators

import v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"

// entrypoint function for Equal operator
func equal(payload string, check *v1.Check) *v1.CheckResult {
	equal, msg, err := compareValues(payload, check)
	return &v1.CheckResult{
		Result:  equal,
		Message: msg,
		Error:   err,
	}
}

// entrypoint function for NotEqual operator
func notEqual(payload string, check *v1.Check) *v1.CheckResult {
	equal, msg, err := compareValues(payload, check)

	return &v1.CheckResult{
		Result:  !equal,
		Message: msg,
		Error:   err,
	}
}

// func lower(payload string, check *v1.Check) *v1.CheckResult {
// 	return compareNumbers(payload, check, false)
// }

// func greater(payload string, check *v1.Check) *v1.CheckResult {

// 	res := compareNumbers(payload, check, false)
// 	res.Match = !res.Match

// 	return res
// }

// func lowerOrEqual(payload string, check *v1.Check) *v1.CheckResult {
// 	return compareNumbers(payload, check, true)
// }

// func greaterOrEqual(payload string, check *v1.Check) *v1.CheckResult {

// 	res := compareNumbers(payload, check, true)
// 	res.Match = !res.Match

// 	return res
// }

// // entrypoint function for IterateEqual operator
// func equalIterate(payload string, check *v1.Check) *v1.CheckResult {

// 	// TODO: not implemented yet
// 	// get values from Array()
// 	// iterate through each item and compare with static value

// 	return nil
// }
