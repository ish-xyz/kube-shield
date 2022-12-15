package cache

/*
This file defines the Cache Controller logic and methods
*/

import (
	"time"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	kcache "k8s.io/client-go/tools/cache"
)

func NewCacheController(clientset dynamic.Interface, c *CacheIndex) *Controller {

	clusterPolicy := schema.GroupVersionResource{
		Group:    v1.GroupVersion.Group,
		Version:  v1.GroupVersion.Version,
		Resource: v1.ClusterPolicyKind,
	}

	policy := schema.GroupVersionResource{
		Group:    v1.GroupVersion.Group,
		Version:  v1.GroupVersion.Version,
		Resource: v1.PolicyKind,
	}

	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(clientset, time.Minute, metav1.NamespaceAll, nil)

	return &Controller{
		ClusterInformer:   factory.ForResource(clusterPolicy).Informer(),
		NamespaceInformer: factory.ForResource(policy).Informer(),
		CacheIndex:        c,
	}
}

// Registers handlers for ClusterPolicies and Policies informers and start them
func (c *Controller) Run(polStopCh <-chan struct{}, clusterPolStopCh <-chan struct{}) {

	// Register handlers
	c.ClusterInformer.AddEventHandler(kcache.ResourceEventHandlerFuncs{
		AddFunc:    c.onClusterPolicyAdd,
		UpdateFunc: c.onClusterPolicyUpdate,
		DeleteFunc: c.onClusterPolicyDelete,
	})

	c.NamespaceInformer.AddEventHandler(kcache.ResourceEventHandlerFuncs{
		AddFunc:    c.onPolicyAdd,
		UpdateFunc: c.onPolicyUpdate,
		DeleteFunc: c.onPolicyDelete,
	})

	go c.ClusterInformer.Run(clusterPolStopCh)
	go c.NamespaceInformer.Run(polStopCh)

	for {
		select {
		case err := <-clusterPolStopCh:
			logrus.Fatalf("cluster policy informer failed %v", err)
		case err := <-polStopCh:
			logrus.Fatalf("namespace policy informer failed %v", err)
		}
	}
}

// Policy update event handler
func (c *Controller) onPolicyUpdate(oldObj interface{}, newObj interface{}) {

	var oldPolicy *v1.Policy
	var newPolicy *v1.Policy

	errOldPol := runtime.DefaultUnstructuredConverter.FromUnstructured(oldObj.(*unstructured.Unstructured).Object, &oldPolicy)
	errNewPol := runtime.DefaultUnstructuredConverter.FromUnstructured(newObj.(*unstructured.Unstructured).Object, &newPolicy)
	if errOldPol != nil || errNewPol != nil {
		logrus.Fatalf("failed to unmarshal unstructured object into Policy %v %v", errOldPol, errNewPol)
		return
	}

	// lock index, remove old policies and add new ones
	c.CacheIndex.Lock()
	for _, def := range oldPolicy.Spec.ApplyOn {
		c.CacheIndex.Delete(
			Verb(def.Verb),
			Namespace(oldPolicy.Namespace),
			getGroup(def.ApiGroup),
			Resource(def.Resource),
			PolicyName(oldPolicy.Name),
		)
	}
	for _, def := range newPolicy.Spec.ApplyOn {
		c.CacheIndex.Delete(
			Verb(def.Verb),
			Namespace(newPolicy.Namespace),
			getGroup(def.ApiGroup),
			Resource(def.Resource),
			PolicyName(newPolicy.Name),
		)
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

	for _, def := range policy.Spec.ApplyOn {
		c.CacheIndex.Lock()
		c.CacheIndex.Add(
			Verb(def.Verb),
			Namespace(policy.Namespace),
			GetGroup(def.ApiGroup),
			Resource(def.Resource),
			PolicyName(policy.Name),
		)
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

	for _, def := range policy.Spec.ApplyOn {
		c.CacheIndex.Lock()
		c.CacheIndex.Delete(
			Verb(def.Verb),
			Namespace(policy.Namespace),
			GetGroup(def.ApiGroup),
			Resource(def.Resource),
			PolicyName(policy.Name),
		)
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
		logrus.Fatalf("failed to unmarshal unstructured object into ClusterPolicy %v %v", errOldPol, errNewPol)
		return
	}

	// lock index, remove old policies and add new ones
	c.CacheIndex.Lock()
	for _, def := range oldPolicy.Spec.ApplyOn {
		c.CacheIndex.Delete(
			Verb(def.Verb),
			Namespace(oldPolicy.Namespace),
			GetGroup(def.ApiGroup),
			Resource(def.Resource),
			PolicyName(oldPolicy.Name),
		)
	}
	for _, def := range newPolicy.Spec.ApplyOn {
		c.CacheIndex.Delete(
			Verb(def.Verb),
			Namespace(newPolicy.Namespace),
			GetGroup(def.ApiGroup),
			Resource(def.Resource),
			PolicyName(newPolicy.Name),
		)
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

	for _, def := range clusterpolicy.Spec.ApplyOn {
		c.CacheIndex.Lock()
		c.CacheIndex.Add(
			Verb(def.Verb),
			Namespace("_ClusterScope"),
			GetGroup(def.ApiGroup),
			Resource(def.Resource),
			PolicyName(clusterpolicy.Name),
		)
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

	for _, def := range clusterpolicy.Spec.ApplyOn {
		c.CacheIndex.Lock()
		c.CacheIndex.Delete(
			Verb(def.Verb),
			Namespace("_ClusterScope"),
			GetGroup(def.ApiGroup),
			Resource(def.Resource),
			PolicyName(clusterpolicy.Name),
		)
		c.CacheIndex.Unlock()
	}
}
