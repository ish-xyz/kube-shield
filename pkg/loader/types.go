package loader

/*
  Use controller-gen to generate the CRDs
  https://github.com/kubernetes-sigs/controller-tools
*/

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +resource:path=policy
// +kubebuilder:object:root=true

type Policy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// +required
	Spec PolicySpec `json:"spec,omitempty"`
}

type PolicySpec struct {
	// +kubebuilder:validation:Enum=AllowAll;DenyAll
	DefaultBehaviour string `json:"defaultBehaviour"`
	ApplyOn          struct {
		// +kubebuilder:validation:UniqueItems=true
		ApiVersions []string `json:"apiVersions"`
		// +kubebuilder:validation:UniqueItems=true
		Kinds []string `json:"kinds"`
	} `json:"applyOn"`
	Policies []struct {
		Name string `json:"name"`
		// +kubebuilder:validation:MinItems=1
		// +kubebuilder:validation:MaxItems=500
		Rules []*Rule `json:"rules"`
	} `json:"policies"`
}

type Rule struct {
	Field string `json:"field"`
	// +kubebuilder:validation:Enum=GreaterThan;LowerThan;Equal;NotEqual;Regex
	Operator string `json:"operator"`
	Value    string `json:"value"`
}
