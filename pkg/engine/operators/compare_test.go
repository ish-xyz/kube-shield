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
	assert.True(t, res.Match)
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
	assert.True(t, res.Match)
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
	assert.Equal(t, res.Match, false)
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
	assert.False(t, res.Match)
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
	assert.False(t, res.Match)
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
	assert.False(t, res.Match)
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
	assert.False(t, res.Match)
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
	assert.False(t, res.Match)
}
