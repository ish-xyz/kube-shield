package engine

import (
	"fmt"

	"github.com/RedLabsPlatform/kube-shield/pkg/webhook/cache"
	admissionv1 "k8s.io/api/admission/v1"
)

func (e *Engine) Run(payload *admissionv1.AdmissionReview) {

	//Policies map[Namespace]map[Group]map[Version]map[Kind][]PolicyName

	index := e.CacheController.CacheIndex
	ns := cache.Namespace(payload.Request.Namespace)
	group, version := e.CacheController.GetGV(payload.APIVersion)
	kind := cache.Kind(payload.Kind)
	store := e.CacheController.NamespaceInformer.GetStore()

	for _, v := range index.Get(ns, group, version, kind) {
		fmt.Println(v)
		obj, exists, err := store.GetByKey(string(v))
		if err != nil {
			e.Logger.Warnln("failed to get policy with name '%v', error: '%v'", v, err)
		}
		if !exists {
			e.Logger.Warnln("object %v is cached in index but the actual resource doesn't exists", v)
			continue
		}
		fmt.Println(obj)
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
