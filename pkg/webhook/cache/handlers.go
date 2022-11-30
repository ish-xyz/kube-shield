package cache

import (
	"github.com/sirupsen/logrus"
)

func (c *CacheController) onAdd(obj interface{}) {
	logrus.Warnln("onAdd")
}

func (c *CacheController) onDelete(obj interface{}) {
	logrus.Warnln("onDelete")
}
