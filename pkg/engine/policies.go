package engine

import (
	"fmt"
	"strings"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/RedLabsPlatform/kube-shield/pkg/cache"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

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
