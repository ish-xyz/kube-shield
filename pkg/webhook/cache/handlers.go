package cache

import (
	"fmt"
	"strings"

	"github.com/RedLabsPlatform/kube-shield/pkg/webhook/engine"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	DELETE = 0
	ADD    = 1
)

func (c *CacheController) onPolicyAdd(obj interface{}) {

	var policy *engine.Policy
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(*unstructured.Unstructured).Object, &policy)
	if err != nil {
		logrus.Errorln("failed to unmarshal unstructured object into Policy")
		return
	}

	err = c.handleIndexRequest(policy.Spec.ApplyOn, ADD, policy.Namespace, policy.Name)
	if err != nil {
		logrus.Errorln("failed to ADD policy to memory index")
		return
	}
}

func (c *CacheController) onPolicyDelete(obj interface{}) {
	var policy *engine.Policy
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(*unstructured.Unstructured).Object, &policy)
	if err != nil {
		logrus.Errorln("failed to unmarshal unstructured object into Policy")
		return
	}

	err = c.handleIndexRequest(policy.Spec.ApplyOn, DELETE, policy.Namespace, policy.Name)
	if err != nil {
		logrus.Errorln("failed to DELETE policy to memory index")
		return
	}
}

func (c *CacheController) onClusterPolicyAdd(obj interface{}) {
	var clusterpolicy *engine.ClusterPolicy
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(*unstructured.Unstructured).Object, &clusterpolicy)
	if err != nil {
		logrus.Errorln("failed to unmarshal unstructured object into Policy")
		return
	}

	err = c.handleIndexRequest(clusterpolicy.Spec.ApplyOn, ADD, "_ClusterScope", clusterpolicy.Name)
	if err != nil {
		logrus.Errorln("failed to ADD policy to memory index")
		return
	}
}

func (c *CacheController) onClusterPolicyDelete(obj interface{}) {

	var clusterpolicy *engine.ClusterPolicy
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(*unstructured.Unstructured).Object, &clusterpolicy)
	if err != nil {
		logrus.Errorln("failed to unmarshal unstructured object into Policy")
		return
	}

	err = c.handleIndexRequest(clusterpolicy.Spec.ApplyOn, DELETE, "_ClusterScope", clusterpolicy.Name)
	if err != nil {
		logrus.Errorln("failed to DELETE policy to memory index")
		return
	}
}

func (c *CacheController) handleIndexRequest(resources []*engine.ResourceAddress, ops int, namespace string, policyName string) error {

	var fn func(ns Namespace, grp Group, ver Version, kind Kind, name PolicyName)

	if ops == ADD {
		fn = c.CacheIndex.Add
	} else if ops == DELETE {
		fn = c.CacheIndex.Delete
	} else {
		return fmt.Errorf("index handler invalid operations")
	}

	for _, res := range resources {

		// Core resources don't have a group specified
		// So we set the group to "_core"
		// "_" is there because it's not possible to have groups that starts with _ in the k8s CRDs,
		// so it won't overwrite any possible CRD
		gvr := strings.Split(res.APIVersion, "/")
		group := "_core"
		version := res.APIVersion
		if len(gvr) > 1 {
			group = gvr[0]
			version = gvr[1]
		}

		c.CacheIndex.Lock()
		fn(
			Namespace(namespace),
			Group(group),
			Version(version),
			Kind(res.Kind),
			PolicyName(policyName),
		)
		c.CacheIndex.Unlock()
	}
	return nil
}
