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
	// +kubebuilder:default=AllowIfMatch
	// +kubebuilder:validation:Enum=AllowIfMatch;DenyIfMatch
	DefaultBehaviour string `json:"defaultBehaviour"`
	// +kubebuilder:validation:Required
	ApplyOn []*Definition `json:"applyOn"`
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
	Checks []*Check `json:"checks"`
}

// +kubebuilder:object:root=true
// PolicyList contains a list of Policy
type PolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Policy `json:"items"`
}

type Definition struct {
	// +kubebuilder:validation:Required
	Group string `json:"apiGroup"`
	// +kubebuilder:validation:Required
	Resource string `json:"resource"`
	// +kubebuilder:validation:Required
	Verb string `json:"verb"`
}

type Check struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^\$_\..*$
	Field string `json:"field"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=GreaterThan;LowerThan;Equal;NotEqual;Regex;Count
	Operator string `json:"operator"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XPreserveUnknownFields
	Value interface{} `json:"value"`
}

type CheckResult struct {
	Status int
	Match  bool
	Error  error
}

func init() {
	SchemeBuilder.Register(&Policy{}, &PolicyList{})
}
