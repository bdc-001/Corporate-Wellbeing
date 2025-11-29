package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/convin/crae/internal/models"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserManagementService struct {
	db *sqlx.DB
}

func NewUserManagementService(db *sqlx.DB) *UserManagementService {
	return &UserManagementService{db: db}
}

// GetDB returns the database connection (for use in handlers)
func (s *UserManagementService) GetDB() *sqlx.DB {
	return s.db
}

// CreateUser creates a new user with auto-generated password
func (s *UserManagementService) CreateUser(tenantID int64, req CreateUserRequest) (*models.User, string, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, "", err
	}
	
	// Check if email already exists
	var existingID int64
	err := s.db.Get(&existingID, `SELECT id FROM users WHERE tenant_id = $1 AND email = $2`, tenantID, req.Email)
	if err == nil {
		return nil, "", fmt.Errorf("user with email %s already exists", req.Email)
	}

	// If user_type is observer, clear role_id
	if req.UserType == "observer" {
		req.RoleID = nil
	}

	// Generate password
	password := req.Password
	if password == "" {
		password = s.generatePassword()
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Construct full name from first_name and last_name for backward compatibility
	fullName := req.FirstName
	if req.LastName != nil && *req.LastName != "" {
		fullName = req.FirstName + " " + *req.LastName
	}
	if req.Name != nil && *req.Name != "" {
		fullName = *req.Name // Use provided name if available (backward compatibility)
	}

	// Insert user
	var user models.User
	var lastName *string
	if req.LastName != nil && *req.LastName != "" {
		lastName = req.LastName
	}
	err = s.db.QueryRowx(
		`INSERT INTO users (
			tenant_id, email, name, first_name, last_name, password_hash, phone, role_id,
			manager_id, auditor_id, team_id, user_type, location, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, tenant_id, email, name, first_name, last_name, phone, role_id, manager_id, auditor_id,
			team_id, user_type, location, is_active, created_at, updated_at`,
		tenantID, req.Email, fullName, req.FirstName, lastName, string(hashedPassword), req.Phone, req.RoleID,
		req.ManagerID, req.AuditorID, req.TeamID, getOrDefault(req.UserType, "product_user"), req.Location, true,
	).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.Name, &user.FirstName, &user.LastName, &user.Phone, &user.RoleID,
		&user.ManagerID, &user.AuditorID, &user.TeamID, &user.UserType, &user.Location,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	return &user, password, nil
}

// BulkCreateUsers creates multiple users from CSV/Excel data
func (s *UserManagementService) BulkCreateUsers(tenantID int64, users []CreateUserRequest) ([]models.User, []string, []error) {
	var createdUsers []models.User
	var passwords []string
	var errors []error

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, nil, []error{fmt.Errorf("failed to begin transaction: %w", err)}
	}
	defer tx.Rollback()

	for i, req := range users {
		// If user_type is observer, clear role_id
		if req.UserType == "observer" {
			req.RoleID = nil
		}

		// Generate password
		password := s.generatePassword()

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			errors = append(errors, fmt.Errorf("row %d: failed to hash password: %w", i+1, err))
			continue
		}

		// Check if email already exists
		var existingID int64
		err = tx.Get(&existingID, `SELECT id FROM users WHERE tenant_id = $1 AND email = $2`, tenantID, req.Email)
		if err == nil {
			// User exists, update instead
			_, err = tx.Exec(
				`UPDATE users SET
					name = $1, phone = $2, role_id = $3, manager_id = $4,
					auditor_id = $5, team_id = $6, user_type = $7, location = $8,
					updated_at = CURRENT_TIMESTAMP
				WHERE tenant_id = $9 AND email = $10`,
				req.Name, req.Phone, req.RoleID, req.ManagerID, req.AuditorID,
				req.TeamID, getOrDefault(req.UserType, "product_user"), req.Location, tenantID, req.Email,
			)
			if err != nil {
				errors = append(errors, fmt.Errorf("row %d: failed to update user: %w", i+1, err))
				continue
			}
			// Get updated user
			var user models.User
			err = tx.Get(&user, `SELECT * FROM users WHERE tenant_id = $1 AND email = $2`, tenantID, req.Email)
			if err == nil {
				createdUsers = append(createdUsers, user)
				passwords = append(passwords, "") // No password for updated users
			}
			continue
		}

		// Insert new user
		var user models.User
		err = tx.QueryRowx(
			`INSERT INTO users (
				tenant_id, email, name, password_hash, phone, role_id,
				manager_id, auditor_id, team_id, user_type, location, is_active
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
			RETURNING id, tenant_id, email, name, phone, role_id, manager_id, auditor_id,
				team_id, user_type, location, is_active, created_at, updated_at`,
			tenantID, req.Email, req.Name, string(hashedPassword), req.Phone, req.RoleID,
			req.ManagerID, req.AuditorID, req.TeamID, getOrDefault(req.UserType, "product_user"), req.Location, true,
		).Scan(
			&user.ID, &user.TenantID, &user.Email, &user.Name, &user.Phone, &user.RoleID,
			&user.ManagerID, &user.AuditorID, &user.TeamID, &user.UserType, &user.Location,
			&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			errors = append(errors, fmt.Errorf("row %d: failed to create user: %w", i+1, err))
			continue
		}

		createdUsers = append(createdUsers, user)
		passwords = append(passwords, password)
	}

	if err = tx.Commit(); err != nil {
		return nil, nil, []error{fmt.Errorf("failed to commit transaction: %w", err)}
	}

	return createdUsers, passwords, errors
}

