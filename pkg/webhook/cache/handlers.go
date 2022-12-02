package cache

import (
	"fmt"
	"strings"

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
		return
	}

	for _, res := range policy.Spec.ApplyOn {
		c.CacheIndex.Lock()

		// Core resources don't have a group specified
		// So we set the group to "nil"
		// "_" is there because it's not possible to have groups that starts with _
		gvr := strings.Split(res.APIVersion, "/")
		group := "_nil"
		version := res.APIVersion
		if len(gvr) > 1 {
			group = gvr[0]
			version = gvr[1]
		}

		c.CacheIndex.CachePolicyResourcesMapping(
			Namespace(policy.Namespace),
			Group(group),
			Version(version),
			Kind(res.Kind),
			RuleName(policy.Name),
		)
		c.CacheIndex.Unlock()
	}

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
