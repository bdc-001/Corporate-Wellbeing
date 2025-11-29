package services

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// ExperimentService handles A/B tests and feature flags
type ExperimentService struct {
	db *sqlx.DB
}

// NewExperimentService creates a new experiment service
func NewExperimentService(db *sqlx.DB) *ExperimentService {
	return &ExperimentService{db: db}
}

// Experiment represents an A/B test or experiment
type Experiment struct {
	ID             int64          `db:"id" json:"id"`
	TenantID       int64          `db:"tenant_id" json:"tenant_id"`
	ExperimentName string         `db:"experiment_name" json:"experiment_name"`
	ExperimentType string         `db:"experiment_type" json:"experiment_type"`
	Hypothesis     sql.NullString `db:"hypothesis" json:"hypothesis"`
	StartDate      time.Time      `db:"start_date" json:"start_date"`
	EndDate        sql.NullTime   `db:"end_date" json:"end_date"`
	Status         string         `db:"status" json:"status"`
	Variants       JSONB          `db:"variants" json:"variants"`
	SuccessMetrics JSONB          `db:"success_metrics" json:"success_metrics"`
	Results        JSONB          `db:"results" json:"results,omitempty"`
	Winner         sql.NullString `db:"winner" json:"winner"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at" json:"updated_at"`
}

// ExperimentAssignment represents user assignment to experiment variant
type ExperimentAssignment struct {
	ID           int64         `db:"id" json:"id"`
	ExperimentID int64         `db:"experiment_id" json:"experiment_id"`
	CustomerID   sql.NullInt64 `db:"customer_id" json:"customer_id"`
	SessionID    sql.NullString `db:"session_id" json:"session_id"`
	Variant      string        `db:"variant" json:"variant"`
	AssignedAt   time.Time     `db:"assigned_at" json:"assigned_at"`
}

