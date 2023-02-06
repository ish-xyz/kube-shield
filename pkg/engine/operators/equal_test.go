package operators

import (
	"testing"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/stretchr/testify/assert"
)

func TestEqualTypeMismatch(t *testing.T) {
	res := equal(`{"example": [1,2,3]}`, &v1.Check{Field: "$_.example", Value: 1})
	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestEqualFail(t *testing.T) {
	res := equal(`{"example": 2}`, &v1.Check{Field: "$_.example", Value: 1})
	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_OK, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestEqualOK(t *testing.T) {
	res := equal(`{"example": 1}`, &v1.Check{Field: "$_.example", Value: 1})
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

func TestEqualInitError(t *testing.T) {
	res := equal(`{"example": 1}`, &v1.Check{Field: "1", Value: 1})
	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_INIT_ERROR, res.Status)
	assert.NotEmpty(t, res.Error)
}
