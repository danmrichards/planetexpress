package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danmrichards/planetexpress/internal/api"
	"github.com/danmrichards/planetexpress/internal/api/middleware"
	"github.com/danmrichards/planetexpress/internal/api/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type handler struct {
	evtSvc   EventService
	shipSvc  ShipService
	pkgIDGen pkgIDFunc
}

// Init initialises a new API handler using the given router.
func Init(r *mux.Router, evtSvc EventService, shipSvc ShipService) error {
	mw, err := validateAPIRequestsMiddleware()
	if err != nil {
		return fmt.Errorf("validation middleware: %w", err)
	}

	sr := r.PathPrefix("/v1").Subrouter()
	r.Use(mw)

	h := &handler{
		evtSvc:   evtSvc,
		shipSvc:  shipSvc,
		pkgIDGen: uuid.NewString,
	}

	sr.Path("/ship/status").Methods(http.MethodGet).HandlerFunc(h.shipStatus)
	sr.Path("/package/allocate").Methods(http.MethodPost).HandlerFunc(h.packageAllocate)
	sr.Path("/package/{package_id}/load").Methods(http.MethodPut).HandlerFunc(h.packageLoad)
	sr.Path("/package/{package_id}/unload").Methods(http.MethodPut).HandlerFunc(h.packageUnload)

	return nil
}

func validateAPIRequestsMiddleware() (mux.MiddlewareFunc, error) {
	swag, err := api.GetSwagger()
	if err != nil {
		return nil, err
	}

	m, err := middleware.NewSwaggerValidationMiddlewareWithErrEnc(swag, swaggerErrorJSON)
	if err != nil {
		return nil, err
	}

	return m.Middleware, nil
}

func swaggerErrorJSON(ctx context.Context, err error, w http.ResponseWriter) {
	response.WriteError(w, http.StatusBadRequest, "Validation error", err.Error())
}
