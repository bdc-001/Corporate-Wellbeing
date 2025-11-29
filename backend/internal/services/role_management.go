package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/convin/crae/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type RoleManagementService struct {
	db *sqlx.DB
}

func NewRoleManagementService(db *sqlx.DB) *RoleManagementService {
	return &RoleManagementService{db: db}
}

// CreateRole creates a new role
func (s *RoleManagementService) CreateRole(tenantID int64, req CreateRoleRequest) (*models.Role, error) {
	// Check if role name already exists
	var existingID int
	err := s.db.Get(&existingID, `SELECT id FROM roles WHERE tenant_id = $1 AND name = $2`, tenantID, req.Name)
	if err == nil {
		return nil, fmt.Errorf("role with name %s already exists", req.Name)
	}

	// Insert role
	var role models.Role
	var codeNamesArr pq.StringArray
	err = s.db.QueryRowx(
		`INSERT INTO roles (tenant_id, name, description, code_names, can_be_edited, is_default)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, tenant_id, name, description, code_names, can_be_edited, is_default, created_at, updated_at`,
		tenantID, req.Name, req.Description, pq.Array(req.CodeNames), true, false,
	).Scan(
		&role.ID, &role.TenantID, &role.Name, &role.Description, &codeNamesArr,
		&role.CanBeEdited, &role.IsDefault, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}
	role.CodeNames = []string(codeNamesArr)

	// Insert team restrictions if provided
	if len(req.AllowedTeamIDs) > 0 {
		for _, teamID := range req.AllowedTeamIDs {
			_, err = s.db.Exec(
				`INSERT INTO role_teams (role_id, team_id) VALUES ($1, $2)
				ON CONFLICT (role_id, team_id) DO NOTHING`,
				role.ID, teamID,
			)
			if err != nil {
				// Log error but continue
				fmt.Printf("Warning: failed to add team restriction: %v\n", err)
			}
		}
	}

	// Fetch team IDs
	role.AllowedTeamIDs = req.AllowedTeamIDs

	return &role, nil
}

