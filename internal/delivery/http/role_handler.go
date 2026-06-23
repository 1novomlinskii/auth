package httpdelivery

import (
	"net/http"

	"github.com/user/auth/internal/usecase"
)

// RoleHandler handles role management HTTP endpoints.
type RoleHandler struct {
	roleUsecase *usecase.RoleUsecase
}

// NewRoleHandler creates a new RoleHandler.
func NewRoleHandler(roleUsecase *usecase.RoleUsecase) *RoleHandler {
	return &RoleHandler{roleUsecase: roleUsecase}
}

// CreateRole handles POST /v1/roles.
func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}

// GetRole handles GET /v1/roles/{role_id}.
func (h *RoleHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}

// ListRoles handles GET /v1/roles.
func (h *RoleHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}

// UpdateRole handles PATCH /v1/roles/{role_id}.
func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}

// DeleteRole handles DELETE /v1/roles/{role_id}.
func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}

// AssignRole handles POST /v1/roles/{role_id}/assign.
func (h *RoleHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}

// RevokeRole handles POST /v1/roles/{role_id}/revoke.
func (h *RoleHandler) RevokeRole(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}

// CheckPermission handles POST /v1/permissions/check.
func (h *RoleHandler) CheckPermission(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}

// GetUserRoles handles GET /v1/users/{user_id}/roles.
func (h *RoleHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, errInvalidArgument, "not implemented")
}
