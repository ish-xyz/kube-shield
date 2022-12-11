package operators

import (
	"testing"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/stretchr/testify/assert"
)

func TestCompareNumbers(t *testing.T) {

	check1 := &v1.Check{
		Field:    "$_.value1",
		Operator: "GreaterThan",
		Value:    "1.25",
	}
	check2 := &v1.Check{
		Field:    "$_.value1",
		Operator: "LowerThan",
		Value:    "1.10",
	}

	check3 := &v1.Check{
		Field:    "$_.value1",
		Operator: "GreaterThan",
		Value:    "1.13",
	}

	payload := `{"value1": 1.25, "value2": 1.10, "value3": 1.12}`

	res1 := compareNumbers(payload, check1)
	res2 := compareNumbers(payload, check2)
	res3 := compareNumbers(payload, check3)

	assert.False(t, res1.Result)
	assert.False(t, res2.Result)
	assert.True(t, res3.Result)
}

func TestBigNumbers(t *testing.T) {
	check1 := &v1.Check{
		Field:    "$_.value1",
		Operator: "LowerThan",
		Value:    "10000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	}

	payload := `{"value1": 99, "value2": 1.10, "value3": 1.12}`

	res1 := compareNumbers(payload, check1)

	assert.Equal(t, res1.Message, "")
	assert.True(t, res1.Result)
}

func TestInvalidNumbers(t *testing.T) {
	check1 := &v1.Check{
		Field:    "$_.value1",
		Operator: "LowerThan",
		Value:    "12FA",
	}

	payload := `{"value1": 99, "value2": 1.10, "value3": 1.12}`

	res1 := compareNumbers(payload, check1)

	assert.NotEqual(t, res1.Message, "")
	assert.False(t, res1.Result)
}

func TestNegativeNumbers(t *testing.T) {
	check1 := &v1.Check{
		Field:    "$_.value1",
		Operator: "GreaterThan",
		Value:    "-10000000",
	}

	payload := `{"value1": 99}`

	res1 := compareNumbers(payload, check1)

	assert.Equal(t, res1.Message, "")
	assert.True(t, res1.Result)
}

func TestZeroNum(t *testing.T) {
	check1 := &v1.Check{
		Field:    "$_.value1",
		Operator: "GreaterThan",
		Value:    "0",
	}

	payload := `{"value1": 99}`

	res1 := compareNumbers(payload, check1)

	assert.Equal(t, res1.Message, "")
	assert.True(t, res1.Result)
}

// test zero