// FeatureFlag represents a feature flag for gradual rollout
type FeatureFlag struct {
	ID                 int64     `db:"id" json:"id"`
	TenantID           int64     `db:"tenant_id" json:"tenant_id"`
	FlagName           string    `db:"flag_name" json:"flag_name"`
	FlagKey            string    `db:"flag_key" json:"flag_key"`
	Description        sql.NullString `db:"description" json:"description"`
	IsEnabled          bool      `db:"is_enabled" json:"is_enabled"`
	RolloutPercentage  int       `db:"rollout_percentage" json:"rollout_percentage"`
	TargetSegments     JSONB     `db:"target_segments" json:"target_segments,omitempty"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
}

// CreateExperiment creates a new experiment
func (s *ExperimentService) CreateExperiment(tenantID int64, experiment *Experiment) error {
	experiment.TenantID = tenantID
	experiment.CreatedAt = time.Now()
	experiment.UpdatedAt = time.Now()

	query := `
		INSERT INTO experiments (
			tenant_id, experiment_name, experiment_type, hypothesis,
			start_date, end_date, status, variants, success_metrics,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`

	err := s.db.QueryRow(
		query,
		experiment.TenantID, experiment.ExperimentName, experiment.ExperimentType,
		experiment.Hypothesis, experiment.StartDate, experiment.EndDate,
		experiment.Status, experiment.Variants, experiment.SuccessMetrics,
		experiment.CreatedAt, experiment.UpdatedAt,
	).Scan(&experiment.ID)

	return err
}

// GetExperiment retrieves an experiment
func (s *ExperimentService) GetExperiment(tenantID, experimentID int64) (*Experiment, error) {
	var experiment Experiment
	query := `SELECT * FROM experiments WHERE id = $1 AND tenant_id = $2`
	err := s.db.Get(&experiment, query, experimentID, tenantID)
	return &experiment, err
}

// ListExperiments lists all experiments
func (s *ExperimentService) ListExperiments(tenantID int64, status string) ([]Experiment, error) {
	query := `SELECT * FROM experiments WHERE tenant_id = $1`
	args := []interface{}{tenantID}

	if status != "" {
		query += ` AND status = $2`
		args = append(args, status)
	}

	query += ` ORDER BY created_at DESC`

	var experiments []Experiment
	err := s.db.Select(&experiments, query, args...)
	return experiments, err
}

// UpdateExperiment updates an experiment
func (s *ExperimentService) UpdateExperiment(tenantID, experimentID int64, updates map[string]interface{}) error {
	query := `UPDATE experiments SET updated_at = NOW()`
	args := []interface{}{}
	argIdx := 1

	for key, value := range updates {
		query += `, ` + key + ` = $` + string(rune(argIdx))
		args = append(args, value)
		argIdx++
	}

	query += ` WHERE id = $` + string(rune(argIdx)) + ` AND tenant_id = $` + string(rune(argIdx+1))
	args = append(args, experimentID, tenantID)

	_, err := s.db.Exec(query, args...)
	return err
}

// AssignVariant assigns a user to an experiment variant
func (s *ExperimentService) AssignVariant(experimentID int64, customerID *int64, sessionID *string) (*ExperimentAssignment, error) {
	// Get experiment
	var experiment Experiment
	err := s.db.Get(&experiment, `SELECT * FROM experiments WHERE id = $1`, experimentID)
	if err != nil {
		return nil, err
	}

	// Check if already assigned
	var existing ExperimentAssignment
	checkQuery := `SELECT * FROM experiment_assignments WHERE experiment_id = $1 AND `
	if customerID != nil {
		checkQuery += `customer_id = $2`
		err = s.db.Get(&existing, checkQuery, experimentID, *customerID)
	} else if sessionID != nil {
		checkQuery += `session_id = $2`
		err = s.db.Get(&existing, checkQuery, experimentID, *sessionID)
	}

	if err == nil {
		// Already assigned
		return &existing, nil
	}

	// Assign to a variant (simplified - use hash-based assignment in production)
	variant := "control" // Default

	assignment := &ExperimentAssignment{
		ExperimentID: experimentID,
		Variant:      variant,
		AssignedAt:   time.Now(),
	}

	if customerID != nil {
		assignment.CustomerID = sql.NullInt64{Int64: *customerID, Valid: true}
	}
	if sessionID != nil {
		assignment.SessionID = sql.NullString{String: *sessionID, Valid: true}
	}

	query := `
		INSERT INTO experiment_assignments (
			experiment_id, customer_id, session_id, variant, assigned_at
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err = s.db.QueryRow(
		query,
		assignment.ExperimentID, assignment.CustomerID,
		assignment.SessionID, assignment.Variant, assignment.AssignedAt,
	).Scan(&assignment.ID)

	return assignment, err
}

// GetExperimentResults calculates experiment results
func (s *ExperimentService) GetExperimentResults(tenantID, experimentID int64) (map[string]interface{}, error) {
	// Get variant performance
	query := `
		SELECT 
			variant,
			COUNT(*) as assignments,
			COUNT(CASE WHEN s.converted = true THEN 1 END) as conversions,
			AVG(CASE WHEN s.converted = true THEN s.conversion_value END) as avg_value
		FROM experiment_assignments ea
		LEFT JOIN sessions s ON ea.session_id = s.session_id
		WHERE ea.experiment_id = $1
		GROUP BY variant`

	rows, err := s.db.Query(query, experimentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := map[string]interface{}{
		"variants": []map[string]interface{}{},
	}

	variantResults := []map[string]interface{}{}
	for rows.Next() {
		var variant string
		var assignments, conversions int
		var avgValue sql.NullFloat64

		rows.Scan(&variant, &assignments, &conversions, &avgValue)

		conversionRate := 0.0
		if assignments > 0 {
			conversionRate = float64(conversions) / float64(assignments) * 100
		}

		variantResults = append(variantResults, map[string]interface{}{
			"variant":         variant,
			"assignments":     assignments,
			"conversions":     conversions,
			"conversion_rate": conversionRate,
			"avg_value":       avgValue.Float64,
		})
	}

	results["variants"] = variantResults
	return results, nil
}

// CreateFeatureFlag creates a new feature flag
func (s *ExperimentService) CreateFeatureFlag(tenantID int64, flag *FeatureFlag) error {
	flag.TenantID = tenantID
	flag.CreatedAt = time.Now()
	flag.UpdatedAt = time.Now()

	query := `
		INSERT INTO feature_flags (
			tenant_id, flag_name, flag_key, description, is_enabled,
			rollout_percentage, target_segments, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	err := s.db.QueryRow(
		query,
		flag.TenantID, flag.FlagName, flag.FlagKey, flag.Description,
		flag.IsEnabled, flag.RolloutPercentage, flag.TargetSegments,
		flag.CreatedAt, flag.UpdatedAt,
	).Scan(&flag.ID)

	return err
}

// GetFeatureFlag retrieves a feature flag by key
func (s *ExperimentService) GetFeatureFlag(tenantID int64, flagKey string) (*FeatureFlag, error) {
	var flag FeatureFlag
	query := `SELECT * FROM feature_flags WHERE flag_key = $1 AND tenant_id = $2`
	err := s.db.Get(&flag, query, flagKey, tenantID)
	return &flag, err
}

// ListFeatureFlags lists all feature flags
func (s *ExperimentService) ListFeatureFlags(tenantID int64) ([]FeatureFlag, error) {
	query := `SELECT * FROM feature_flags WHERE tenant_id = $1 ORDER BY created_at DESC`
	var flags []FeatureFlag
	err := s.db.Select(&flags, query, tenantID)
	return flags, err
}

// UpdateFeatureFlag updates a feature flag
func (s *ExperimentService) UpdateFeatureFlag(tenantID int64, flagKey string, updates map[string]interface{}) error {
	query := `UPDATE feature_flags SET updated_at = NOW()`
	args := []interface{}{}
	argIdx := 1

	for key, value := range updates {
		query += `, ` + key + ` = $` + string(rune(argIdx))
		args = append(args, value)
		argIdx++
	}

	query += ` WHERE flag_key = $` + string(rune(argIdx)) + ` AND tenant_id = $` + string(rune(argIdx+1))
	args = append(args, flagKey, tenantID)

	_, err := s.db.Exec(query, args...)
	return err
}

// IsFeatureEnabled checks if a feature is enabled for a user
func (s *ExperimentService) IsFeatureEnabled(tenantID int64, flagKey string, customerID *int64) (bool, error) {
	flag, err := s.GetFeatureFlag(tenantID, flagKey)
	if err != nil {
		return false, err
	}

	if !flag.IsEnabled {
		return false, nil
	}

	// Check rollout percentage (simplified - use consistent hashing in production)
	if flag.RolloutPercentage >= 100 {
		return true, nil
	}

	if flag.RolloutPercentage <= 0 {
		return false, nil
	}

	// Simplified rollout check
	return true, nil
}

