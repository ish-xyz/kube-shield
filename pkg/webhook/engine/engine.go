package engine

import (
	"fmt"

	admissionv1 "k8s.io/api/admission/v1"
)

func Load(resource *admissionv1.AdmissionReview) {
	fmt.Println("TODO")

	/*
		Load all cluster policies first from index
		Load all namespaced policies

	*/
}
