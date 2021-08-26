package api

import (

	// Import to force an explicit dependency and version in go.mod, allowing the generation script to work.
	_ "github.com/deepmap/oapi-codegen/pkg/codegen"
)

//go:generate ../../scripts/genapi.sh ../../docs/api.yaml types types.go
//go:generate ../../scripts/genapi.sh ../../docs/api.yaml spec spec.go
