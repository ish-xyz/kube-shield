package operators

import (
	"encoding/json"
	"errors"
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

func compareNumbers(payload string, check *v1.Check) *v1.CheckResult {

	// Initialisation
	res := &v1.CheckResult{Status: CHECK_INIT_ERROR, Match: false, Error: fmt.Errorf("init error")}

	checkJson, err := json.Marshal(check)
	if err != nil {
		res.Error = errors.New("failed to initialise check json data")
		return res
	}

	expected, _ := getValue(string(checkJson), "$_.value")
	actual, err := getValue(payload, check.Field)
	if err != nil {
		res.Error = err
		return res
	}

	// Check processesing
	res.Status = CHECK_DONE

	if actual.Type != gjson.Number || expected.Type != gjson.Number {
		res.Error = errors.New("not a number")
		return res
	}

	return nil
}

// function to run both Equal and NotEqual operator
// pass payload as JSON string value and use behaviour to switch between the two operators
// behaviour:
// -> true == Equal
// -> false == NotEqual
func compareValues(payload string, check *v1.Check) *v1.CheckResult {

	// Initialisation
	res := &v1.CheckResult{Status: CHECK_INIT_ERROR, Match: false, Error: fmt.Errorf("init error")}

	checkJson, err := json.Marshal(check)
	if err != nil {
		res.Error = errors.New("failed to initialise check json data")
		return res
	}

	expected, _ := getValue(string(checkJson), "$_.value")
	actual, err := getValue(payload, check.Field)
	if err != nil {
		res.Error = err
		return res
	}

	// Check processesing
	res.Status = CHECK_DONE
	if match := compareTypes(expected, actual); !match {
		res.Error = errors.New("type mismatch")
		return res
	}

	res.Match = compareComplex(expected, actual)
	if !res.Match {
		res.Error = fmt.Errorf("different values")
	}

	return res
}
