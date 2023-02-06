package operators

import (
	"testing"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/stretchr/testify/assert"
)

func TestGreaterWithEqualValue(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: 1}
	res := greaterThan(`{"example": 1}`, check)

	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
}
func TestGreaterThanOK(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: 1}
	res := greaterThan(`{"example": 0}`, check)

	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
}
func TestGreaterFloatInt(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: 1.2}
	res := greaterThan(`{"example": 1}`, check)

	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
}
func TestGreaterFloat(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: 1.2}
	res := greaterThan(`{"example": 1.1}`, check)

	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestGreaterInvalidType(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: "1111"}
	res := greaterThan(`{"example": "wrong-type"}`, check)

	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
}