// UpdateUser updates an existing user
func (s *UserManagementService) UpdateUser(tenantID int64, userID int64, req UpdateUserRequest) (*models.User, error) {
	// If user_type is being set to observer, clear role_id
	if req.UserType != nil && *req.UserType == "observer" {
		req.RoleID = nil
	}

	// Build update query dynamically
	updates := []string{}
	args := []interface{}{}
	argPos := 1

	// Handle name updates: if first_name or last_name is provided, update both name and first_name/last_name
	if req.FirstName != nil || req.LastName != nil {
		// Get current user to construct full name
		var currentUser models.User
		err := s.db.Get(&currentUser, `SELECT first_name, last_name, name FROM users WHERE tenant_id = $1 AND id = $2`, tenantID, userID)
		if err == nil {
			firstName := currentUser.FirstName
			lastName := currentUser.LastName
			if req.FirstName != nil {
				firstName = req.FirstName
			}
			if req.LastName != nil {
				lastName = req.LastName
			}
			
			// Construct full name
			fullName := ""
			if firstName != nil && *firstName != "" {
				fullName = *firstName
				if lastName != nil && *lastName != "" {
					fullName = *firstName + " " + *lastName
				}
			}
			
			updates = append(updates, fmt.Sprintf("name = $%d", argPos))
			args = append(args, fullName)
			argPos++
			
			updates = append(updates, fmt.Sprintf("first_name = $%d", argPos))
			args = append(args, firstName)
			argPos++
			
			updates = append(updates, fmt.Sprintf("last_name = $%d", argPos))
			args = append(args, lastName)
			argPos++
		}
	} else if req.Name != nil {
		// If only name is provided (backward compatibility), update name only
		updates = append(updates, fmt.Sprintf("name = $%d", argPos))
		args = append(args, *req.Name)
		argPos++
	}
	if req.Phone != nil {
		updates = append(updates, fmt.Sprintf("phone = $%d", argPos))
		args = append(args, *req.Phone)
		argPos++
	}
	if req.RoleID != nil {
		updates = append(updates, fmt.Sprintf("role_id = $%d", argPos))
		args = append(args, *req.RoleID)
		argPos++
	}
	if req.ManagerID != nil {
		updates = append(updates, fmt.Sprintf("manager_id = $%d", argPos))
		args = append(args, *req.ManagerID)
		argPos++
	}
	if req.AuditorID != nil {
		updates = append(updates, fmt.Sprintf("auditor_id = $%d", argPos))
		args = append(args, *req.AuditorID)
		argPos++
	}
	if req.TeamID != nil {
		updates = append(updates, fmt.Sprintf("team_id = $%d", argPos))
		args = append(args, *req.TeamID)
		argPos++
	}
	if req.UserType != nil {
		updates = append(updates, fmt.Sprintf("user_type = $%d", argPos))
		args = append(args, *req.UserType)
		argPos++
	}
	if req.Location != nil {
		updates = append(updates, fmt.Sprintf("location = $%d", argPos))
		args = append(args, *req.Location)
		argPos++
	}
	if req.Timezone != nil {
		updates = append(updates, fmt.Sprintf("timezone = $%d", argPos))
		args = append(args, *req.Timezone)
		argPos++
	}
	if req.IsActive != nil {
		updates = append(updates, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, *req.IsActive)
		argPos++
	}
	if req.Password != nil && *req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		updates = append(updates, fmt.Sprintf("password_hash = $%d", argPos))
		args = append(args, string(hashedPassword))
		argPos++
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// Always update the updated_at timestamp (no placeholder needed)
	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")

	// Add WHERE clause arguments
	// argPos is the next available position, so tenant_id = $argPos, id = $argPos+1
	args = append(args, tenantID, userID)

	query := fmt.Sprintf(
		`UPDATE users SET %s WHERE tenant_id = $%d AND id = $%d
		RETURNING id, tenant_id, email, name, first_name, last_name, phone, role_id, manager_id, auditor_id,
			team_id, user_type, location, timezone, is_active, created_at, updated_at`,
		strings.Join(updates, ", "), argPos, argPos+1,
	)

	var user models.User
	err := s.db.QueryRowx(query, args...).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.Name, &user.FirstName, &user.LastName, &user.Phone, &user.RoleID,
		&user.ManagerID, &user.AuditorID, &user.TeamID, &user.UserType, &user.Location,
		&user.Timezone, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}

