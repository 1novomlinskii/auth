package httpdelivery

import (
	"context"
	"net/http"
)

type contextKey string

const (
	// ContextUserID is the key for storing the authenticated user ID in the request context.
	ContextUserID contextKey = "user_id"
	// ContextTenantID is the key for storing the tenant ID in the request context.
	ContextTenantID contextKey = "tenant_id"
)

// CORS is a simple CORS middleware that allows all origins.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AuthRequired is a stub middleware that extracts user info from the token.
// Currently it passes through without verification.
func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: parse JWT from Authorization header and set ContextUserID / ContextTenantID
		ctx := context.WithValue(r.Context(), ContextUserID, "unknown")
		ctx = context.WithValue(ctx, ContextTenantID, "default")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
