package httpdelivery

import (
	"net/http"

	"github.com/user/auth/internal/usecase"
)

// UserHandler handles user management HTTP endpoints.
type UserHandler struct {
	userUsecase *usecase.UserUsecase
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userUsecase *usecase.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

// GetUser handles GET /v1/users/{user_id}.
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}

// UpdateUser handles PATCH /v1/users/{user_id}.
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}

// ListUsers handles GET /v1/users.
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}
