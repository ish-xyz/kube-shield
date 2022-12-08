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

func TestEqual(t *testing.T) {
	jsonData := `{"map": {"with": {"some":"values"}}, "other": ["values", "values"]}`
	chk := &v1.Check{
		Field:    "$_.map.with.some",
		Operator: "Equal",
		Value:    "values",
	}

	multiChk := &v1.Check{
		Field:    "$_.other",
		Operator: "Equal",
		Value:    "values",
	}

	assert.Equal(t, Equal(jsonData, chk), true)
	assert.Equal(t, Equal(jsonData, multiChk), true)
}

func TestEqualTypeMismatch(t *testing.T) {
	jsonData := `{"value": true}`
	check1 := &v1.Check{
		Field:    "$_.value",
		Operator: "Equal",
		Value:    "true",
	}

	assert.Equal(t, Equal(jsonData, check1), true)
}

func TestEqualEmptyValue(t *testing.T) {
	jsonData := `{"value": true}`
	check1 := &v1.Check{
		Field:    "$_.EMPTYVAL",
		Operator: "Equal",
		Value:    "true",
	}

	assert.Equal(t, Equal(jsonData, check1), false)
}

// func TestNotEqual(t *testing.T) {
// 	jsonData := `{"some":"values", "other": [{"status": "someval"}, {"status": "myval"}]}`
// 	chk := &v1.Check{
// 		Field:    "$_.some",
// 		Operator: "NotEqual",
// 		Value:    "not-equal-value",
// 	}

// 	multiChk := &v1.Check{
// 		Field:    "$_.other.#.status",
// 		Operator: "NotEqual",
// 		Value:    "myval",
// 	}

// 	assert.Equal(
// 		t,
// 		NotEqual(jsonData, chk),
// 		true,
// 	)
// 	assert.Equal(
// 		t,
// 		NotEqual(jsonData, multiChk),
// 		false,
// 	)
// }