// UpdateRole updates an existing role
func (s *RoleManagementService) UpdateRole(tenantID int64, roleID int, req UpdateRoleRequest) (*models.Role, error) {
	// Check if role exists and can be edited
	var canBeEdited bool
	var isDefault bool
	err := s.db.Get(&canBeEdited, `SELECT can_be_edited FROM roles WHERE tenant_id = $1 AND id = $2`, tenantID, roleID)
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}
	if !canBeEdited {
		return nil, fmt.Errorf("role cannot be edited")
	}

	err = s.db.Get(&isDefault, `SELECT is_default FROM roles WHERE tenant_id = $1 AND id = $2`, tenantID, roleID)
	if err == nil && isDefault && req.CodeNames != nil {
		// Prevent changing permissions on default roles
		return nil, fmt.Errorf("cannot modify permissions of default role")
	}

	// Build update query
	updates := []string{}
	args := []interface{}{}
	argPos := 1

	if req.Name != nil {
		// Check if new name conflicts
		var existingID int
		err = s.db.Get(&existingID, `SELECT id FROM roles WHERE tenant_id = $1 AND name = $2 AND id != $3`,
			tenantID, *req.Name, roleID)
		if err == nil {
			return nil, fmt.Errorf("role with name %s already exists", *req.Name)
		}
		updates = append(updates, fmt.Sprintf("name = $%d", argPos))
		args = append(args, *req.Name)
		argPos++
	}
	if req.Description != nil {
		updates = append(updates, fmt.Sprintf("description = $%d", argPos))
		args = append(args, *req.Description)
		argPos++
	}
	if req.CodeNames != nil {
		updates = append(updates, fmt.Sprintf("code_names = $%d", argPos))
		args = append(args, pq.Array(req.CodeNames))
		argPos++
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	updates = append(updates, fmt.Sprintf("updated_at = $%d", argPos))
	args = append(args, time.Now())
	argPos++

	// Add WHERE clause
	args = append(args, tenantID, roleID)
	wherePos1 := argPos
	wherePos2 := argPos + 1

	query := fmt.Sprintf(
		`UPDATE roles SET %s WHERE tenant_id = $%d AND id = $%d
		RETURNING id, tenant_id, name, description, code_names, can_be_edited, is_default, created_at, updated_at`,
		strings.Join(updates, ", "), wherePos1, wherePos2,
	)

	var role models.Role
	var codeNamesArr pq.StringArray
	err = s.db.QueryRowx(query, args...).Scan(
		&role.ID, &role.TenantID, &role.Name, &role.Description, &codeNamesArr,
		&role.CanBeEdited, &role.IsDefault, &role.CreatedAt, &role.UpdatedAt,
	)
	if err == nil {
		role.CodeNames = []string(codeNamesArr)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	// Update team restrictions if provided
	if req.AllowedTeamIDs != nil {
		// Delete existing restrictions
		_, err = s.db.Exec(`DELETE FROM role_teams WHERE role_id = $1`, roleID)
		if err != nil {
			return nil, fmt.Errorf("failed to clear team restrictions: %w", err)
		}

		// Insert new restrictions
		for _, teamID := range req.AllowedTeamIDs {
			_, err = s.db.Exec(
				`INSERT INTO role_teams (role_id, team_id) VALUES ($1, $2)`,
				roleID, teamID,
			)
			if err != nil {
				// Log error but continue
				fmt.Printf("Warning: failed to add team restriction: %v\n", err)
			}
		}
		role.AllowedTeamIDs = req.AllowedTeamIDs
	} else {
		// Fetch existing team IDs
		var teamIDs []int
		err = s.db.Select(&teamIDs, `SELECT team_id FROM role_teams WHERE role_id = $1`, roleID)
		if err == nil {
			role.AllowedTeamIDs = teamIDs
		}
	}

	return &role, nil
}

// GetRole retrieves a role by ID
func (s *RoleManagementService) GetRole(tenantID int64, roleID int) (*models.Role, error) {
	var role models.Role
	var codeNamesArr pq.StringArray
	err := s.db.QueryRowx(
		`SELECT id, tenant_id, name, description, code_names, can_be_edited, is_default, created_at, updated_at
		FROM roles WHERE tenant_id = $1 AND id = $2`,
		tenantID, roleID,
	).Scan(
		&role.ID, &role.TenantID, &role.Name, &role.Description, &codeNamesArr,
		&role.CanBeEdited, &role.IsDefault, &role.CreatedAt, &role.UpdatedAt,
	)
	if err == nil {
		role.CodeNames = []string(codeNamesArr)
	}
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}

	// Get team restrictions
	var teamIDs []int
	err = s.db.Select(&teamIDs, `SELECT team_id FROM role_teams WHERE role_id = $1`, roleID)
	if err == nil {
		role.AllowedTeamIDs = teamIDs
	}

	// Get user count
	var userCount int
	err = s.db.Get(&userCount, `SELECT COUNT(*) FROM users WHERE role_id = $1 AND is_active = true`, roleID)
	if err == nil {
		role.UserCount = &userCount
	}

	return &role, nil
}

