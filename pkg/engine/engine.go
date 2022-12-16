package engine

import (
	"fmt"
	"strings"

	"github.com/RedLabsPlatform/kube-shield/pkg/cache"
	admissionv1 "k8s.io/api/admission/v1"
)

func (e *Engine) RunNamespacedPolicies(payload *admissionv1.AdmissionReview) {

	index := e.CacheController.CacheIndex
	namespacedStore := e.CacheController.NamespaceInformer.GetStore()
	clusterStore := e.CacheController.ClusterInformer.GetStore()

	req := payload.Request
	verb := cache.Verb(strings.ToLower(string(req.Operation)))
	ns := cache.Namespace(req.Namespace)
	clusterNs := cache.Namespace(cache.CLUSTER_NAMESPACE)
	group := cache.GetGroup(req.Resource.Group)
	res := cache.GetResource(req.RequestResource.Resource, req.SubResource)

	for _, name := range index.Get(verb, ns, group, res) {
		policyKey := fmt.Sprintf("%s/%s", ns, name)
		obj, exists, err := namespacedStore.GetByKey(policyKey)
		fmt.Println(obj, exists, err)
	}

	// TODO
	fmt.Println(verb, ns, group, res)
	fmt.Println(clusterStore.ListKeys())
	fmt.Println(clusterNs)

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
