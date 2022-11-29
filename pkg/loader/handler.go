package loader

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type PolicyHandler struct{}

func NewPolicyHandler() *PolicyHandler {
	return &PolicyHandler{}
}

func (p PolicyHandler) OnAdd(obj interface{}) {
	policy, ok := obj.(*unstructured.Unstructured)
	if !ok {
		logrus.Warningln("Unable to parse policy")
	}
	fmt.Println(policy)
	// Add On Add Logic here. We could inject it on the constructor function and call it from here
}

func (p PolicyHandler) OnUpdate(oldObj, newObj interface{}) {
	policy, ok := newObj.(*unstructured.Unstructured)
	if !ok {
		logrus.Warningln("Unable to parse policy")
	}
	fmt.Println(policy)
	// Add On Add Logic here. We could inject it on the constructor function and call it from here
}

func (p PolicyHandler) OnDelete(obj interface{}) {
	policy, ok := obj.(*unstructured.Unstructured)
	if !ok {
		logrus.Warningln("Unable to parse policy")
	}
	fmt.Println(policy)
	// Add On Add Logic here. We could inject it on the constructor function and call it from here
}
