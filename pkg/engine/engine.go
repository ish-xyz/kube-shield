package engine

import (
	"encoding/json"
	"fmt"
	"strings"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/RedLabsPlatform/kube-shield/pkg/cache"
	"github.com/RedLabsPlatform/kube-shield/pkg/engine/lua"
	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Run policies: Namespaced, Cluster, Global policies
func (e *Engine) Run(req *admissionv1.AdmissionRequest) error {

	store := e.CacheController.NamespaceInformer.GetStore()
	operation := cache.Operation(strings.ToLower(string(req.Operation)))
	ns := cache.Namespace(req.Namespace)
	group := cache.GetGroup(req.Resource.Group)
	res := cache.GetResource(req.RequestResource.Resource, req.SubResource)
	policies := e.CacheController.CacheIndex.Get(operation, ns, group, res)
	policies = append(policies, e.CacheController.CacheIndex.Get(operation, cache.CLUSTER_SCOPE, group, res)...)

	for _, name := range policies {

		spec := &v1.PolicySpec{}
		key := strings.TrimPrefix(fmt.Sprintf("%s/%s", req.Namespace, name), "/")
		obj, exists, err := store.GetByKey(key)
		if !exists || obj == nil {
			logrus.Warningln("cache mismatch error", key)
			//TODO: add metrics "cache_mismatch{policy_name="", etc}"
		}
		if err != nil {
			logrus.Errorf("failed to get policy with key '%s'", key)
			//TODO: should try to reconcile and clean the index or wait for the cache to refresh
			//TODO: add metrics "failed_loading{policy_name="", etc}"
			continue
		}

		err = runtime.DefaultUnstructuredConverter.FromUnstructured(
			obj.(*unstructured.Unstructured).Object["spec"].(map[string]interface{}),
			&spec,
		)
		if err != nil {
			logrus.Errorf("failed to convert policy with key '%s' into object", key)
			//TODO: add metrics "failed_loading{policy_name="", etc}"
			continue
		}

		err = evaluateRules(req, spec.Rules)
		if err != nil {
			return fmt.Errorf("\n\nDenied by policy: '%s'\n%v", name, err)
		}
	}

	return nil
}

// Rules are in OR, so if any of the rules have passed, the function returns a nil error
func evaluateRules(req *admissionv1.AdmissionRequest, rules []*v1.Rule) error {

	var lastRule string

	for _, rule := range rules {
		lastRule = rule.Name
		jsonReq, err := json.Marshal(req)
		if err != nil {
			return err
		}

		res, err := lua.Execute(string(jsonReq), rule.Script)
		if !res {
			return fmt.Errorf("rule: '%s'\nerror: '%v'\n ", lastRule, err)
		}
	}

	return nil
}
