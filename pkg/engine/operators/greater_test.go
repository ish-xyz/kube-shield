package operators

// import (
// 	"testing"

// 	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
// 	"github.com/stretchr/testify/assert"
// )

// func TestGreaterWithEqualValue(t *testing.T) {
// 	payload := `{"example": 1}`
// 	check := &v1.Check{Field: "$_.example", Value: 1}
// 	res := greater(payload, check)

// 	assert.Equal(t, false, res.Match)
// 	assert.Equal(t, CHECK_DONE, res.Status)
// 	assert.NotEmpty(t, res.Error)
// }

// func TestGreaterInt(t *testing.T) {
// 	payload := `{"example": 1}`
// 	check := &v1.Check{Field: "$_.example", Value: 0}
// 	res := greater(payload, check)

// 	assert.Equal(t, true, res.Match)
// 	assert.Equal(t, CHECK_DONE, res.Status)
// 	assert.NotEmpty(t, res.Error)
// }
// func TestGreaterFloatVsInt(t *testing.T) {
// 	payload := `{"example": 1.2}`
// 	check := &v1.Check{Field: "$_.example", Value: 1}
// 	res := greater(payload, check)

// 	assert.Equal(t, true, res.Match)
// 	assert.Equal(t, CHECK_DONE, res.Status)
// 	assert.NotEmpty(t, res.Error)
// }
// func TestGreaterFloat(t *testing.T) {
// 	payload := `{"example": 1.3}`
// 	check := &v1.Check{Field: "$_.example", Value: 1.2}
// 	res := greater(payload, check)

// 	assert.Equal(t, true, res.Match)
// 	assert.Equal(t, CHECK_DONE, res.Status)
// 	assert.NotEmpty(t, res.Error)
// }

// func TestGreaterInvalidType(t *testing.T) {
// 	payload := `{"example": "wrong-type"}`
// 	check := &v1.Check{Field: "$_.example", Value: "1111"}
// 	res := greater(payload, check)

// 	assert.Equal(t, false, res.Match)
// 	assert.Equal(t, CHECK_DONE, res.Status)
// 	assert.NotEmpty(t, res.Error)
// }
