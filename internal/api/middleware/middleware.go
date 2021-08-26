package middleware

import (
	"context"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/gorilla/mux"
)

// SwaggerValidationMiddleware is a gorilla/mux middleware that validates requests against an OpenAPI spec.
type SwaggerValidationMiddleware struct {
	ErrorEncoder       openapi3filter.ErrorEncoder
	AuthenticationFunc openapi3filter.AuthenticationFunc
	router             routers.Router
}

// NewSwaggerValidationMiddleware returns a new SwaggerValidationMiddleware
func NewSwaggerValidationMiddleware(swag *openapi3.T) (*SwaggerValidationMiddleware, error) {
	return NewSwaggerValidationMiddlewareWithErrEnc(swag, openapi3filter.DefaultErrorEncoder)
}

// NewSwaggerValidationMiddlewareWithErrEnc returns a new SwaggerValidationMiddleware
func NewSwaggerValidationMiddlewareWithErrEnc(swag *openapi3.T, errenc openapi3filter.ErrorEncoder) (*SwaggerValidationMiddleware, error) {
	if err := swag.Validate(context.Background()); err != nil {
		return nil, err
	}

	router, err := gorillamux.NewRouter(swag)
	if err != nil {
		return nil, err
	}

	return &SwaggerValidationMiddleware{
		ErrorEncoder:       errenc,
		AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
		router:             router,
	}, nil
}

// Middleware implements gorilla/mux MiddlewareFunc
var _ mux.MiddlewareFunc = (&SwaggerValidationMiddleware{}).Middleware

// Middleware implements mux.Middleware.
func (m *SwaggerValidationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handled := m.before(w, r); handled {
			return
		}
		// TODO: validateResponse
		next.ServeHTTP(w, r)
	})
}

func (m *SwaggerValidationMiddleware) before(w http.ResponseWriter, r *http.Request) (handled bool) {
	if err := m.validateRequest(r); err != nil {
		m.ErrorEncoder(r.Context(), err, w)
		return true
	}
	return false
}

func (m *SwaggerValidationMiddleware) validateRequest(r *http.Request) error {
	// Find route
	route, pathParams, err := m.router.FindRoute(r)
	if err != nil {
		return err
	}

	options := &openapi3filter.Options{
		AuthenticationFunc: m.AuthenticationFunc,
	}

	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    r,
		PathParams: pathParams,
		Route:      route,
		Options:    options,
	}
	if err = openapi3filter.ValidateRequest(r.Context(), requestValidationInput); err != nil {
		return err
	}

	return nil
}
