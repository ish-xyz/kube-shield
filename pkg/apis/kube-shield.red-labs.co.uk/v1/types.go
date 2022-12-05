package engine

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Policy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PolicySpec `json:"spec,omitempty"`
}

type PolicySpec struct {
	DefaultBehaviour string             `json:"defaultBehaviour"`
	ApplyOn          []*ResourceAddress `json:"applyOn"`
	Policies         []struct {
		Name  string  `json:"name"`
		Rules []*Rule `json:"rules"`
	} `json:"policies"`
}

type ClusterPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ClusterPolicySpec `json:"spec,omitempty"`
}

type ClusterPolicySpec struct {
	DefaultBehaviour string             `json:"defaultBehaviour"`
	ApplyOn          []*ResourceAddress `json:"applyOn"`
	Policies         []struct {
		Name  string  `json:"name"`
		Rules []*Rule `json:"rules"`
	} `json:"policies"`
}

type ResourceAddress struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
}

type Rule struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type Engine struct{}
