//go:build tools

package hack

import (
	_ "github.com/ogen-go/ogen/cmd/ogen"
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen"
)
