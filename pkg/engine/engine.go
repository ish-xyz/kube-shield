package engine

import (
	"fmt"
	"strings"

	"github.com/RedLabsPlatform/kube-shield/pkg/cache"
	admissionv1 "k8s.io/api/admission/v1"
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
