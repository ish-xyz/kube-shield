package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
type Policy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PolicySpec   `json:"spec,omitempty"`
	Status            PolicyStatus `json:"status,omitempty"`
}

type PolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

type PolicySpec struct {
	// +kubebuilder:default=IfMatchDeny
	// +kubebuilder:validation:Enum=IfMatchAllow;IfMatchDeny
	DefaultBehaviour string `json:"defaultBehaviour"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxItems=500
	// +kubebuilder:validation:MinItems=1
	ApplyOn []*Resource `json:"applyOn"`
	// +kubebuilder:validation:MaxItems=500
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	Rules []*Rule `json:"rules"`
}

type Rule struct {
	// +kubebuilder:validation:Pattern=^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxItems=500
	// +kubebuilder:validation:MinItems=1
	Checks []Check `json:"checks"`
}

// +kubebuilder:object:root=true
// PolicyList contains a list of Policy
type PolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Policy `json:"items"`
}

type Resource struct {
	// +kubebuilder:validation:Required
	apiGroups []string `json:"apiGroups"`
	// +kubebuilder:validation:Required
	resources []string `json:"resources"`
	// +kubebuilder:validation:Required
	verbs []string `json:"verbs"`
}

type Check struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^\$_\..*$
	Field string `json:"field"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=GreaterThan;LowerThan;Equal;NotEqual;Regex;Count
	Operator string `json:"operator"`

	// +kubebuilder:validation:Required
	Value interface{} `json:"value"`
}

type CheckResult struct {
	Result  bool
	Message string
}

func init() {
	SchemeBuilder.Register(&Policy{}, &PolicyList{})
}
