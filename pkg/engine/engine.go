package engine

import (
	"fmt"
	"strings"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/RedLabsPlatform/kube-shield/pkg/cache"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func (e *Engine) RunClusterPolicies(payload *admissionv1.AdmissionReview) {

	index := e.CacheController.CacheIndex
	store := e.CacheController.ClusterInformer.GetStore()

	req := payload.Request
	verb := cache.Verb(strings.ToLower(string(req.Operation)))
	ns := cache.Namespace(cache.CLUSTER_SCOPE)
	group := cache.GetGroup(req.Resource.Group)
	res := cache.GetResource(req.RequestResource.Resource, req.SubResource)

	for _, name := range index.Get(verb, ns, group, res) {
		policyKey := fmt.Sprintf("%s/%s", ns, name)
		obj, exists, err := store.GetByKey(policyKey)
		fmt.Println(obj, exists, err)
	}

	// TODO
	fmt.Println(verb, ns, group, res)

	/*
		for each clusterpolicy
			get policy object from cache
			load into policy
			for each check
		for each policy
			get policy object from cache
		Load all namespaced policies
	*/
}

func (e *Engine) RunNamespacePolicies(payload *admissionv1.AdmissionReview) {

	index := e.CacheController.CacheIndex
	store := e.CacheController.NamespaceInformer.GetStore()

	req := payload.Request
	verb := cache.Verb(strings.ToLower(string(req.Operation)))
	ns := cache.Namespace(req.Namespace)
	group := cache.GetGroup(req.Resource.Group)
	res := cache.GetResource(req.RequestResource.Resource, req.SubResource)

	for _, name := range index.Get(verb, ns, group, res) {
		var policy *v1.Policy
		policyKey := fmt.Sprintf("%s/%s", ns, name)
		obj, exists, err := store.GetByKey(policyKey)

		if err != nil || !exists {
			e.Logger.Errorf("failed to get policy with key '%s'", policyKey)
			continue
		}

		err = runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(*unstructured.Unstructured).Object, &policy)
		if err != nil {
			e.Logger.Errorf("failed to convert policy with key '%s' into object", policyKey)
			continue
		}
		for _, rule := range policy.Spec.Rules {
			fmt.Println(rule)
		}
	}

	// TODO
	fmt.Println(verb, ns, group, res)

	/*
		for each clusterpolicy
			get policy object from cache
			load into policy
			for each check
		for each policy
			get policy object from cache
		Load all namespaced policies
	*/
}
