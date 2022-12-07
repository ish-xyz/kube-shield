package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	index := NewCacheIndex()

	index.Add(Namespace("default"), Group("_core"), Version("v1"), Kind("Pod"), PolicyName("myPol"))
	assert.NotNil(t, index.Policies["default"])
	assert.NotNil(t, index.Policies["default"]["_core"])
	assert.NotNil(t, index.Policies["default"]["_core"]["v1"])
	assert.NotNil(t, index.Policies["default"]["_core"]["v1"]["Pod"])
	assert.Equal(t, index.Policies["default"]["_core"]["v1"]["Pod"][0], PolicyName("myPol"))
}

func TestDelete(t *testing.T) {
	index := NewCacheIndex()

	index.Add(Namespace("default"), Group("_core"), Version("v1"), Kind("Pod"), PolicyName("myPol"))
	index.Delete(Namespace("default"), Group("_core"), Version("v1"), Kind("Pod"), PolicyName("myPol"))
	assert.NotNil(t, index.Policies["default"])
	assert.NotNil(t, index.Policies["default"]["_core"])
	assert.NotNil(t, index.Policies["default"]["_core"]["v1"])
	assert.Nil(t, index.Policies["default"]["_core"]["v1"]["Pod"])
}
