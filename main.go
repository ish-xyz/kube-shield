package main

import "github.com/RedLabsPlatform/kube-shield/cmd"

//go:generate go run sigs.k8s.io/controller-tools/cmd/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./pkg/apis/..."
//go:generate go run sigs.k8s.io/controller-tools/cmd/controller-gen crd paths=./pkg/apis/... crd:crdVersions=v1 output:crd:dir=./crd-gen

func main() {
	cmd.Execute()
}
