package operators

import (
	"testing"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/stretchr/testify/assert"
)

func TestNotEqualTypeMismatchIntArray(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: 1}
	checkRes := notEqual(`{"example": [1,2,3]}`, check)

	assert.Equal(t, true, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "type mismatch", checkRes.Message)
}

func TestNotEqualTypeMismatchArrayMap(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: []int{1, 2, 3}}
	checkRes := notEqual(`{"example": {"key":"val"}}`, check)

	assert.Equal(t, true, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "type mismatch", checkRes.Message)
}

func TestNotEqualTypeMismatchMapArray(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: map[string]interface{}{"key": "val"}}
	checkRes := notEqual(`{"example": [1,2,3]}`, check)

	assert.Equal(t, true, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "type mismatch", checkRes.Message)
}

func TestNotEqualArrayTypeMismatch(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: []int{1, 2, 3}}
	checkRes := notEqual(`{"example": ["1","2","3"]}`, check)

	assert.Equal(t, true, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "different values", checkRes.Message)
}

func TestNotEqualFail(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: 1}
	checkRes := notEqual(`{"example": 2}`, check)

	assert.Equal(t, true, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "different values", checkRes.Message)
}

func TestNotEqualNumber(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: 1}
	checkRes := notEqual(`{"example": 1}`, check)

	assert.Equal(t, false, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "equal values", checkRes.Message)
}

func TestNotEqualBool(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: true}
	checkRes := notEqual(`{"example": true}`, check)

	assert.Equal(t, false, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "equal values", checkRes.Message)
}

func TestNotEqualString(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: "string"}
	checkRes := notEqual(`{"example": "string"}`, check)

	assert.Equal(t, false, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "equal values", checkRes.Message)
}

func TestNotEqualRaw(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: nil}
	checkRes := notEqual(`{"example": null}`, check)

	assert.Equal(t, false, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "equal values", checkRes.Message)
}

func TestNotEqualArray(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: []int{1, 2, 3}}
	checkRes := notEqual(`{"example": [1,2,3]}`, check)

	assert.Equal(t, false, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "equal values", checkRes.Message)
}

func TestNotEqualCountArray(t *testing.T) {
	check := &v1.Check{Field: "$_.example.#", Value: 3}
	checkRes := notEqual(`{"example": [1,2,3]}`, check)

	assert.Equal(t, false, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "equal values", checkRes.Message)
}

func TestNotEqualMapKeysArray(t *testing.T) {
	check := &v1.Check{Field: `$_.example.#.key`, Value: []string{"val1", "val2", "val3"}}
	checkRes := notEqual(`{"example": [{"key":"val1"}, {"key":"val2"}, {"key":"val3"}]}`, check)

	assert.Equal(t, false, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "equal values", checkRes.Message)
}

func TestNotEqualOKNestedMap(t *testing.T) {
	payload := `{"example": {"key": {"subkey": "subval"}, "key2": {"subkey2":"subval2"}}}`
	expected := map[string]interface{}{"key": map[string]interface{}{"subkey": "subval"}, "key2": map[string]interface{}{"subkey2": "subval2"}}
	check := &v1.Check{Field: "$_.example", Value: expected}
	checkRes := notEqual(payload, check)

	assert.Equal(t, false, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "equal values", checkRes.Message)
}

func TestNotEqualOKMap(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: map[string]string{"key": "val", "nokey": "noval"}}
	checkRes := notEqual(`{"example": {"key":"val", "nokey":"noval"}}`, check)

	assert.Equal(t, false, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "equal values", checkRes.Message)
}

func TestNotEqualArrayFail(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: []int{1, 4, 3}}
	checkRes := notEqual(`{"example": [1,2,3]}`, check)

	assert.Equal(t, true, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "different values", checkRes.Message)
}

func TestNotEqualArrayDifferentLen(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: []int{1, 2, 3}}
	checkRes := notEqual(`{"example": [1,2]}`, check)

	assert.Equal(t, true, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "different values", checkRes.Message)
}

func TestNotEqualMapWrongKey(t *testing.T) {
	payload := `{"example": {"wrong-key": "my-custom-value"}}`
	check := &v1.Check{Field: "$_.example", Value: map[string]interface{}{"key": "my-custom-value"}}
	checkRes := notEqual(payload, check)

	assert.Equal(t, true, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "different values", checkRes.Message)
}

func TestNotEqualInitError(t *testing.T) {
	check := &v1.Check{Field: "KEY_PATH_WITHOUT_PREFIX_$_.", Value: 1}
	checkRes := notEqual(`{"example": 1}`, check)

	assert.Equal(t, true, checkRes.Result)
	assert.NotNil(t, checkRes.Error)
	assert.Empty(t, checkRes.Message)
}
