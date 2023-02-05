package operators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqualTypeMismatch(t *testing.T) {
	res := equal(`{"example": [1,2,3]}`, "$_.example", 1)
	assert.Equal(t, false, res.Match)
	assert.NotEqual(t, "", res.Message)
}

func TestEqualFail(t *testing.T) {
	res := equal(`{"example": 2}`, "$_.example", 1)
	assert.Equal(t, false, res.Match)
	assert.NotEqual(t, "", res.Message)
}

func TestEqualOK(t *testing.T) {
	res := equal(`{"example": 1}`, "$_.example", 1)
	assert.Equal(t, true, res.Match)
	assert.Equal(t, "", res.Message)
}
