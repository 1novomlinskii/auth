package entity

// APIKey represents an API key for service-to-service authentication.
type APIKey struct {
	APIKeyID    string
	TenantID    string
	ServiceName string
	KeyPrefix   string
	KeyHash     string
	Permissions []string
	CreatedBy   string
	ExpiresAt   *string
	IsActive    bool
	CreatedAt   string
}
