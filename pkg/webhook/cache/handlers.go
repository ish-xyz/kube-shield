package cache

import (
	"strings"

	"github.com/RedLabsPlatform/kube-shield/pkg/webhook/engine"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func getGV(info string) (Group, Version) {

	gv := strings.Split(info, "/")
	group := "_core"
	version := info
	if len(gv) > 1 {
		group = gv[0]
		version = gv[1]
	}

	return Group(group), Version(version)
}

func (c *CacheController) onPolicyUpdate(oldObj interface{}, newObj interface{}) {

	var oldPolicy *engine.Policy
	var newPolicy *engine.Policy

	errOldPol := runtime.DefaultUnstructuredConverter.FromUnstructured(oldObj.(*unstructured.Unstructured).Object, &oldPolicy)
	errNewPol := runtime.DefaultUnstructuredConverter.FromUnstructured(newObj.(*unstructured.Unstructured).Object, &newPolicy)
	if errOldPol != nil || errNewPol != nil {
		logrus.Fatal("failed to unmarshal unstructured object into Policy %v %v", errOldPol, errNewPol)
		return
	}

	// lock index, remove old policies and add new ones
	c.CacheIndex.Lock()
	for _, res := range oldPolicy.Spec.ApplyOn {
		group, version := getGV(res.APIVersion)
		c.CacheIndex.Delete(Namespace(oldPolicy.Namespace), group, version, Kind(res.Kind), PolicyName(oldPolicy.Name))
	}
	for _, res := range newPolicy.Spec.ApplyOn {
		group, version := getGV(res.APIVersion)
		c.CacheIndex.Add(Namespace(newPolicy.Namespace), group, version, Kind(res.Kind), PolicyName(newPolicy.Name))
	}
	c.CacheIndex.Unlock()
}

func (c *CacheController) onPolicyAdd(obj interface{}) {

	var policy *engine.Policy
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(*unstructured.Unstructured).Object, &policy)
	if err != nil {
		logrus.Fatal("failed to unmarshal unstructured object into Policy")
		return
	}

	for _, res := range policy.Spec.ApplyOn {
		group, version := getGV(res.APIVersion)
		c.CacheIndex.Lock()
		c.CacheIndex.Add(Namespace(policy.Namespace), group, version, Kind(res.Kind), PolicyName(policy.Name))
		c.CacheIndex.Unlock()
	}
}

func (c *CacheController) onPolicyDelete(obj interface{}) {
	var policy *engine.Policy
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(*unstructured.Unstructured).Object, &policy)
	if err != nil {
		logrus.Fatal("failed to unmarshal unstructured object into Policy")
		return
	}

	for _, res := range policy.Spec.ApplyOn {
		group, version := getGV(res.APIVersion)
		c.CacheIndex.Lock()
		c.CacheIndex.Delete(Namespace(policy.Namespace), group, version, Kind(res.Kind), PolicyName(policy.Name))
		c.CacheIndex.Unlock()
	}
}

func (c *CacheController) onClusterPolicyUpdate(oldObj interface{}, newObj interface{}) {

	var oldPolicy *engine.ClusterPolicy
	var newPolicy *engine.ClusterPolicy

	errOldPol := runtime.DefaultUnstructuredConverter.FromUnstructured(oldObj.(*unstructured.Unstructured).Object, &oldPolicy)
	errNewPol := runtime.DefaultUnstructuredConverter.FromUnstructured(newObj.(*unstructured.Unstructured).Object, &newPolicy)
	if errOldPol != nil || errNewPol != nil {
		logrus.Fatal("failed to unmarshal unstructured object into ClusterPolicy %v %v", errOldPol, errNewPol)
		return
	}

	// lock index, remove old policies and add new ones
	c.CacheIndex.Lock()
	for _, res := range oldPolicy.Spec.ApplyOn {
		group, version := getGV(res.APIVersion)
		c.CacheIndex.Delete(Namespace("_ClusterScope"), group, version, Kind(res.Kind), PolicyName(oldPolicy.Name))
	}
	for _, res := range newPolicy.Spec.ApplyOn {
		group, version := getGV(res.APIVersion)
		c.CacheIndex.Add(Namespace("_ClusterScope"), group, version, Kind(res.Kind), PolicyName(newPolicy.Name))
	}
	c.CacheIndex.Unlock()
}

func (c *CacheController) onClusterPolicyAdd(obj interface{}) {
	var clusterpolicy *engine.ClusterPolicy
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(*unstructured.Unstructured).Object, &clusterpolicy)
	if err != nil {
		logrus.Fatal("failed to unmarshal unstructured object into ClusterPolicy")
		return
	}

	for _, res := range clusterpolicy.Spec.ApplyOn {
		group, version := getGV(res.APIVersion)
		c.CacheIndex.Lock()
		c.CacheIndex.Add(Namespace("_ClusterScope"), group, version, Kind(res.Kind), PolicyName(clusterpolicy.Name))
		c.CacheIndex.Unlock()
	}
}

func (c *CacheController) onClusterPolicyDelete(obj interface{}) {

	var clusterpolicy *engine.ClusterPolicy
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(*unstructured.Unstructured).Object, &clusterpolicy)
	if err != nil {
		logrus.Errorln("failed to unmarshal unstructured object into ClusterPolicy")
		return
	}

	for _, res := range clusterpolicy.Spec.ApplyOn {
		group, version := getGV(res.APIVersion)
		c.CacheIndex.Lock()
		c.CacheIndex.Delete(Namespace("_ClusterScope"), group, version, Kind(res.Kind), PolicyName(clusterpolicy.Name))
		c.CacheIndex.Unlock()
	}
}
