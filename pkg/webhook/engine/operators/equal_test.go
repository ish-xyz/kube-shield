package operators

import (
	"testing"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/stretchr/testify/assert"
)

func TestEqual(t *testing.T) {

	payload := `{ "my": {"data":"mydata"}}`
	check := &v1.Check{
		Field:    "$_.my.data",
		Operator: "Equal",
		Value:    "mydata",
	}

	res := equal(payload, check)

	assert.Equal(t, res.Error, "")
	assert.Equal(t, res.Result, true)
}

func TestEqualFailed(t *testing.T) {

	payload := `{ "my": {"data":"mydata"}}`
	check := &v1.Check{
		Field:    "$_.my.data",
		Operator: "Equal",
		Value:    "FAILED",
	}

	res := equal(payload, check)

	assert.NotEqual(t, res.Error, "")
	assert.Equal(t, res.Result, false)
}

func TestEqualTypeMismatchBool(t *testing.T) {

	payload := `{ "my": {"data": false }}`
	check := &v1.Check{
		Field:    "$_.my.data",
		Operator: "Equal",
		Value:    "false",
	}

	res := equal(payload, check)

	assert.Equal(t, res.Error, "")
	assert.Equal(t, res.Result, true)
}

func TestEqualTypeMismatchInt(t *testing.T) {

	payload := `{ "my": {"data": 1 }}`
	check := &v1.Check{
		Field:    "$_.my.data",
		Operator: "Equal",
		Value:    "1",
	}

	res := equal(payload, check)

	assert.Equal(t, res.Error, "")
	assert.Equal(t, res.Result, true)
}

func TestEqualTypeMismatchFloat(t *testing.T) {

	payload := `{ "my": {"data": 1.25 }}`
	check := &v1.Check{
		Field:    "$_.my.data",
		Operator: "Equal",
		Value:    "1.25",
	}

	res := equal(payload, check)

	assert.Equal(t, res.Error, "")
	assert.Equal(t, res.Result, true)
}
