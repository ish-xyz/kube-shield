package operators

import (
	"errors"
	"testing"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/stretchr/testify/assert"
)

func TestEqualTypeMismatch(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: 1}
	res := equal(`{"example": [1,2,3]}`, check)

	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
	assert.Equal(t, res.Error, errors.New("type mismatch"))
}
func TestEqualTypeMismatchMapArray(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: []int{1, 2, 3}}
	res := equal(`{"example": {"key":"val"}}`, check)

	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
	assert.Equal(t, res.Error, errors.New("type mismatch"))
}

func TestEqualTypeMismatchArrayMap(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: map[string]interface{}{"key": "val"}}
	res := equal(`{"example": [1,2,3]}`, check)

	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
	assert.Equal(t, res.Error, errors.New("type mismatch"))
}

func TestEqualFail(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: 1}
	res := equal(`{"example": 2}`, check)

	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
	assert.Equal(t, res.Error, errors.New("different values"))
}

func TestEqualNumber(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: 1}
	res := equal(`{"example": 1}`, check)

	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestEqualBool(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: true}
	res := equal(`{"example": true}`, check)

	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestEqualString(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: "true"}
	res := equal(`{"example": "true"}`, check)

	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestEqualRaw(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: nil}
	res := equal(`{"example": null}`, check)

	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestEqualArray(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: []int{1, 2, 3}}
	res := equal(`{"example": [1,2,3]}`, check)

	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestEqualCountArray(t *testing.T) {
	check := &v1.Check{Field: "$_.example.#", Value: 3}
	res := equal(`{"example": [1,2,3]}`, check)

	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestEqualMapKeysArray(t *testing.T) {
	check := &v1.Check{Field: `$_.example.#.key`, Value: []string{"val1", "val2", "val3"}}
	res := equal(`{"example": [{"key":"val1"}, {"key":"val2"}, {"key":"val3"}]}`, check)

	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestEqualArrayFail(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: []int{1, 4, 3}}
	res := equal(`{"example": [1,2,3]}`, check)

	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestEqualArrayTypeMismatch(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: []int{1, 2, 3}}
	res := equal(`{"example": ["1","2","3"]}`, check)

	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestEqualArrayDifferentLen(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: []int{1, 2, 3}}
	res := equal(`{"example": [1,2]}`, check)

	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestEqualOKMap(t *testing.T) {
	check := &v1.Check{Field: "$_.example", Value: map[string]string{"key": "val", "nokey": "noval"}}
	res := equal(`{"example": {"key":"val", "nokey":"noval"}}`, check)

	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestEqualOKNestedMap(t *testing.T) {
	payload := `{"example": {"key": {"subkey": "subval"}, "key2": {"subkey2":"subval2"}}}`
	expected := map[string]interface{}{"key": map[string]interface{}{"subkey": "subval"}, "key2": map[string]interface{}{"subkey2": "subval2"}}
	check := &v1.Check{Field: "$_.example", Value: expected}
	res := equal(payload, check)

	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestEqualMapWrongKey(t *testing.T) {
	payload := `{"example": {"wrong-key": "my-custom-value"}}`
	check := &v1.Check{Field: "$_.example", Value: map[string]interface{}{"key": "my-custom-value"}}
	res := equal(payload, check)

	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.Equal(t, res.Error, errors.New("different values"))
}

func TestEqualInitError(t *testing.T) {
	check := &v1.Check{Field: "KEY_PATH_WITHOUT_PREFIX_$_.", Value: 1}
	res := equal(`{"example": 1}`, check)

	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_INIT_ERROR, res.Status)
	assert.NotEmpty(t, res.Error)
}
