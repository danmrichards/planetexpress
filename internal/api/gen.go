package api

// Import the oapi-codegen package for side effects:
//
// This forces an explicit dependency and version in go.mod, ensuring that any
// consumers of this template package stays in sync with this third-party
// dependency
import _ "github.com/deepmap/oapi-codegen/pkg/codegen"

//go:generate ../../scripts/genapi.sh ../../docs/api.yaml types types.go
//go:generate ../../scripts/genapi.sh ../../docs/api.yaml spec spec.go
