package cache

import (
	"fmt"

	"github.com/RedLabsPlatform/kube-shield/pkg/webhook/engine"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func (c *CacheController) onPolicyAdd(obj interface{}) {

	var policy *engine.Policy
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(*unstructured.Unstructured).Object, &policy)
	if err != nil {
		logrus.Errorln("failed to unmarshal unstructured object into Policy")
	}

	// hard coded for testing
	c.CacheIndex.AddPolicyFromAddress("default//v1/deployment/rule1")
	c.CacheIndex.AddPolicyFromAddress("default//v1/deployment/rule2")
	c.CacheIndex.AddPolicyFromAddress("default/apps/v1/deployment/rule2")

	fmt.Println(c.CacheIndex)

	fmt.Println(policy)
}

func (c *CacheController) onPolicyDelete(obj interface{}) {
	logrus.Warnln("onPolicyDelete")
}

func (c *CacheController) onClusterPolicyAdd(obj interface{}) {
	logrus.Warnln("onClusterPolicyAdd")
}

func (c *CacheController) onClusterPolicyDelete(obj interface{}) {
	logrus.Warnln("onClusterPolicyDelete")
}
