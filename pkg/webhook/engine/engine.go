package engine

import (
	"fmt"

	"github.com/RedLabsPlatform/kube-shield/pkg/webhook/cache"
	admissionv1 "k8s.io/api/admission/v1"
)

func (e *Engine) RunNamespacedPolicies(payload *admissionv1.AdmissionReview) {

	index := e.CacheController.CacheIndex
	ns := cache.Namespace(payload.Request.Namespace)
	group, version := e.CacheController.GetGV(payload.Request.Kind.Group + payload.Request.Kind.Version)
	kind := cache.Kind(payload.Request.Kind.Kind)
	store := e.CacheController.NamespaceInformer.GetStore()

	fmt.Println(ns, group, version, kind)
	fmt.Println(index)

	for _, v := range index.Get(ns, group, version, kind) {
		fmt.Println(v)
		obj, exists, err := store.GetByKey(string(v))
		fmt.Println("here before")
		if err != nil {
			e.Logger.Warnln("failed to get policy with name '%v', error: '%v'", v, err)
		}

		if !exists {
			e.Logger.Warnln("object %v is cached in index but the actual resource doesn't exists", v)
			continue
		}
		fmt.Println(obj, exists)
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
