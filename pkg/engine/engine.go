package engine

import (
	"fmt"
	"strings"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/RedLabsPlatform/kube-shield/pkg/cache"
	"github.com/RedLabsPlatform/kube-shield/pkg/engine/operators"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func (e *Engine) RunNamespacePolicies(req *admissionv1.AdmissionRequest) {

	store := e.CacheController.NamespaceInformer.GetStore()
	verb := cache.Verb(strings.ToLower(string(req.Operation)))
	ns := cache.Namespace(req.Namespace)
	group := cache.GetGroup(req.Resource.Group)
	res := cache.GetResource(req.RequestResource.Resource, req.SubResource)

	for _, name := range e.CacheController.CacheIndex.Get(verb, ns, group, res) {
		var policy *v1.Policy
		policyKey := fmt.Sprintf("%s/%s", ns, name)
		obj, exists, err := store.GetByKey(policyKey)

		if err != nil || !exists {
			e.Logger.Errorf("failed to get policy with key '%s'", policyKey)
			continue
		}

		err = runtime.DefaultUnstructuredConverter.FromUnstructured(
			obj.(*unstructured.Unstructured).Object,
			&policy,
		)
		if err != nil {
			e.Logger.Errorf("failed to convert policy with key '%s' into object", policyKey)
			continue
		}
		for _, rule := range policy.Spec.Rules {
			for _, check := range rule.Checks {
				reqObject := string(req.Object.Raw)
				checkResult, err := operators.Run(reqObject, &check)
				fmt.Printf("%+v\n", checkResult)
				fmt.Println(err)
				fmt.Println(reqObject)
			}

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
