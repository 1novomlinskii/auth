package httpdelivery

import (
	"net/http"

	"github.com/user/auth/internal/usecase"
)

// AuthHandler handles authentication-related HTTP endpoints.
type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

// Register handles POST /v1/auth/register.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}

// Login handles POST /v1/auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}

// RefreshToken handles POST /v1/auth/refresh.
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}

// Logout handles POST /v1/auth/logout.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}

// GetOAuth2URL handles GET /v1/auth/oauth2/{provider}/url.
func (h *AuthHandler) GetOAuth2URL(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}

// LoginOAuth2 handles POST /v1/auth/oauth2/login.
func (h *AuthHandler) LoginOAuth2(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}
