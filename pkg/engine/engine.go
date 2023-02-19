package engine

import (
	"encoding/json"
	"fmt"

	v1 "github.com/RedLabsPlatform/kube-shield/pkg/apis/v1"
	"github.com/RedLabsPlatform/kube-shield/pkg/engine/lua"
	admissionv1 "k8s.io/api/admission/v1"
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

// Run policies: Namespaced, Cluster, Global policies
func (e *Engine) Run(req *admissionv1.AdmissionRequest) error {
	err := e.RunClusterPolicies(req)
	if err != nil {
		return err
	}
	err = e.RunNamespacePolicies(req)
	if err != nil {
		return err
	}
	return nil
}
