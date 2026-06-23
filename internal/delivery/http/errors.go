package httpdelivery

import (
	"encoding/json"
	"net/http"
)

// Google-style error response.
type apiError struct {
	Error apiErrorItem `json:"error"`
}

type apiErrorItem struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// writeError sends a JSON error response in Google API style.
func writeError(w http.ResponseWriter, httpStatus int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(apiError{
		Error: apiErrorItem{Code: code, Message: message},
	})
}

// writeJSON sends a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

// Common error codes.
const (
	errInvalidArgument  = "INVALID_ARGUMENT"
	errNotFound         = "NOT_FOUND"
	errInternal         = "INTERNAL"
	errUnauthenticated  = "UNAUTHENTICATED"
	errPermissionDenied = "PERMISSION_DENIED"
)

// HTTP status to error code mapping helpers.
func writeInvalidArgument(w http.ResponseWriter, message string) {
	writeError(w, http.StatusBadRequest, errInvalidArgument, message)
}

func writeNotFound(w http.ResponseWriter, message string) {
	writeError(w, http.StatusNotFound, errNotFound, message)
}

func writeInternal(w http.ResponseWriter, message string) {
	writeError(w, http.StatusInternalServerError, errInternal, message)
}

func writeUnauthenticated(w http.ResponseWriter, message string) {
	writeError(w, http.StatusUnauthorized, errUnauthenticated, message)
}

func writePermissionDenied(w http.ResponseWriter, message string) {
	writeError(w, http.StatusForbidden, errPermissionDenied, message)
}
