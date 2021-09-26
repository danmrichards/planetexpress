package handler

import (
	"encoding/json"
	"net/http"

	"github.com/danmrichards/planetexpress/internal/api"
	"github.com/danmrichards/planetexpress/internal/api/response"
)

func (h handler) packageAllocate(w http.ResponseWriter, r *http.Request) {
	// Note that no other validation is being done here, as it is handled by the
	// swagger middleware.
	var allocateRequest api.Package
	if err := json.NewDecoder(r.Body).Decode(&allocateRequest); err != nil {
		response.WriteDefaultStatusError(w, http.StatusBadRequest)
		return
	}

	// TODO: Check current state of ship to see if can be allocated.
	// State of ship is equal to the current snapshot of state plus any newer
	// events.

	packageID := h.pkgIDGen()

	// Dispatch allocate event.
	if err := h.evtSvc.PackageAllocate(r.Context(), packageID, allocateRequest.Size); err != nil {
		response.WriteError(
			w,
			http.StatusInternalServerError,
			"could not allocate package",
			err.Error(),
		)
		return
	}

	if err := response.WriteStatusJSON(
		w,
		http.StatusCreated,
		api.AllocatePackage{
			Data: api.AllocatedPackage{
				Package:   allocateRequest,
				PackageId: packageID,
			},
		},
	); err != nil {
		response.WriteError(
			w,
			http.StatusInternalServerError,
			"could not write response",
			err.Error(),
		)
		return
	}
}
