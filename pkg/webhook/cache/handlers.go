package cache

import (
	"fmt"
	"strings"

	"github.com/RedLabsPlatform/kube-shield/pkg/webhook/engine"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

//TODO: too much code repetition here, clean up.

func (c *CacheController) onPolicyAdd(obj interface{}) {

	var policy *engine.Policy
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(*unstructured.Unstructured).Object, &policy)
	if err != nil {
		logrus.Errorln("failed to unmarshal unstructured object into Policy")
		return
	}

	for _, res := range policy.Spec.ApplyOn {

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
		c.CacheIndex.Add(
			Namespace(policy.Namespace),
			Group(group),
			Version(version),
			Kind(res.Kind),
			PolicyName(policy.Name),
		)
		c.CacheIndex.Unlock()
	}
	fmt.Println(c.CacheIndex.Policies)
}

func (c *CacheController) onPolicyDelete(obj interface{}) {
	var policy *engine.Policy
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(*unstructured.Unstructured).Object, &policy)
	if err != nil {
		logrus.Errorln("failed to unmarshal unstructured object into Policy")
		return
	}

	for _, res := range policy.Spec.ApplyOn {

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
		c.CacheIndex.Remove(
			Namespace(policy.Namespace),
			Group(group),
			Version(version),
			Kind(res.Kind),
			PolicyName(policy.Name),
		)
		c.CacheIndex.Unlock()
	}
	fmt.Println(c.CacheIndex.Policies)
}

func (c *CacheController) onClusterPolicyAdd(obj interface{}) {
	var clusterpolicy *engine.ClusterPolicy
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(*unstructured.Unstructured).Object, &clusterpolicy)
	if err != nil {
		logrus.Errorln("failed to unmarshal unstructured object into Policy")
		return
	}

	for _, res := range clusterpolicy.Spec.ApplyOn {

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
		c.CacheIndex.Add(
			Namespace("_ClusterScope"),
			Group(group),
			Version(version),
			Kind(res.Kind),
			PolicyName(clusterpolicy.Name),
		)
		c.CacheIndex.Unlock()
	}
	fmt.Println(c.CacheIndex.Policies)
}

func (c *CacheController) onClusterPolicyDelete(obj interface{}) {

	var clusterpolicy *engine.ClusterPolicy
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(*unstructured.Unstructured).Object, &clusterpolicy)
	if err != nil {
		logrus.Errorln("failed to unmarshal unstructured object into Policy")
		return
	}

	for _, res := range clusterpolicy.Spec.ApplyOn {

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
		c.CacheIndex.Remove(
			Namespace("_ClusterScope"),
			Group(group),
			Version(version),
			Kind(res.Kind),
			PolicyName(clusterpolicy.Name),
		)
		c.CacheIndex.Unlock()
	}
	fmt.Println(c.CacheIndex.Policies)
}
