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

	assert.Equal(t, res1.Result, false)
	assert.Equal(t, res2.Result, false)
	assert.Equal(t, res3.Result, true)
}

// test invalid numbers e.g.: "13fa"
// test very big numbers e.g.: 10000000000000000
// test negative numbers
// test zero
