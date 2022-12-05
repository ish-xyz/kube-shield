package cache

import (
	"strings"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/kube-shield.red-labs.co.uk/v1"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Extract API Group & Version from apiVersion field
// the core api group, which is usually equal to "" (empty) is set as "_core" here
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

// ** Policies scope

// Policy update event handler
func (c *Controller) onPolicyUpdate(oldObj interface{}, newObj interface{}) {

	var oldPolicy *v1.Policy
	var newPolicy *v1.Policy

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

// Policy add event handler
func (c *Controller) onPolicyAdd(obj interface{}) {

	var policy *v1.Policy
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

// Policy delete event handler
func (c *Controller) onPolicyDelete(obj interface{}) {
	var policy *v1.Policy
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

// ** Cluster scope policies

// ClusterPolicy update event handler
func (c *Controller) onClusterPolicyUpdate(oldObj interface{}, newObj interface{}) {

	var oldPolicy *v1.ClusterPolicy
	var newPolicy *v1.ClusterPolicy

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

// ClusterPolicy add event handler
func (c *Controller) onClusterPolicyAdd(obj interface{}) {
	var clusterpolicy *v1.ClusterPolicy
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

// ClusterPolicy delete event handler
func (c *Controller) onClusterPolicyDelete(obj interface{}) {

	var clusterpolicy *v1.ClusterPolicy
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
