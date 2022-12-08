package operators

import (
	"strings"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/tidwall/gjson"
)

func getValues(address string, jsonData string) []gjson.Result {
	address = strings.TrimPrefix(address, "$_.")
	return gjson.Get(jsonData, address).Array()
}

func Equal(rawPayload string, check *v1.Check) bool {

	var val interface{}

	values := getValues(check.Field, rawPayload)

	// Empty values should be returned immediately
	if len(values) < 1 && check.Value != "" {
		return false
	}

	// Exit at the first non matching value of the array
	for _, retrievedVal := range values {
		val = retrievedVal.Raw
		if val != check.Value {
			return false
		}
	}

	return true
}

// func NotEqual(rawPayload string, check *v1.Check) bool {

// 	var val interface{}

// 	values := getValues(check.Field, rawPayload)
// 	if len(values) < 1 && check.Value == "" {
// 		return false
// 	}

// 	for _, retrievedVal := range values {
// 		val = retrievedVal.Str
// 		if val == check.Value {
// 			return false
// 		}
// 	}

// 	return true
// }
