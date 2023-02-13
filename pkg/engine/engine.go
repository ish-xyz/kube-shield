package engine

import (
	"encoding/json"
	"fmt"
	"strings"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/RedLabsPlatform/kube-shield/pkg/cache"
	"github.com/RedLabsPlatform/kube-shield/pkg/engine/lua"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Rules are in OR, so if any of the rules have passed, the function returns a nil error
func runRules(req *admissionv1.AdmissionRequest, policy *v1.Policy) error {

	/*
		* AllowIfMatch  == AllowIfMatch (default) -> true
		  (if all checks are true, pass)

		* DenyIfMatch  == AllowIfMatch (default) -> false
		  (if all checks are false, pass)
	*/
	var lastRule string

	for _, rule := range policy.Spec.Rules {
		lastRule = rule.Name
		jsonReq, err := json.Marshal(req)
		if err != nil {
			return err
		}

		res, err := lua.Execute(string(jsonReq), rule.Script)
		if !res {
			return fmt.Errorf("\n\nDenied by policy: '%s'\nrule: '%s'\nerror: '%v'\n ", policy.Name, lastRule, err)
		}
	}

	return nil
}

func (e *Engine) RunNamespacePolicies(req *admissionv1.AdmissionRequest) error {

	store := e.CacheController.NamespaceInformer.GetStore()
	operation := cache.Operation(strings.ToLower(string(req.Operation)))
	ns := cache.Namespace(req.Namespace)
	group := cache.GetGroup(req.Resource.Group)
	res := cache.GetResource(req.RequestResource.Resource, req.SubResource)

	for _, name := range e.CacheController.CacheIndex.Get(operation, ns, group, res) {

		var policy *v1.Policy

		policyKey := fmt.Sprintf("%s/%s", ns, name)
		obj, exists, err := store.GetByKey(policyKey)
		if err != nil || !exists {
			e.Logger.Errorf("failed to get policy with key '%s'", policyKey)
			//TODO: I think it should exit here
			continue
		}

		err = runtime.DefaultUnstructuredConverter.FromUnstructured(
			obj.(*unstructured.Unstructured).Object,
			&policy,
		)
		if err != nil {
			e.Logger.Errorf("failed to convert policy with key '%s' into object", policyKey)
			//TODO: I think it should exit here
			continue
		}

		err = runRules(req, policy)
		if err != nil {
			return err
		}

	}

	return nil
}