// GetUser retrieves a user by ID
func (s *UserManagementService) GetUser(tenantID int64, userID int64) (*models.User, error) {
	var user models.User
	err := s.db.Get(&user,
		`SELECT 
			u.id, u.tenant_id, u.email, u.name, u.first_name, u.last_name, u.phone, u.role_id, 
			u.manager_id, u.auditor_id, u.team_id, u.user_type, u.location, u.timezone,
			u.is_active, u.last_login_at, u.created_at, u.updated_at,
			r.name as role_name, 
			m.name as manager_name,
			a.name as auditor_name, 
			t.name as team_name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		LEFT JOIN users m ON u.manager_id = m.id
		LEFT JOIN users a ON u.auditor_id = a.id
		LEFT JOIN teams t ON u.team_id = t.id
		WHERE u.tenant_id = $1 AND u.id = $2`,
		tenantID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

// ListUsers retrieves all users for a tenant
func (s *UserManagementService) ListUsers(tenantID int64, filters UserFilters) ([]models.User, error) {
	query := `
		SELECT 
			u.id, u.tenant_id, u.email, u.name, u.first_name, u.last_name, u.phone, u.role_id, 
			u.manager_id, u.auditor_id, u.team_id, u.user_type, u.location, u.timezone,
			u.is_active, u.last_login_at, u.created_at, u.updated_at,
			r.name as role_name, 
			m.name as manager_name,
			a.name as auditor_name, 
			t.name as team_name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		LEFT JOIN users m ON u.manager_id = m.id
		LEFT JOIN users a ON u.auditor_id = a.id
		LEFT JOIN teams t ON u.team_id = t.id
		WHERE u.tenant_id = $1
	`
	args := []interface{}{tenantID}
	argPos := 2

	if filters.RoleID != nil {
		query += fmt.Sprintf(" AND u.role_id = $%d", argPos)
		args = append(args, *filters.RoleID)
		argPos++
	}
	if filters.TeamID != nil {
		query += fmt.Sprintf(" AND u.team_id = $%d", argPos)
		args = append(args, *filters.TeamID)
		argPos++
	}
	if filters.IsActive != nil {
		query += fmt.Sprintf(" AND u.is_active = $%d", argPos)
		args = append(args, *filters.IsActive)
		argPos++
	}
	if filters.Search != nil && *filters.Search != "" {
		query += fmt.Sprintf(" AND (u.name ILIKE $%d OR u.email ILIKE $%d)", argPos, argPos)
		searchTerm := "%" + *filters.Search + "%"
		args = append(args, searchTerm)
		argPos++
	}

	query += " ORDER BY u.created_at DESC"

	type UserWithJoins struct {
		models.User
		RoleName    *string `db:"role_name"`
		ManagerName *string `db:"manager_name"`
		AuditorName *string `db:"auditor_name"`
		TeamName    *string `db:"team_name"`
	}

	var usersWithJoins []UserWithJoins
	err := s.db.Select(&usersWithJoins, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Convert to models.User
	users := make([]models.User, len(usersWithJoins))
	for i, u := range usersWithJoins {
		users[i] = u.User
		users[i].RoleName = u.RoleName
		users[i].ManagerName = u.ManagerName
		users[i].AuditorName = u.AuditorName
		users[i].TeamName = u.TeamName
	}

	return users, nil
}

// DeleteUser deletes a user (soft delete by setting is_active = false)
func (s *UserManagementService) DeleteUser(tenantID int64, userID int64) error {
	_, err := s.db.Exec(
		`UPDATE users SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE tenant_id = $1 AND id = $2`,
		tenantID, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// generatePassword generates a random password
func (s *UserManagementService) generatePassword() string {
	b := make([]byte, 12)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:16] // 16 character password
}

// getOrDefault returns the value if not empty, otherwise returns the default
func getOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// Request/Response types
type CreateUserRequest struct {
	Email     string  `json:"email" binding:"required,email"`
	Name      *string `json:"name"` // Kept for backward compatibility, will be constructed from first_name + last_name (optional)
	FirstName string  `json:"first_name" binding:"required"`
	LastName  *string `json:"last_name"` // Optional
	Password  string  `json:"password"` // Optional, auto-generated if empty
	Phone     *string `json:"phone"`
	RoleID    *int    `json:"role_id"`
	ManagerID *int64  `json:"manager_id"`
	AuditorID *int64  `json:"auditor_id"`
	TeamID    *int    `json:"team_id"`
	UserType  string  `json:"user_type"` // product_user, observer (default: product_user)
	Location  *string `json:"location"`
}

// Validate validates the CreateUserRequest
func (r *CreateUserRequest) Validate() error {
	if r.FirstName == "" && (r.Name == nil || *r.Name == "") {
		return fmt.Errorf("first_name is required")
	}
	return nil
}

type UpdateUserRequest struct {
	Name      *string `json:"name"` // Kept for backward compatibility
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Phone     *string `json:"phone"`
	RoleID    *int    `json:"role_id"`
	ManagerID *int64  `json:"manager_id"`
	AuditorID *int64  `json:"auditor_id"`
	TeamID    *int    `json:"team_id"`
	UserType  *string `json:"user_type"`
	Location  *string `json:"location"`
	Timezone  *string `json:"timezone"`
	IsActive  *bool   `json:"is_active"`
	Password  *string `json:"password"`
}

type UserFilters struct {
	RoleID   *int    `json:"role_id"`
	TeamID   *int    `json:"team_id"`
	IsActive *bool   `json:"is_active"`
	Search   *string `json:"search"`
}

// AuthenticateUser authenticates a user by email and password
func (s *UserManagementService) AuthenticateUser(tenantID int64, email, password string) (*models.User, error) {
	var user models.User
	err := s.db.Get(&user,
		`SELECT 
			u.id, u.tenant_id, u.email, u.name, u.first_name, u.last_name, u.password_hash, u.phone, u.role_id, 
			u.manager_id, u.auditor_id, u.team_id, u.user_type, u.location, u.timezone,
			u.is_active, u.last_login_at, u.created_at, u.updated_at
		FROM users u
		WHERE u.tenant_id = $1 AND u.email = $2 AND u.is_active = true`,
		tenantID, email,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Update last login
	_, err = s.db.Exec(
		`UPDATE users SET last_login_at = CURRENT_TIMESTAMP WHERE id = $1`,
		user.ID,
	)
	if err != nil {
		// Log error but don't fail authentication
		fmt.Printf("Failed to update last_login_at: %v\n", err)
	}

	// Clear password hash from response
	user.PasswordHash = ""

	return &user, nil
}