// ListRoles retrieves all roles for a tenant
func (s *RoleManagementService) ListRoles(tenantID int64) ([]models.Role, error) {
	type RoleRow struct {
		ID          int            `db:"id"`
		TenantID    int64          `db:"tenant_id"`
		Name        string         `db:"name"`
		Description *string       `db:"description"`
		CodeNames   pq.StringArray `db:"code_names"`
		CanBeEdited bool           `db:"can_be_edited"`
		IsDefault   bool           `db:"is_default"`
		CreatedAt   time.Time      `db:"created_at"`
		UpdatedAt   time.Time      `db:"updated_at"`
	}

	var roleRows []RoleRow
	err := s.db.Select(&roleRows,
		`SELECT id, tenant_id, name, description, code_names, can_be_edited, is_default, created_at, updated_at
		FROM roles WHERE tenant_id = $1 ORDER BY is_default DESC, name ASC`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}

	roles := make([]models.Role, len(roleRows))
	for i, row := range roleRows {
		roles[i] = models.Role{
			ID:          row.ID,
			TenantID:    row.TenantID,
			Name:        row.Name,
			Description: row.Description,
			CodeNames:   []string(row.CodeNames),
			CanBeEdited: row.CanBeEdited,
			IsDefault:   row.IsDefault,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
		}
	}

	// Get team restrictions and user counts for each role
	for i := range roles {
		var teamIDs []int
		err = s.db.Select(&teamIDs, `SELECT team_id FROM role_teams WHERE role_id = $1`, roles[i].ID)
		if err == nil {
			roles[i].AllowedTeamIDs = teamIDs
		}

		var userCount int
		err = s.db.Get(&userCount, `SELECT COUNT(*) FROM users WHERE role_id = $1 AND is_active = true`, roles[i].ID)
		if err == nil {
			roles[i].UserCount = &userCount
		}
	}

	return roles, nil
}

// DeleteRole deletes a role (only if it can be edited and is not default)
func (s *RoleManagementService) DeleteRole(tenantID int64, roleID int) error {
	// Check if role can be deleted
	var canBeEdited bool
	var isDefault bool
	err := s.db.Get(&canBeEdited, `SELECT can_be_edited FROM roles WHERE tenant_id = $1 AND id = $2`, tenantID, roleID)
	if err != nil {
		return fmt.Errorf("role not found: %w", err)
	}
	if !canBeEdited {
		return fmt.Errorf("role cannot be deleted")
	}

	err = s.db.Get(&isDefault, `SELECT is_default FROM roles WHERE tenant_id = $1 AND id = $2`, tenantID, roleID)
	if err == nil && isDefault {
		return fmt.Errorf("default role cannot be deleted")
	}

	// Check if role is in use
	var userCount int
	err = s.db.Get(&userCount, `SELECT COUNT(*) FROM users WHERE role_id = $1`, roleID)
	if err == nil && userCount > 0 {
		return fmt.Errorf("cannot delete role: %d users are assigned to this role", userCount)
	}

	// Delete team restrictions
	_, err = s.db.Exec(`DELETE FROM role_teams WHERE role_id = $1`, roleID)
	if err != nil {
		return fmt.Errorf("failed to delete team restrictions: %w", err)
	}

	// Delete role
	_, err = s.db.Exec(`DELETE FROM roles WHERE tenant_id = $1 AND id = $2`, tenantID, roleID)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}

// ListPermissions retrieves all available permissions grouped by permission groups
func (s *RoleManagementService) ListPermissions() ([]models.PermissionGroup, error) {
	// Get all permission groups
	var groups []models.PermissionGroup
	err := s.db.Select(&groups,
		`SELECT id, name, description, display_order, created_at
		FROM permission_groups ORDER BY display_order ASC, name ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list permission groups: %w", err)
	}

	// Get permissions for each group
	for i := range groups {
		var permissions []models.Permission
		err = s.db.Select(&permissions,
			`SELECT id, code_name, name, description, group_id, created_at
			FROM permissions WHERE group_id = $1 ORDER BY name ASC`,
			groups[i].ID,
		)
		if err == nil {
			groups[i].Permissions = permissions
		}
	}

	return groups, nil
}

// Request/Response types
type CreateRoleRequest struct {
	Name           string   `json:"name" binding:"required"`
	Description    *string  `json:"description"`
	CodeNames      []string `json:"code_names"` // Array of permission code names
	AllowedTeamIDs []int    `json:"allowed_team_ids"` // Team restrictions
}

type UpdateRoleRequest struct {
	Name           *string  `json:"name"`
	Description    *string  `json:"description"`
	CodeNames      []string `json:"code_names"`
	AllowedTeamIDs []int    `json:"allowed_team_ids"`
}

