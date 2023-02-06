package operators

import (
	"errors"
	"testing"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/stretchr/testify/assert"
)

func TestNotEqualTypeMismatch(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: 1}
	res := notEqual(`{"example": [1,2,3]}`, check)

	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_DONE, res.Status)
	assert.NotEmpty(t, res.Error)
	assert.Equal(t, res.Error, errors.New("type mismatch"))
}

func TestNotEqualTypeMismatchMapArray(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: []int{1, 2, 3}}
	res := notEqual(`{"example": {"key":"val"}}`, check)

	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_DONE, res.Status)
	assert.NotEmpty(t, res.Error)
	assert.Equal(t, res.Error, errors.New("type mismatch"))
}

func TestNotEqualTypeMismatchArrayMap(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: map[string]interface{}{"key": "val"}}
	res := notEqual(`{"example": [1,2,3]}`, check)

	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_DONE, res.Status)
	assert.NotEmpty(t, res.Error)
	assert.Equal(t, res.Error, errors.New("type mismatch"))
}

func TestNotEqualFail(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: 1}
	res := notEqual(`{"example": 2}`, check)

	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_DONE, res.Status)
	assert.NotEmpty(t, res.Error)
	assert.Equal(t, res.Error, errors.New("different values"))
}

func TestNotEqualKNumber(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: 1}
	res := notEqual(`{"example": 1}`, check)

	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_DONE, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestNotEqualBool(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: true}
	res := notEqual(`{"example": true}`, check)

	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_DONE, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestNotEqualString(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: "true"}
	res := notEqual(`{"example": "true"}`, check)

	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_DONE, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestNotEqualRaw(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: nil}
	res := notEqual(`{"example": null}`, check)

	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_DONE, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestNotEqualArray(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: []int{1, 2, 3}}
	res := notEqual(`{"example": [1,2,3]}`, check)

	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_DONE, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestNotEqualCountArray(t *testing.T) {
	check := &v1.Check{Field: "$_.example.#", Value: 3}
	res := notEqual(`{"example": [1,2,3]}`, check)

	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_DONE, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestNotEqualMapKeysArray(t *testing.T) {
	check := &v1.Check{Field: `$_.example.#.key`, Value: []string{"val1", "val2", "val3"}}
	res := notEqual(`{"example": [{"key":"val1"}, {"key":"val2"}, {"key":"val3"}]}`, check)

	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_DONE, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestNotEqualArrayTypeMismatch(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: []int{1, 2, 3}}
	res := notEqual(`{"example": ["1","2","3"]}`, check)

	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_DONE, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestNotEqualArrayDifferentLen(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: []int{1, 2, 3}}
	res := notEqual(`{"example": [1,2]}`, check)

	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_DONE, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestNotEqualMap(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: map[string]string{"key": "val", "nokey": "noval"}}
	res := notEqual(`{"example": {"key":"val", "nokey":"noval"}}`, check)

	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_DONE, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestNotEqualNestedMap(t *testing.T) {
	payload := `{"example": {"key": {"subkey": "subval"}, "key2": {"subkey2":"subval2"}}}`
	expected := map[string]interface{}{"key": map[string]interface{}{"subkey": "subval"}, "key2": map[string]interface{}{"subkey2": "subval2"}}
	check := &v1.Check{Field: "$_.example", Value: expected}
	res := notEqual(payload, check)

	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_DONE, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestNotEqualMapWrongKey(t *testing.T) {
	payload := `{"example": {"wrong-key": "my-custom-value"}}`
	check := &v1.Check{Field: "$_.example", Value: map[string]interface{}{"key": "my-custom-value"}}
	res := notEqual(payload, check)

	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_DONE, res.Status)
	assert.Equal(t, res.Error, errors.New("different values"))
}

func TestNotEqualInitError(t *testing.T) {
	check := &v1.Check{Field: "KEY_PATH_WITHOUT_PREFIX_$_.", Value: 1}
	res := notEqual(`{"example": 1}`, check)

	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_INIT_ERROR, res.Status)
	assert.NotEmpty(t, res.Error)
}
