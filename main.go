package main

import (
	"github.com/RedLabsPlatform/kube-shield/cmd"
)

//DISABLED go:generate go run sigs.k8s.io/controller-tools/cmd/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./pkg/apis/v1/..."
//DISABLED go:generate go run sigs.k8s.io/controller-tools/cmd/controller-gen crd paths=./pkg/apis/v1/... crd:crdVersions=v1 output:crd:dir=./crds

func main() {
	cmd.Execute()
}
