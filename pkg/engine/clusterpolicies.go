package engine

import (
	"strings"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/RedLabsPlatform/kube-shield/pkg/cache"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

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
