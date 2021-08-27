//go:build generate
// +build generate

package handler

import (

	// Import to force a dependency
	_ "github.com/matryer/moq"
)

//go:generate go run github.com/matryer/moq -out handler_test.go . EventService
