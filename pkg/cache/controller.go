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
		Resource: v1.CLUSTER_POLICY_KIND,
	}

	policy := schema.GroupVersionResource{
		Group:    v1.GroupVersion.Group,
		Version:  v1.GroupVersion.Version,
		Resource: v1.POLICY_KIND,
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
	c.CacheIndex.Delete(oldPolicy.Spec.ApplyOn, oldPolicy.Namespace, oldPolicy.Name)
	c.CacheIndex.Add(newPolicy.Spec.ApplyOn, newPolicy.Namespace, newPolicy.Name)
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

	c.CacheIndex.Lock()
	c.CacheIndex.Add(policy.Spec.ApplyOn, policy.Namespace, policy.Name)
	c.CacheIndex.Unlock()
}

// Policy delete event handler
func (c *Controller) onPolicyDelete(obj interface{}) {
	var policy *v1.Policy
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(*unstructured.Unstructured).Object, &policy)
	if err != nil {
		logrus.Fatal("failed to unmarshal unstructured object into Policy")
		return
	}

	c.CacheIndex.Lock()
	c.CacheIndex.Delete(policy.Spec.ApplyOn, policy.Namespace, policy.Name)
	c.CacheIndex.Unlock()
}

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
	c.CacheIndex.Delete(oldPolicy.Spec.ApplyOn, CLUSTER_SCOPE, oldPolicy.Name)
	c.CacheIndex.Add(newPolicy.Spec.ApplyOn, CLUSTER_SCOPE, newPolicy.Name)
	c.CacheIndex.Unlock()
}

// ClusterPolicy add event handler
func (c *Controller) onClusterPolicyAdd(obj interface{}) {
	var policy *v1.ClusterPolicy
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(*unstructured.Unstructured).Object, &policy)
	if err != nil {
		logrus.Fatal("failed to unmarshal unstructured object into ClusterPolicy")
		return
	}

	c.CacheIndex.Lock()
	c.CacheIndex.Add(policy.Spec.ApplyOn, CLUSTER_SCOPE, policy.Name)
	c.CacheIndex.Unlock()
}

// ClusterPolicy delete event handler
func (c *Controller) onClusterPolicyDelete(obj interface{}) {

	var policy *v1.ClusterPolicy
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(*unstructured.Unstructured).Object, &policy)
	if err != nil {
		logrus.Errorln("failed to unmarshal unstructured object into ClusterPolicy")
		return
	}

	c.CacheIndex.Lock()
	c.CacheIndex.Delete(policy.Spec.ApplyOn, CLUSTER_SCOPE, policy.Name)
	c.CacheIndex.Unlock()
}
