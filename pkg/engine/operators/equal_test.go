package operators

import (
	"testing"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/stretchr/testify/assert"
)

func TestEqualTypeMismatch(t *testing.T) {
	res := equal(`{"example": [1,2,3]}`, &v1.Check{Field: "$_.example", Value: 1})
	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_EXECUTED, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestEqualFail(t *testing.T) {
	res := equal(`{"example": 2}`, &v1.Check{Field: "$_.example", Value: 1})
	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_EXECUTED, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestEqualOK(t *testing.T) {
	res := equal(`{"example": 1}`, &v1.Check{Field: "$_.example", Value: 1})
	assert.Equal(t, true, res.Match)
	assert.Equal(t, CHECK_EXECUTED, res.Status)
	assert.NotEmpty(t, res.Error)
}

func TestEqualInitError(t *testing.T) {
	res := equal(`{"example": 1}`, &v1.Check{Field: "1", Value: 1})
	assert.Equal(t, false, res.Match)
	assert.Equal(t, CHECK_INIT_ERROR, res.Status)
	assert.NotEmpty(t, res.Error)
}
