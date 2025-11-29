package services

import (
	"fmt"

	"github.com/convin/crae/internal/models"
	"github.com/jmoiron/sqlx"
)

type PermissionService struct {
	db *sqlx.DB
}

func NewPermissionService(db *sqlx.DB) *PermissionService {
	return &PermissionService{db: db}
}

// CheckPermission checks if a user has a specific permission
// Returns true if user has permission, false otherwise
func (s *PermissionService) CheckPermission(userID int64, permissionCodeName string) (bool, error) {
	// Get user details
	var user models.User
	err := s.db.Get(&user,
		`SELECT id, tenant_id, user_type, role_id, is_active
		FROM users WHERE id = $1`,
		userID,
	)
	if err != nil {
		return false, fmt.Errorf("user not found: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return false, nil
	}

	// Observers have no functional access regardless of role
	if user.UserType == "observer" {
		return false, nil
	}

	// Product users need a role to have permissions
	if user.UserType != "product_user" {
		return false, nil
	}

	// If no role assigned, no permissions
	if user.RoleID == nil {
		return false, nil
	}

	// Check if role has the permission
	var hasPermission bool
	err = s.db.Get(&hasPermission,
		`SELECT $1 = ANY(code_names)
		FROM roles
		WHERE id = $2 AND tenant_id = $3`,
		permissionCodeName, *user.RoleID, user.TenantID,
	)
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}

	return hasPermission, nil
}

// GetUserPermissions returns all permissions for a user
func (s *PermissionService) GetUserPermissions(userID int64) ([]string, error) {
	// Get user details
	var user models.User
	err := s.db.Get(&user,
		`SELECT id, tenant_id, user_type, role_id, is_active
		FROM users WHERE id = $1`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return []string{}, nil
	}

	// Observers have no permissions
	if user.UserType == "observer" {
		return []string{}, nil
	}

	// Product users need a role to have permissions
	if user.UserType != "product_user" {
		return []string{}, nil
	}

	// If no role assigned, no permissions
	if user.RoleID == nil {
		return []string{}, nil
	}

	// Get permissions from role
	var codeNames []string
	err = s.db.Get(&codeNames,
		`SELECT code_names
		FROM roles
		WHERE id = $1 AND tenant_id = $2`,
		*user.RoleID, user.TenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	return codeNames, nil
}

// HasAnyPermission checks if user has any of the given permissions
func (s *PermissionService) HasAnyPermission(userID int64, permissionCodeNames []string) (bool, error) {
	permissions, err := s.GetUserPermissions(userID)
	if err != nil {
		return false, err
	}

	// Create a map for quick lookup
	permMap := make(map[string]bool)
	for _, perm := range permissions {
		permMap[perm] = true
	}

	// Check if user has any of the required permissions
	for _, requiredPerm := range permissionCodeNames {
		if permMap[requiredPerm] {
			return true, nil
		}
	}

	return false, nil
}

