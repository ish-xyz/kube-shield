package operators

import (
	"testing"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/stretchr/testify/assert"
)

func TestRegexExactMatchTypes(t *testing.T) {
	jsonData := `{
		"bool": true,
		"float": 1.25,
		"int": 1,
		"map": {"some": "value"},
		"array": [1,2,3],
		"string": "here"
	}`
	checkBool := &v1.Check{
		Field:    "$_.bool",
		Operator: "Regex",
		Value:    "^true$",
	}

	checkFloat := &v1.Check{
		Field:    "$_.float",
		Operator: "Regex",
		Value:    "^1.25$",
	}

	checkInt := &v1.Check{
		Field:    "$_.int",
		Operator: "Regex",
		Value:    "^[0-9]$",
	}

	checkMap := &v1.Check{
		Field:    "$_.map",
		Operator: "Regex",
		Value:    `^{"some": "value"}$`,
	}

	checkArr := &v1.Check{
		Field:    "$_.array",
		Operator: "Regex",
		Value:    `^[1,2,3]$`,
	}

	checkStr := &v1.Check{
		Field:    "$_.string",
		Operator: "Regex",
		Value:    `^here$`,
	}

	resBool := regex(jsonData, checkBool)
	assert.Equal(t, true, resBool.Result)
	assert.Equal(t, resBool.Errors, []error{error(nil)})

	resInt := regex(jsonData, checkInt)
	assert.Equal(t, true, resInt.Result)
	assert.Equal(t, resInt.Errors, []error{error(nil)})

	resFloat := regex(jsonData, checkFloat)
	assert.Equal(t, true, resFloat.Result)
	assert.Equal(t, resFloat.Errors, []error{error(nil)})

	resMap := regex(jsonData, checkMap)
	assert.Equal(t, true, resMap.Result)
	assert.Equal(t, resMap.Errors, []error{error(nil)})

	resArr := regex(jsonData, checkArr)
	assert.Equal(t, true, resArr.Result)
	assert.Equal(t, resArr.Errors, []error{error(nil)})

	resStr := regex(jsonData, checkStr)
	assert.Equal(t, true, resStr.Result)
	assert.Equal(t, resStr.Errors, []error{error(nil)})

}

func TestRegexDontMatch(t *testing.T) {
	jsonData := `{
		"bool": true,
		"float": 1.25,
		"int": 1,
		"map": {"some": "value"},
		"array": [1,2,3],
		"string": "here"
	}`
	checkStr := &v1.Check{
		Field:    "$_.string",
		Operator: "Regex",
		Value:    `^hereDONTMATCH$`,
	}

	res := regex(jsonData, checkStr)
	assert.Equal(t, false, res.Result)
	assert.NotNil(t, res.Errors)
}
