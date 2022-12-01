package cache

import (
	"fmt"

	"github.com/RedLabsPlatform/kube-shield/pkg/webhook/engine"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
)

func (c *CacheController) onAdd(obj interface{}) {

	var policy *engine.Policy

	fmt.Println(obj)

	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(map[string]interface{}), &policy)
	fmt.Println(err)
}

func (c *CacheController) onDelete(obj interface{}) {
	logrus.Warnln("onDelete")
}
