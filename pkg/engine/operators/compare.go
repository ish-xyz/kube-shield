package operators

import (
	"encoding/json"
	"fmt"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/tidwall/gjson"
)

// compare types of two gjson.Result objects
func compareTypes(expected, actual gjson.Result) bool {

	if expected.Type != actual.Type {
		return false
	}

	if expected.Type == gjson.Type(gjson.JSON) {
		if expected.IsArray() && !actual.IsArray() {
			return false
		}

		if expected.IsObject() && !actual.IsObject() {
			return false
		}
	}

	return true
}

// compare complex results -> map/key
func compareComplex(expected, actual gjson.Result) bool {

	if expected.Type != gjson.JSON {
		return compareSimple(expected, actual)
	}

	if expected.IsArray() {
		expectedArr := expected.Array()
		actualArr := actual.Array()
		// checking length of the arrays
		if len(expectedArr) != len(actualArr) {
			return false
		}

		for i := range expectedArr {
			// compare types of the values in array
			if !compareTypes(expectedArr[i], actualArr[i]) {
				return false
			}

			// compare values re-iterating this function
			if !compareComplex(expectedArr[i], actualArr[i]) {
				return false
			}
		}
	}

	if expected.IsObject() {
		expectedMap := expected.Map()
		actualMap := actual.Map()

		// checking number of keys in maps is the same
		if len(expectedMap) != len(actualMap) {
			return false
		}
		for k, _ := range expectedMap {

			// check that the key we are looking for exists in the other map
			if _, ok := actualMap[k]; !ok {
				return false
			}

			// compare value types
			if !compareTypes(expectedMap[k], actualMap[k]) {
				return false
			}

			// compare values re-iterating this function
			if !compareComplex(expectedMap[k], actualMap[k]) {
				return false
			}
		}
	}

	return true
}

// Compare basic, non-complex, results
func compareSimple(expected, actual gjson.Result) bool {

	switch atype := expected.Type; atype {

	case gjson.True, gjson.False:
		return expected.Bool() == actual.Bool()
	case gjson.Number:
		return expected.Num == actual.Num
	case gjson.String:
		return expected.String() == actual.String()
	case gjson.Null:
		return expected.Raw == actual.Raw
	}

	return false
}

// function to compare 2 gjson.Result numbers
// returns => isEqual (bool), message (string), executionErrors (error)
func compareNumbers(payload string, check *v1.Check) (bool, string, error) {

	checkJson, err := json.Marshal(check)
	if err != nil {
		return false, "", fmt.Errorf("init error")
	}

	expected, _ := getValue(string(checkJson), "$_.value")
	actual, err := getValue(payload, check.Field)
	if err != nil {
		return false, "", fmt.Errorf("failed to get values from payload")
	}

	// Check processesing
	if actual.Type != gjson.Number || expected.Type != gjson.Number {
		return false, "values passed are not numbers", nil
	}

	return true, "", nil
}

// function to compare 2 gjson.Result
// returns => isEqual (bool), message (string), executionErrors (error)
func compareValues(payload string, check *v1.Check) (bool, string, error) {

	checkJson, err := json.Marshal(check)
	if err != nil {
		return false, "", fmt.Errorf("init error")
	}

	expected, _ := getValue(string(checkJson), "$_.value")
	actual, err := getValue(payload, check.Field)
	if err != nil {
		return false, "", fmt.Errorf("failed to get values from payload")
	}

	// Need to perform types comparison before
	if res := compareTypes(expected, actual); !res {
		return false, "type mismatch", nil
	}

	if !compareComplex(expected, actual) {
		return false, "different values", nil
	}
	return true, "equal values", nil

}
