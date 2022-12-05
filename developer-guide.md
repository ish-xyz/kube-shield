
## Generating new CRDs:

if you made some changes to pkg/apis/*_types.go you most likely want to generate new CRDs. 

To do that you can run:
```
go generate
```
which will create the new CRDs under `./crds/*`.


If the above command doesn't work, you can troubleshoot by running:

```
go run sigs.k8s.io/controller-tools/cmd/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./pkg/apis/..."
go run sigs.k8s.io/controller-tools/cmd/controller-gen crd paths=./pkg/apis/... crd:crdVersions=v1 output:crd:dir=./crd-gen
```
