package operators

import (
	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
)

func equal(rawPayload string, check *v1.Check) bool {

	var val interface{}

	values := getValues(check.Field, rawPayload)

	// Empty values should be returned immediately
	if len(values) < 1 && check.Value != "" {
		return false
	}

	// Exit at the first non matching value of the array
	for _, v := range values {
		val = getTypedValue(v)
		if val != check.Value {
			return false
		}
	}

	return true
}

func notEqual(rawPayload string, check *v1.Check) bool {

	var val interface{}

	values := getValues(check.Field, rawPayload)

	// Matching empty values should be returned immediately
	if len(values) < 1 && check.Value == "" {
		return false
	}

	// Exit at the first matching value of the array
	for _, v := range values {
		val = getTypedValue(v)
		if val == check.Value {
			return false
		}
	}

	return true
}
