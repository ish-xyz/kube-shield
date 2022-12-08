package operators

import (
	"testing"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/stretchr/testify/assert"
)

func TestNotEqualStrings(t *testing.T) {
	jsonData := `{"map": {"with": {"some": "values" }}, "other": ["values", "values"]}`
	check1 := &v1.Check{
		Field:    "$_.map.with.some",
		Operator: "NotEqual",
		Value:    "not-equal",
	}

	check2 := &v1.Check{
		Field:    "$_.other",
		Operator: "NotEqual",
		Value:    "not-equal",
	}

	assert.Equal(t, true, notEqual(jsonData, check1))
	assert.Equal(t, true, notEqual(jsonData, check2))
}

func TestNotEqualTypeMismatch(t *testing.T) {
	jsonData := `{"value": true}`
	check1 := &v1.Check{
		Field:    "$_.value",
		Operator: "NotEqual",
		Value:    true,
	}
	check2 := &v1.Check{
		Field:    "$_.value",
		Operator: "NotEqual",
		Value:    "true",
	}

	assert.Equal(t, false, notEqual(jsonData, check1))
	assert.Equal(t, true, notEqual(jsonData, check2))
}

func TestNotEqualNumbers(t *testing.T) {
	jsonData := `{"value": 1, "float": 1.25}`
	check1 := &v1.Check{
		Field:    "$_.value",
		Operator: "NotEqual",
		Value:    1,
	}
	check2 := &v1.Check{
		Field:    "$_.value",
		Operator: "NotEqual",
		Value:    "1",
	}
	check3 := &v1.Check{
		Field:    "$_.float",
		Operator: "NotEqual",
		Value:    1.25,
	}

	assert.Equal(t, false, notEqual(jsonData, check1))
	assert.Equal(t, true, notEqual(jsonData, check2))
	assert.Equal(t, false, notEqual(jsonData, check3))
}

func TestNotEqualEmptyValue(t *testing.T) {
	jsonData := `{"value": true}`
	check1 := &v1.Check{
		Field:    "$_.EMPTYVAL",
		Operator: "NotEqual",
		Value:    "true",
	}

	assert.Equal(t, true, notEqual(jsonData, check1))
}
