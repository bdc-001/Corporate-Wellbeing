package services

import (
	"fmt"
	"strings"

	"github.com/convin/crae/internal/models"
	"github.com/jmoiron/sqlx"
)

type TeamManagementService struct {
	db *sqlx.DB
}

func NewTeamManagementService(db *sqlx.DB) *TeamManagementService {
	return &TeamManagementService{db: db}
}

// GetDB returns the database connection
func (s *TeamManagementService) GetDB() *sqlx.DB {
	return s.db
}

// CreateTeamRequest represents a request to create a team
type CreateTeamRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	VendorID    *int     `json:"vendor_id"` // Optional vendor ID
	GroupID     *int     `json:"group_id"` // Parent team ID
	Members     []int64  `json:"members"` // User IDs to assign to team (optional, not used in UI)
	Subteams    []CreateTeamRequest `json:"subteams"` // Recursive subteams (optional, not used in UI)
}

// UpdateTeamRequest represents a request to update a team
type UpdateTeamRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Members     []int64  `json:"members"` // Optional, not used in UI
	Subteams    []CreateTeamRequest `json:"subteams"` // Optional, not used in UI
}

// CreateTeam creates a new team with optional subteams and members
func (s *TeamManagementService) CreateTeam(tenantID int64, req CreateTeamRequest) (*models.Team, error) {
	// Validation: Can't add members to parent team if it has subteams
	if len(req.Subteams) > 0 && len(req.Members) > 0 {
		return nil, fmt.Errorf("cannot add members to parent team if it has subteams")
	}

	// Start transaction
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Validate vendor_id if provided
	if req.VendorID != nil {
		var vendorExists bool
		err = tx.Get(&vendorExists,
			`SELECT EXISTS(SELECT 1 FROM vendors WHERE id = $1 AND tenant_id = $2)`,
			*req.VendorID, tenantID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to validate vendor: %w", err)
		}
		if !vendorExists {
			return nil, fmt.Errorf("vendor with id %d does not exist for this tenant", *req.VendorID)
		}
	}

	// Create main team
	var team models.Team
	err = tx.QueryRowx(
		`INSERT INTO teams (tenant_id, vendor_id, name, description, group_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, tenant_id, vendor_id, name, description, about, group_id, manager_id, use_case_id, created_at, updated_at`,
		tenantID, req.VendorID, req.Name, req.Description, req.GroupID,
	).Scan(
		&team.ID, &team.TenantID, &team.VendorID, &team.Name, &team.Description, &team.About,
		&team.GroupID, &team.ManagerID, &team.UseCaseID, &team.CreatedAt, &team.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	// Update members (assign to team and set manager)
	if len(req.Members) > 0 {
		err = s.updateMembers(tx, req.Members, team.ID, team.ManagerID)
		if err != nil {
			return nil, err
		}
	}

	// Create subteams recursively
	if len(req.Subteams) > 0 {
		for _, subteamReq := range req.Subteams {
			subteamReq.GroupID = &team.ID
			_, err = s.createTeamInTx(tx, tenantID, subteamReq)
			if err != nil {
				return nil, err
			}
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &team, nil
}

// createTeamInTx creates a team within a transaction
func (s *TeamManagementService) createTeamInTx(tx *sqlx.Tx, tenantID int64, req CreateTeamRequest) (*models.Team, error) {
	// Validate vendor_id if provided
	if req.VendorID != nil {
		var vendorExists bool
		err := tx.Get(&vendorExists,
			`SELECT EXISTS(SELECT 1 FROM vendors WHERE id = $1 AND tenant_id = $2)`,
			*req.VendorID, tenantID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to validate vendor: %w", err)
		}
		if !vendorExists {
			return nil, fmt.Errorf("vendor with id %d does not exist for this tenant", *req.VendorID)
		}
	}

	var team models.Team
	err := tx.QueryRowx(
		`INSERT INTO teams (tenant_id, vendor_id, name, description, group_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, tenant_id, vendor_id, name, description, about, group_id, manager_id, use_case_id, created_at, updated_at`,
		tenantID, req.VendorID, req.Name, req.Description, req.GroupID,
	).Scan(
		&team.ID, &team.TenantID, &team.VendorID, &team.Name, &team.Description, &team.About,
		&team.GroupID, &team.ManagerID, &team.UseCaseID, &team.CreatedAt, &team.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create subteam: %w", err)
	}

	// Update members
	if len(req.Members) > 0 {
		err = s.updateMembers(tx, req.Members, team.ID, team.ManagerID)
		if err != nil {
			return nil, err
		}
	}

	// Create nested subteams
	if len(req.Subteams) > 0 {
		for _, subteamReq := range req.Subteams {
			subteamReq.GroupID = &team.ID
			_, err = s.createTeamInTx(tx, tenantID, subteamReq)
			if err != nil {
				return nil, err
			}
		}
	}

	return &team, nil
}

// updateMembers assigns members to a team and sets their manager
func (s *TeamManagementService) updateMembers(tx *sqlx.Tx, memberIDs []int64, teamID int, managerID *int64) error {
	if len(memberIDs) == 0 {
		return nil
	}

	query := `UPDATE users SET team_id = $1, manager_id = $2, updated_at = CURRENT_TIMESTAMP WHERE id = ANY($3)`
	_, err := tx.Exec(query, teamID, managerID, memberIDs)
	if err != nil {
		return fmt.Errorf("failed to update members: %w", err)
	}

	return nil
}

// GetTeam retrieves a team with subteams and members
func (s *TeamManagementService) GetTeam(tenantID int64, teamID int) (*models.Team, []models.Team, []models.User, error) {
	var team models.Team
	err := s.db.Get(&team,
		`SELECT id, tenant_id, vendor_id, name, description, about, group_id, manager_id, use_case_id, created_at, updated_at
		FROM teams WHERE id = $1 AND tenant_id = $2`,
		teamID, tenantID,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("team not found: %w", err)
	}

	// Get subteams
	var subteams []models.Team
	err = s.db.Select(&subteams,
		`SELECT id, tenant_id, vendor_id, name, description, about, group_id, manager_id, use_case_id, created_at, updated_at
		FROM teams WHERE group_id = $1 AND tenant_id = $2 ORDER BY name ASC`,
		teamID, tenantID,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get subteams: %w", err)
	}

	// Get members
	var members []models.User
	err = s.db.Select(&members,
		`SELECT id, tenant_id, email, name, phone, role_id, manager_id, auditor_id, team_id, user_type, location, is_active, created_at, updated_at
		FROM users WHERE team_id = $1 AND tenant_id = $2 ORDER BY name ASC`,
		teamID, tenantID,
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get members: %w", err)
	}

	return &team, subteams, members, nil
}

// ListTeams retrieves all teams for a tenant (with optional filtering)
func (s *TeamManagementService) ListTeams(tenantID int64, includeSubteams bool) ([]models.Team, error) {
	var teams []models.Team
	query := `SELECT id, tenant_id, vendor_id, name, description, about, group_id, manager_id, use_case_id, created_at, updated_at
		FROM teams WHERE tenant_id = $1`
	
	if !includeSubteams {
		query += " AND group_id IS NULL"
	}
	
	query += " ORDER BY name ASC"

	err := s.db.Select(&teams, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list teams: %w", err)
	}

	return teams, nil
}

// UpdateTeam updates a team
func (s *TeamManagementService) UpdateTeam(tenantID int64, teamID int, req UpdateTeamRequest) (*models.Team, error) {
	// Start transaction
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Build update query
	updates := []string{}
	args := []interface{}{}
	argPos := 1

	if req.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argPos))
		args = append(args, *req.Name)
		argPos++
	}
	if req.Description != nil {
		updates = append(updates, fmt.Sprintf("description = $%d", argPos))
		args = append(args, *req.Description)
		argPos++
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, tenantID, teamID)

	query := fmt.Sprintf(
		`UPDATE teams SET %s WHERE tenant_id = $%d AND id = $%d
		RETURNING id, tenant_id, vendor_id, name, description, about, group_id, manager_id, use_case_id, created_at, updated_at`,
		strings.Join(updates, ", "), argPos, argPos+1,
	)

	var team models.Team
	err = tx.QueryRowx(query, args...).Scan(
		&team.ID, &team.TenantID, &team.VendorID, &team.Name, &team.Description, &team.About,
		&team.GroupID, &team.ManagerID, &team.UseCaseID, &team.CreatedAt, &team.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update team: %w", err)
	}

	// Update members if provided
	if req.Members != nil {
		err = s.updateMembers(tx, req.Members, teamID, team.ManagerID)
		if err != nil {
			return nil, err
		}
	}

	// Update subteams if provided
	if req.Subteams != nil {
		// Delete existing subteams (or handle update logic)
		// For simplicity, we'll just create new ones
		// In production, you'd want to handle updates/deletes more carefully
		for _, subteamReq := range req.Subteams {
			subteamReq.GroupID = &teamID
			_, err = s.createTeamInTx(tx, tenantID, subteamReq)
			if err != nil {
				return nil, err
			}
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &team, nil
}

// DeleteTeam deletes a team with safe member migration
func (s *TeamManagementService) DeleteTeam(tenantID int64, teamID int, transferToTeamID *int) error {
	// Start transaction
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if team has subteams
	var subteamCount int
	err = tx.Get(&subteamCount, `SELECT COUNT(*) FROM teams WHERE group_id = $1 AND tenant_id = $2`, teamID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to check subteams: %w", err)
	}

	if subteamCount > 0 {
		return fmt.Errorf("cannot delete team with subteams. Please delete or migrate subteams first")
	}

	// Check if team has members
	var memberCount int
	err = tx.Get(&memberCount, `SELECT COUNT(*) FROM users WHERE team_id = $1 AND tenant_id = $2`, teamID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to check members: %w", err)
	}

	if memberCount > 0 {
		if transferToTeamID == nil {
			return fmt.Errorf("cannot delete team with members. Please provide transfer_to_team_id")
		}

		// Transfer members to destination team
		_, err = tx.Exec(
			`UPDATE users SET team_id = $1, updated_at = CURRENT_TIMESTAMP WHERE team_id = $2 AND tenant_id = $3`,
			*transferToTeamID, teamID, tenantID,
		)
		if err != nil {
			return fmt.Errorf("failed to transfer members: %w", err)
		}
	}

	// Delete team
	_, err = tx.Exec(`DELETE FROM teams WHERE id = $1 AND tenant_id = $2`, teamID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// AddTeamMembers adds members to a team
func (s *TeamManagementService) AddTeamMembers(tenantID int64, teamID int, memberIDs []int64) error {
	// Get team manager
	var managerID *int64
	err := s.db.Get(&managerID, `SELECT manager_id FROM teams WHERE id = $1 AND tenant_id = $2`, teamID, tenantID)
	if err != nil {
		return fmt.Errorf("team not found: %w", err)
	}

	// Update members
	_, err = s.db.Exec(
		`UPDATE users SET team_id = $1, manager_id = $2, updated_at = CURRENT_TIMESTAMP WHERE id = ANY($3) AND tenant_id = $4`,
		teamID, managerID, memberIDs, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to add members: %w", err)
	}

	return nil
}

// Use Case Management

// CreateUseCase creates a new use case
func (s *TeamManagementService) CreateUseCase(tenantID int64, name, description string) (*models.UseCase, error) {
	var useCase models.UseCase
	err := s.db.QueryRowx(
		`INSERT INTO use_cases (tenant_id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, tenant_id, name, description, created_at, updated_at`,
		tenantID, name, description,
	).Scan(
		&useCase.ID, &useCase.TenantID, &useCase.Name, &useCase.Description, &useCase.CreatedAt, &useCase.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create use case: %w", err)
	}

	return &useCase, nil
}

// ListUseCases lists all use cases for a tenant
func (s *TeamManagementService) ListUseCases(tenantID int64) ([]models.UseCase, error) {
	var useCases []models.UseCase
	err := s.db.Select(&useCases,
		`SELECT id, tenant_id, name, description, created_at, updated_at
		FROM use_cases WHERE tenant_id = $1 ORDER BY name ASC`,
		tenantID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list use cases: %w", err)
	}

	return useCases, nil
}

// DeleteUseCase deletes a use case (only if not in use)
func (s *TeamManagementService) DeleteUseCase(tenantID int64, useCaseID int) error {
	// Check if use case is in use
	var count int
	err := s.db.Get(&count, `SELECT COUNT(*) FROM teams WHERE use_case_id = $1 AND tenant_id = $2`, useCaseID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to check use case usage: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("cannot delete use case that is in use by teams")
	}

	// Delete use case
	_, err = s.db.Exec(`DELETE FROM use_cases WHERE id = $1 AND tenant_id = $2`, useCaseID, tenantID)
	if err != nil {
		return fmt.Errorf("failed to delete use case: %w", err)
	}

	return nil
}

