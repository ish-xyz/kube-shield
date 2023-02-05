package main

import (
	"github.com/RedLabsPlatform/kube-shield/cmd"
)

//go:generate echo "!!! IMPORTANT have you removed all interfaces from the CRDs structs?"
//go:generate sleep 3
//go:generate go run sigs.k8s.io/controller-tools/cmd/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./pkg/apis/v1/..."
//go:generate go run sigs.k8s.io/controller-tools/cmd/controller-gen crd paths=./pkg/apis/v1/... crd:crdVersions=v1 output:crd:dir=./crds
func main() {
	cmd.Execute()
}
