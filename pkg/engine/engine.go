package engine

import (
	"encoding/json"
	"fmt"
	"strings"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/RedLabsPlatform/kube-shield/pkg/cache"
	"github.com/RedLabsPlatform/kube-shield/pkg/engine/operators"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func runChecks(req *admissionv1.AdmissionRequest, desiredMatchResult bool, checks []*v1.Check) (bool, error) {

	for _, check := range checks {
		jsonReq, err := json.Marshal(req)
		if err != nil {
			return false, err
		}
		checkRes, err := operators.Run(string(jsonReq), check)
		if checkRes.Match != desiredMatchResult {
			return checkRes.Match, fmt.Errorf("%s", checkRes.Message)
		}
	}
	return true, nil
}

func runRules(req *admissionv1.AdmissionRequest, policy *v1.Policy) {

	/*
		* AllowIfMatch  == AllowIfMatch (default) -> true
		  (if all checks are true, pass)

		* DenyIfMatch  == AllowIfMatch (default) -> false
		  (if all checks are false, pass)
	*/
	desiredMatchResult := policy.Spec.DefaultBehaviour == v1.DEFAULT_BEHAVIOUR

	for _, rule := range policy.Spec.Rules {

		runChecks(req, desiredMatchResult, rule.Checks)

	}
}

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

		runRules(req, policy)

	}
}
