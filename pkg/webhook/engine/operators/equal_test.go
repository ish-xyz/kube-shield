package operators

import (
	"testing"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/stretchr/testify/assert"
)

func TestGetValue(t *testing.T) {
	jsonData := `{"map": {"with": {"some":"values"}}}`
	assert.Equal(t, getValues("$_.map.with.some", jsonData)[0].Str, "values")
}

func TestEqualStrings(t *testing.T) {
	jsonData := `{"map": {"with": {"some": "values" }}, "other": ["values", "values"]}`
	check1 := &v1.Check{
		Field:    "$_.map.with.some",
		Operator: "Equal",
		Value:    "values",
	}

	check2 := &v1.Check{
		Field:    "$_.other",
		Operator: "Equal",
		Value:    "values",
	}

	assert.Equal(t, true, equal(jsonData, check1))
	assert.Equal(t, true, equal(jsonData, check2))
}

func TestEqualTypeMismatch(t *testing.T) {
	jsonData := `{"value": true}`
	check1 := &v1.Check{
		Field:    "$_.value",
		Operator: "Equal",
		Value:    true,
	}
	check2 := &v1.Check{
		Field:    "$_.value",
		Operator: "Equal",
		Value:    "true",
	}

	assert.Equal(t, true, equal(jsonData, check1))
	assert.Equal(t, false, equal(jsonData, check2))
}

func TestEqualNumbers(t *testing.T) {
	jsonData := `{"value": 1, "float": 1.25}`
	check1 := &v1.Check{
		Field:    "$_.value",
		Operator: "Equal",
		Value:    1,
	}
	check2 := &v1.Check{
		Field:    "$_.value",
		Operator: "Equal",
		Value:    "1",
	}
	check3 := &v1.Check{
		Field:    "$_.float",
		Operator: "Equal",
		Value:    1.25,
	}

	assert.Equal(t, true, equal(jsonData, check1))
	assert.Equal(t, false, equal(jsonData, check2))
	assert.Equal(t, true, equal(jsonData, check3))
}

func TestEqualEmptyValue(t *testing.T) {
	jsonData := `{"value": true}`
	check1 := &v1.Check{
		Field:    "$_.EMPTYVAL",
		Operator: "Equal",
		Value:    "true",
	}

	assert.Equal(t, false, equal(jsonData, check1))
}
