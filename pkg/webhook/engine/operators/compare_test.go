package operators

import (
	"testing"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/stretchr/testify/assert"
)

func TestEqual2DynamicValues(t *testing.T) {

	payload := `{ "my": {"data":"mydata", "data2":"mydata"}}`
	check := &v1.Check{
		Field:    "$_.my.data",
		Operator: "Equal",
		Value:    "$_.my.data2",
	}

	res := compare(payload, check)

	assert.Equal(t, res.Message, "")
	assert.True(t, res.Result)
}

func TestNotEqual2DynamicValues(t *testing.T) {

	payload := `{ "my": {"data":"mydata", "data2":"mydata2"}}`
	check := &v1.Check{
		Field:    "$_.my.data",
		Operator: "NotEqual",
		Value:    "$_.my.data2",
	}

	res := compare(payload, check)

	assert.Equal(t, res.Message, "")
	assert.True(t, res.Result)
}

func TestEqualFailed(t *testing.T) {

	payload := `{ "my": {"data":"mydata"}}`
	check := &v1.Check{
		Field:    "$_.my.data",
		Operator: "Equal",
		Value:    "FAILED",
	}

	res := compare(payload, check)

	assert.NotEqual(t, res.Message, "")
	assert.Equal(t, res.Result, false)
}

func TestNotEqualFailed(t *testing.T) {

	payload := `{ "my": {"data":"mydata"}}`
	check := &v1.Check{
		Field:    "$_.my.data",
		Operator: "NotEqual",
		Value:    "mydata",
	}

	res := compare(payload, check)

	assert.NotEqual(t, res.Message, "")
	assert.False(t, res.Result)
}

func TestEqualTypeMismatchBool(t *testing.T) {

	payload := `{ "my": {"data": false }}`
	check := &v1.Check{
		Field:    "$_.my.data",
		Operator: "Equal",
		Value:    "false",
	}

	res := compare(payload, check)

	assert.NotEqual(t, res.Message, "")
	assert.False(t, res.Result)
}

func TestEqualTypeMismatchInt(t *testing.T) {

	payload := `{ "my": {"data": "1" }}`
	check := &v1.Check{
		Field:    "$_.my.data",
		Operator: "Equal",
		Value:    1,
	}

	res := compare(payload, check)

	assert.NotEqual(t, res.Message, "")
	assert.False(t, res.Result)
}

func TestEqualEmpty(t *testing.T) {

	payload := `{ "my": {"data": "1" }}`
	check := &v1.Check{
		Field:    "$_.my.dataEMPTY",
		Operator: "Equal",
		Value:    1,
	}

	res := compare(payload, check)

	assert.NotEqual(t, res.Message, "")
	assert.False(t, res.Result)
}

func TestEqualTypeMismatchFloat(t *testing.T) {

	payload := `{ "my": {"data": 1 }}`
	check := &v1.Check{
		Field:    "$_.my.data",
		Operator: "Equal",
		Value:    "1.0",
	}

	res := compare(payload, check)

	assert.NotEqual(t, res.Message, "")
	assert.False(t, res.Result)
}
