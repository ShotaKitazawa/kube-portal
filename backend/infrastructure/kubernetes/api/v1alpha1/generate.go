package v1alpha1

//go:generate go run sigs.k8s.io/controller-tools/cmd/controller-gen object:headerFile="boilerplate.go.txt" paths="."
//go:generate go run sigs.k8s.io/controller-tools/cmd/controller-gen rbac:roleName=manager-role crd paths="./..." output:crd:artifacts:config=../../../../../manifests/crd
