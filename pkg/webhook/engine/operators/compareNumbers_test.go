package operators

import (
	"fmt"
	"testing"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/stretchr/testify/assert"
)

func TestCompareNumbersNoNum(t *testing.T) {
	check := &v1.Check{
		Field:    "$_.value",
		Operator: "GreaterThan",
		Value:    "1.25",
	}

	payload := `{"value": 1.25}`

	res := compareNumbers(payload, check)

	fmt.Println(res)

	assert.Equal(t, 1, 2)
}
