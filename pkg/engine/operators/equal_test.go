package operators

import (
	"testing"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/stretchr/testify/assert"
)

func TestEqualTypeMismatchIntArray(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: 1}
	checkRes := equal(`{"example": [1,2,3]}`, check)

	assert.Equal(t, false, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "type mismatch", checkRes.Message)
}

func TestEqualTypeMismatchArrayMap(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: []int{1, 2, 3}}
	checkRes := equal(`{"example": {"key":"val"}}`, check)

	assert.Equal(t, false, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "type mismatch", checkRes.Message)
}

func TestEqualTypeMismatchMapArray(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: map[string]interface{}{"key": "val"}}
	checkRes := equal(`{"example": [1,2,3]}`, check)

	assert.Equal(t, false, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "type mismatch", checkRes.Message)
}

func TestEqualArrayTypeMismatch(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: []int{1, 2, 3}}
	checkRes := equal(`{"example": ["1","2","3"]}`, check)

	assert.Equal(t, false, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "different values", checkRes.Message)
}

func TestEqualFail(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: 1}
	checkRes := equal(`{"example": 2}`, check)

	assert.Equal(t, false, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "different values", checkRes.Message)
}

func TestEqualNumber(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: 1}
	checkRes := equal(`{"example": 1}`, check)

	assert.Equal(t, true, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "equal values", checkRes.Message)
}

func TestEqualBool(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: true}
	checkRes := equal(`{"example": true}`, check)

	assert.Equal(t, true, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "equal values", checkRes.Message)
}

func TestEqualString(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: "string"}
	checkRes := equal(`{"example": "string"}`, check)

	assert.Equal(t, true, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "equal values", checkRes.Message)
}

func TestEqualRaw(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: nil}
	checkRes := equal(`{"example": null}`, check)

	assert.Equal(t, true, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "equal values", checkRes.Message)
}

func TestEqualArray(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: []int{1, 2, 3}}
	checkRes := equal(`{"example": [1,2,3]}`, check)

	assert.Equal(t, true, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "equal values", checkRes.Message)
}

func TestEqualCountArray(t *testing.T) {
	check := &v1.Check{Field: "$_.example.#", Value: 3}
	checkRes := equal(`{"example": [1,2,3]}`, check)

	assert.Equal(t, true, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "equal values", checkRes.Message)
}

func TestEqualMapKeysArray(t *testing.T) {
	check := &v1.Check{Field: `$_.example.#.key`, Value: []string{"val1", "val2", "val3"}}
	checkRes := equal(`{"example": [{"key":"val1"}, {"key":"val2"}, {"key":"val3"}]}`, check)

	assert.Equal(t, true, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "equal values", checkRes.Message)
}

func TestEqualOKNestedMap(t *testing.T) {
	payload := `{"example": {"key": {"subkey": "subval"}, "key2": {"subkey2":"subval2"}}}`
	expected := map[string]interface{}{"key": map[string]interface{}{"subkey": "subval"}, "key2": map[string]interface{}{"subkey2": "subval2"}}
	check := &v1.Check{Field: "$_.example", Value: expected}
	checkRes := equal(payload, check)

	assert.Equal(t, true, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "equal values", checkRes.Message)
}

func TestEqualOKMap(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: map[string]string{"key": "val", "nokey": "noval"}}
	checkRes := equal(`{"example": {"key":"val", "nokey":"noval"}}`, check)

	assert.Equal(t, true, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "equal values", checkRes.Message)
}

func TestEqualArrayFail(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: []int{1, 4, 3}}
	checkRes := equal(`{"example": [1,2,3]}`, check)

	assert.Equal(t, false, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "different values", checkRes.Message)
}

func TestEqualArrayDifferentLen(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: []int{1, 2, 3}}
	checkRes := equal(`{"example": [1,2]}`, check)

	assert.Equal(t, false, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "different values", checkRes.Message)
}

func TestEqualMapWrongKey(t *testing.T) {
	payload := `{"example": {"wrong-key": "my-custom-value"}}`
	check := &v1.Check{Field: "$_.example", Value: map[string]interface{}{"key": "my-custom-value"}}
	checkRes := equal(payload, check)

	assert.Equal(t, false, checkRes.Result)
	assert.Nil(t, checkRes.Error)
	assert.Equal(t, "different values", checkRes.Message)
}

func TestEqualInitError(t *testing.T) {
	check := &v1.Check{Field: "KEY_PATH_WITHOUT_PREFIX_$_.", Value: 1}
	checkRes := equal(`{"example": 1}`, check)

	assert.Equal(t, false, checkRes.Result)
	assert.NotNil(t, checkRes.Error)
	assert.Empty(t, checkRes.Message)
}
