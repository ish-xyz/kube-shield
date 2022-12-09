package v1

type ResourceAddress struct {
	// +kubebuilder:validation:Required
	APIVersion string `json:"apiVersion"`
	// +kubebuilder:validation:Required
	Kind string `json:"kind"`
}

type Check struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=^\$*$
	Field string `json:"field"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=GreaterThan;LowerThan;Equal;NotEqual;Regex;Count
	Operator string `json:"operator"`
	// +kubebuilder:validation:Required
	Value interface{} `json:"value"`
}

type CheckResult struct {
	Result bool    `json:"result"`
	Errors []error `json:"errors"`
}
