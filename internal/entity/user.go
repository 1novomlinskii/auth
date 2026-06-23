// Package entity contains domain types for the auth service.
package entity

// User represents a user in the system.
type User struct {
	UserID       string
	TenantID     string
	Login        string
	DisplayName  string
	EmailPrimary string
	Locale       string
	Timezone     string
	IsActive     bool
	CreatedAt    string
	UpdatedAt    *string
}

// Session represents an active user session (refresh token).
type Session struct {
	SessionID        string
	UserID           string
	IdentityID       *string
	RefreshTokenHash string
	DeviceInfo       *string
	IPAddress        *string
	ExpiresAt        string
	CreatedAt        string
	RevokedAt        *string
}
