package engine

import (
	"fmt"

	"github.com/RedLabsPlatform/kube-shield/pkg/webhook/cache"
	admissionv1 "k8s.io/api/admission/v1"
)

func (e *Engine) Load(payload *admissionv1.AdmissionReview) {

	//Policies map[Namespace]map[Group]map[Version]map[Kind][]PolicyName

	index := e.CacheController.CacheIndex.Policies
	ns := cache.Namespace(payload.Request.Namespace)
	group, version := e.CacheController.GetGV(payload.APIVersion)
	kind := cache.Kind(payload.Kind)

	e.Logger.Debugln("policies for payload")
	e.Logger.Debugln(fmt.Println(index[ns][group][version][kind]))
	/*
		Load all cluster policies first from index
		Load all namespaced policies
	*/
}
