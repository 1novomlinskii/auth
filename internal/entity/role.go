package entity

// Role represents a role with associated permissions.
type Role struct {
	RoleID      string
	TenantID    string
	Name        string
	Permissions []string
	IsSystem    bool
	Description string
	CreatedAt   string
	UpdatedAt   *string
}
