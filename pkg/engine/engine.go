package engine

import (
	"encoding/json"
	"fmt"
	"strings"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/RedLabsPlatform/kube-shield/pkg/cache"
	"github.com/RedLabsPlatform/kube-shield/pkg/engine/operators"
	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Checks are in AND, so if any of the checks don't match the desired result it returns an error
// Nil error means that the checks have been successful
func runChecks(req *admissionv1.AdmissionRequest, desiredResult bool, checks []*v1.Check) error {

	for _, check := range checks {
		jsonReq, err := json.Marshal(req)
		if err != nil {
			return err
		}
		checkRes := operators.Run(string(jsonReq), check)
		if checkRes.Error != nil {
			logrus.Errorf("engine failed while executing check %v", checkRes.Error)
			return fmt.Errorf("engine failed while executing check %v", checkRes.Error)
		}
		if checkRes.Result != desiredResult {
			return checkRes.Error
		}
	}
	return nil
}

// Rules are in OR, so if any of the rules have passed, the function returns a nil error
func runRules(req *admissionv1.AdmissionRequest, policy *v1.Policy) error {

	/*
		* AllowIfMatch  == AllowIfMatch (default) -> true
		  (if all checks are true, pass)

		* DenyIfMatch  == AllowIfMatch (default) -> false
		  (if all checks are false, pass)
	*/
	var err error
	desiredResult := policy.Spec.DefaultBehaviour == v1.DEFAULT_BEHAVIOUR

	for _, rule := range policy.Spec.Rules {

		err = runChecks(req, desiredResult, rule.Checks)
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("\nDenied by policy '%s'.\nbehaviour is '%s', %v", policy.Name, policy.Spec.DefaultBehaviour, err)
}

func (e *Engine) RunNamespacePolicies(req *admissionv1.AdmissionRequest) error {

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

		err = runRules(req, policy)
		if err != nil {
			return err
		}

	}

	return nil
}
