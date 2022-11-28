package loader

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// KubernetesClient
type Loader struct {
	KubeClient    *kubernetes.Clientset
	KubeDynClient dynamic.Interface
	KubeConfig    *rest.Config
	// TODO: Add other required data
}

// Policy resource definition
type Policy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PolicySpec `json:"spec,omitempty"`
}

type PolicySpec struct {
	DefaultBehaviour string `json:"defaultBehaviour"`
	ApplyOn          struct {
		ApiVersions []string `json:"apiVersions"`
		Kinds       []string `json:"kinds"`
	} `json:"applyOn"`
	Policies []struct {
		Name  string  `json:"name"`
		Rules []*Rule `json:"rules"`
	} `json:"policies"`
}

type Rule struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}
