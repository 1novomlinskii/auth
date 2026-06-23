package httpdelivery

import (
	"net/http"

	"github.com/user/auth/internal/usecase"
)

// APIKeyHandler handles API key management HTTP endpoints.
type APIKeyHandler struct {
	apikeyUsecase *usecase.APIKeyUsecase
}

// NewAPIKeyHandler creates a new APIKeyHandler.
func NewAPIKeyHandler(apikeyUsecase *usecase.APIKeyUsecase) *APIKeyHandler {
	return &APIKeyHandler{apikeyUsecase: apikeyUsecase}
}

// CreateAPIKey handles POST /v1/api-keys.
func (h *APIKeyHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}

// ListAPIKeys handles GET /v1/api-keys.
func (h *APIKeyHandler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}

// ValidateAPIKey handles POST /v1/api-keys/validate.
func (h *APIKeyHandler) ValidateAPIKey(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}

// RevokeAPIKey handles POST /v1/api-keys/{api_key_id}/revoke.
func (h *APIKeyHandler) RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}
