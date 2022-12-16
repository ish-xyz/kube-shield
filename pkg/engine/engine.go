package engine

import (
	"fmt"
	"strings"

	"github.com/RedLabsPlatform/kube-shield/pkg/cache"
	admissionv1 "k8s.io/api/admission/v1"
)

func (e *Engine) RunNamespacedPolicies(payload *admissionv1.AdmissionReview) {

	index := e.CacheController.CacheIndex
	store := e.CacheController.NamespaceInformer.GetStore()
	req := payload.Request

	verb := cache.Verb(strings.ToLower(string(req.Operation)))
	ns := cache.Namespace(req.Namespace)
	group := cache.GetGroup(req.Resource.Group)
	res := cache.GetResource(req.RequestResource.Resource, req.SubResource)

	fmt.Println(verb, ns, group, res)
	for _, v := range index.Get(verb, ns, group, res) {
		fmt.Println(v)
		obj, exists, err := store.GetByKey(string(v))
		fmt.Println(obj, exists, err)
		if err != nil {
			e.Logger.Warnln("failed to get policy with name '%v', error: '%v'", v, err)
		}

		if !exists {
			e.Logger.Warnln("object %v is cached in index but the actual resource doesn't exists", v)
			continue
		}
	}

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
