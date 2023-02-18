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
func runRules(req *admissionv1.AdmissionRequest, name string, rules []*v1.Rule) error {

	var lastRule string

	for _, rule := range rules {
		lastRule = rule.Name
		jsonReq, err := json.Marshal(req)
		if err != nil {
			return err
		}

		res, err := lua.Execute(string(jsonReq), rule.Script)
		if !res {
			return fmt.Errorf("\n\nDenied by policy: '%s'\nrule: '%s'\nerror: '%v'\n ", name, lastRule, err)
		}
	}

	return nil
}

// Run Cluster Policies
func (e *Engine) RunClusterPolicies(req *admissionv1.AdmissionRequest) error {

	store := e.CacheController.ClusterInformer.GetStore()
	operation := cache.Operation(strings.ToLower(string(req.Operation)))
	group := cache.GetGroup(req.Resource.Group)
	res := cache.GetResource(req.RequestResource.Resource, req.SubResource)

	for _, name := range e.CacheController.CacheIndex.Get(operation, cache.CLUSTER_SCOPE, group, res) {
		var policy *v1.ClusterPolicy

		obj, exists, err := store.GetByKey(string(name))
		if err != nil || !exists {
			e.Logger.Errorf("failed to get policy with key '%s'", name)
			//TODO: I think it should exit here
			continue
		}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(
			obj.(*unstructured.Unstructured).Object,
			&policy,
		)
		if err != nil {
			e.Logger.Errorf("failed to convert policy with key '%s' into object", name)
			//TODO: I think it should exit here
			continue
		}

		err = runRules(req, policy.Name, policy.Spec.Rules)
		if err != nil {
			return err
		}
	}
	return nil
}

// Run Policies
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
			//TODO: should try to reconcile and clean the index or wait for the cache to refresh
			//TODO: add metrics "failed_loading{policy_name="", etc}"
			continue
		}

		err = runtime.DefaultUnstructuredConverter.FromUnstructured(
			obj.(*unstructured.Unstructured).Object,
			&policy,
		)
		if err != nil {
			e.Logger.Errorf("failed to convert policy with key '%s' into object", policyKey)
			//TODO: add metrics "failed_loading{policy_name="", etc}"
			continue
		}

		err = runRules(req, policy.Name, policy.Spec.Rules)
		if err != nil {
			return err
		}

	}

	return nil
}
