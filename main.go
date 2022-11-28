package main

import "github.com/RedLabsPlatform/kube-shield/cmd"

//go:generate go run sigs.k8s.io/controller-tools/cmd/controller-gen crd paths=./pkg/loader/... output:dir=./crds

func main() {
	cmd.Execute()
}
