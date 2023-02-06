package operators

// import (
// 	"testing"

// 	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
// 	"github.com/stretchr/testify/assert"
// )

// func TestLowerWithEqualValue(t *testing.T) {
// 	check := &v1.Check{Field: "$_.example", Value: 1}
// 	res := lower(`{"example": 1}`, check)

// 	assert.Equal(t, false, res.Match)
// 	assert.Equal(t, CHECK_DONE, res.Status)
// 	assert.NotEmpty(t, res.Error)
// }

// func TestLowerInt(t *testing.T) {
// 	check := &v1.Check{Field: "$_.example", Value: 1}
// 	res := lower(`{"example": 0}`, check)

// 	assert.Equal(t, true, res.Match)
// 	assert.Equal(t, CHECK_DONE, res.Status)
// 	assert.NotEmpty(t, res.Error)
// }

// func TestLowerNegativeNumber(t *testing.T) {
// 	check := &v1.Check{Field: "$_.example", Value: 1}
// 	res := lower(`{"example": -1}`, check)

// 	assert.Equal(t, true, res.Match)
// 	assert.Equal(t, CHECK_DONE, res.Status)
// 	assert.NotEmpty(t, res.Error)
// }

// func TestLowerFloatVsInt(t *testing.T) {
// 	check := &v1.Check{Field: "$_.example", Value: 1.0}
// 	res := lower(`{"example": 0.2}`, check)

// 	assert.Equal(t, true, res.Match)
// 	assert.Equal(t, CHECK_DONE, res.Status)
// 	assert.NotEmpty(t, res.Error)
// }

// func TestLowerFloat(t *testing.T) {
// 	check := &v1.Check{Field: "$_.example", Value: 1.2}
// 	res := lower(`{"example": 1.1}`, check)

// 	assert.Equal(t, true, res.Match)
// 	assert.Equal(t, CHECK_DONE, res.Status)
// 	assert.NotEmpty(t, res.Error)
// }
