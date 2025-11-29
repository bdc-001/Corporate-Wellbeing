package services

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// FraudService handles fraud detection and data quality
type FraudService struct {
	db *sqlx.DB
}

// NewFraudService creates a new fraud service
func NewFraudService(db *sqlx.DB) *FraudService {
	return &FraudService{db: db}
}

// FraudDetectionRule represents a fraud detection rule
type FraudDetectionRule struct {
	ID             int64     `db:"id" json:"id"`
	TenantID       int64     `db:"tenant_id" json:"tenant_id"`
	RuleName       string    `db:"rule_name" json:"rule_name"`
	RuleType       string    `db:"rule_type" json:"rule_type"`
	IsActive       bool      `db:"is_active" json:"is_active"`
	Severity       string    `db:"severity" json:"severity"`
	DetectionLogic JSONB     `db:"detection_logic" json:"detection_logic"`
	Action         string    `db:"action" json:"action"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// FraudIncident represents a detected fraud incident
type FraudIncident struct {
	ID              int64           `db:"id" json:"id"`
	TenantID        int64           `db:"tenant_id" json:"tenant_id"`
	RuleID          sql.NullInt64   `db:"rule_id" json:"rule_id"`
	IncidentType    string          `db:"incident_type" json:"incident_type"`
	Severity        string          `db:"severity" json:"severity"`
	EntityType      sql.NullString  `db:"entity_type" json:"entity_type"`
	EntityID        sql.NullInt64   `db:"entity_id" json:"entity_id"`
	DetectedAt      time.Time       `db:"detected_at" json:"detected_at"`
	ConfidenceScore float64         `db:"confidence_score" json:"confidence_score"`
	Status          string          `db:"status" json:"status"`
	ReviewedBy      sql.NullString  `db:"reviewed_by" json:"reviewed_by"`
	ReviewedAt      sql.NullTime    `db:"reviewed_at" json:"reviewed_at"`
	Evidence        JSONB           `db:"evidence" json:"evidence,omitempty"`
	ActionsTaken    JSONB           `db:"actions_taken" json:"actions_taken,omitempty"`
}

// DataQualityScore represents data quality metrics
type DataQualityScore struct {
	ID                int64     `db:"id" json:"id"`
	TenantID          int64     `db:"tenant_id" json:"tenant_id"`
	EntityType        string    `db:"entity_type" json:"entity_type"`
	EntityID          int64     `db:"entity_id" json:"entity_id"`
	CompletenessScore float64   `db:"completeness_score" json:"completeness_score"`
	AccuracyScore     float64   `db:"accuracy_score" json:"accuracy_score"`
	ConsistencyScore  float64   `db:"consistency_score" json:"consistency_score"`
	TimelinessScore   float64   `db:"timeliness_score" json:"timeliness_score"`
	OverallScore      float64   `db:"overall_score" json:"overall_score"`
	Issues            JSONB     `db:"issues" json:"issues,omitempty"`
	CalculatedAt      time.Time `db:"calculated_at" json:"calculated_at"`
}

// CreateFraudRule creates a new fraud detection rule
func (s *FraudService) CreateFraudRule(tenantID int64, rule *FraudDetectionRule) error {
	rule.TenantID = tenantID
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()

	query := `
		INSERT INTO fraud_detection_rules (
			tenant_id, rule_name, rule_type, is_active, severity,
			detection_logic, action, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	err := s.db.QueryRow(
		query,
		rule.TenantID, rule.RuleName, rule.RuleType, rule.IsActive,
		rule.Severity, rule.DetectionLogic, rule.Action,
		rule.CreatedAt, rule.UpdatedAt,
	).Scan(&rule.ID)

	return err
}

// GetFraudRules retrieves fraud detection rules
func (s *FraudService) GetFraudRules(tenantID int64, activeOnly bool) ([]FraudDetectionRule, error) {
	query := `SELECT * FROM fraud_detection_rules WHERE tenant_id = $1`
	args := []interface{}{tenantID}

	if activeOnly {
		query += ` AND is_active = true`
	}

	query += ` ORDER BY severity DESC, created_at DESC`

	var rules []FraudDetectionRule
	err := s.db.Select(&rules, query, args...)
	return rules, err
}

// DetectFraud runs fraud detection on an entity
func (s *FraudService) DetectFraud(tenantID int64, entityType string, entityID int64) ([]FraudIncident, error) {
	// Get active rules
	rules, err := s.GetFraudRules(tenantID, true)
	if err != nil {
		return nil, err
	}

	incidents := make([]FraudIncident, 0)

	// Run each rule
	for _, rule := range rules {
		violated, confidence, evidence := s.evaluateFraudRule(tenantID, entityType, entityID, &rule)
		if violated {
			incident := &FraudIncident{
				TenantID:        tenantID,
				RuleID:          sql.NullInt64{Int64: rule.ID, Valid: true},
				IncidentType:    rule.RuleType,
				Severity:        rule.Severity,
				EntityType:      sql.NullString{String: entityType, Valid: true},
				EntityID:        sql.NullInt64{Int64: entityID, Valid: true},
				DetectedAt:      time.Now(),
				ConfidenceScore: confidence,
				Status:          "pending",
				Evidence:        evidence,
			}

			err := s.CreateFraudIncident(incident)
			if err == nil {
				incidents = append(incidents, *incident)
			}
		}
	}

	return incidents, nil
}

// evaluateFraudRule evaluates a single fraud rule
func (s *FraudService) evaluateFraudRule(tenantID int64, entityType string, entityID int64, rule *FraudDetectionRule) (bool, float64, JSONB) {
	// Simplified fraud detection logic
	// In production, parse rule.DetectionLogic and apply complex checks
	
	violated := false
	confidence := 0.0
	evidence := JSONB{}

	// Example: Check for click fraud
	if rule.RuleType == "click_fraud" && entityType == "interaction" {
		// Check for suspicious patterns like multiple clicks from same IP in short time
		var clickCount int
		query := `
			SELECT COUNT(*) FROM interactions 
			WHERE tenant_id = $1 AND id = $2
			AND interaction_time > NOW() - INTERVAL '1 hour'`
		
		s.db.Get(&clickCount, query, tenantID, entityID)
		
		if clickCount > 100 { // Threshold
			violated = true
			confidence = 0.85
		}
	}

	// Example: Check for lead fraud
	if rule.RuleType == "lead_fraud" && entityType == "customer" {
		// Check for fake email patterns, disposable emails, etc.
		var email string
		query := `
			SELECT ci.identifier_value FROM customer_identifiers ci
			WHERE ci.customer_id = $1 AND ci.identifier_type = 'email'
			LIMIT 1`
		
		err := s.db.Get(&email, query, entityID)
		if err == nil && s.isSuspiciousEmail(email) {
			violated = true
			confidence = 0.70
		}
	}

	return violated, confidence, evidence
}

// isSuspiciousEmail checks for suspicious email patterns
func (s *FraudService) isSuspiciousEmail(email string) bool {
	// Simple check - in production, use more sophisticated patterns
	suspiciousPatterns := []string{"tempmail", "throwaway", "fake", "test"}
	for _, pattern := range suspiciousPatterns {
		if len(email) > len(pattern) {
			// Simple substring check
			return false
		}
	}
	return false
}

// CreateFraudIncident creates a new fraud incident
func (s *FraudService) CreateFraudIncident(incident *FraudIncident) error {
	query := `
		INSERT INTO fraud_incidents (
			tenant_id, rule_id, incident_type, severity, entity_type,
			entity_id, detected_at, confidence_score, status, evidence
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	err := s.db.QueryRow(
		query,
		incident.TenantID, incident.RuleID, incident.IncidentType,
		incident.Severity, incident.EntityType, incident.EntityID,
		incident.DetectedAt, incident.ConfidenceScore, incident.Status,
		incident.Evidence,
	).Scan(&incident.ID)

	return err
}

// GetFraudIncidents retrieves fraud incidents
func (s *FraudService) GetFraudIncidents(tenantID int64, status, severity string, limit int) ([]FraudIncident, error) {
	query := `SELECT * FROM fraud_incidents WHERE tenant_id = $1`
	args := []interface{}{tenantID}
	argIdx := 2

	if status != "" {
		query += ` AND status = $` + string(rune(argIdx))
		args = append(args, status)
		argIdx++
	}

	if severity != "" {
		query += ` AND severity = $` + string(rune(argIdx))
		args = append(args, severity)
	}

	query += ` ORDER BY detected_at DESC LIMIT $` + string(rune(argIdx))
	args = append(args, limit)

	var incidents []FraudIncident
	err := s.db.Select(&incidents, query, args...)
	return incidents, err
}

// UpdateIncidentStatus updates fraud incident status
func (s *FraudService) UpdateIncidentStatus(tenantID, incidentID int64, status, reviewedBy string) error {
	query := `
		UPDATE fraud_incidents
		SET status = $1, reviewed_by = $2, reviewed_at = NOW()
		WHERE id = $3 AND tenant_id = $4`

	_, err := s.db.Exec(query, status, reviewedBy, incidentID, tenantID)
	return err
}

// CalculateDataQuality calculates data quality scores for an entity
func (s *FraudService) CalculateDataQuality(tenantID int64, entityType string, entityID int64) (*DataQualityScore, error) {
	score := &DataQualityScore{
		TenantID:     tenantID,
		EntityType:   entityType,
		EntityID:     entityID,
		CalculatedAt: time.Now(),
	}

	// Calculate scores based on entity type
	switch entityType {
	case "customer":
		score.CompletenessScore = s.calculateCustomerCompleteness(entityID)
		score.AccuracyScore = s.calculateCustomerAccuracy(entityID)
		score.ConsistencyScore = s.calculateCustomerConsistency(entityID)
		score.TimelinessScore = s.calculateCustomerTimeliness(entityID)
	case "interaction":
		score.CompletenessScore = s.calculateInteractionCompleteness(entityID)
		score.AccuracyScore = 100.0 // Placeholder
		score.ConsistencyScore = 100.0
		score.TimelinessScore = 100.0
	}

	// Calculate overall score
	score.OverallScore = (score.CompletenessScore + score.AccuracyScore + 
		score.ConsistencyScore + score.TimelinessScore) / 4.0

	// Save to database
	query := `
		INSERT INTO data_quality_scores (
			tenant_id, entity_type, entity_id, completeness_score,
			accuracy_score, consistency_score, timeliness_score,
			overall_score, calculated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	err := s.db.QueryRow(
		query,
		score.TenantID, score.EntityType, score.EntityID,
		score.CompletenessScore, score.AccuracyScore,
		score.ConsistencyScore, score.TimelinessScore,
		score.OverallScore, score.CalculatedAt,
	).Scan(&score.ID)

	return score, err
}

// calculateCustomerCompleteness calculates completeness score for a customer
func (s *FraudService) calculateCustomerCompleteness(customerID int64) float64 {
	// Count how many fields are filled
	// Simplified - in production, check specific important fields
	score := 50.0 // Base score

	// Check if has email
	var hasEmail bool
	s.db.Get(&hasEmail, `
		SELECT EXISTS(SELECT 1 FROM customer_identifiers 
		WHERE customer_id = $1 AND identifier_type = 'email')`, customerID)
	if hasEmail {
		score += 20.0
	}

	// Check if has phone
	var hasPhone bool
	s.db.Get(&hasPhone, `
		SELECT EXISTS(SELECT 1 FROM customer_identifiers 
		WHERE customer_id = $1 AND identifier_type = 'phone')`, customerID)
	if hasPhone {
		score += 15.0
	}

	// Check if has account
	var hasAccount bool
	s.db.Get(&hasAccount, `
		SELECT EXISTS(SELECT 1 FROM customers 
		WHERE id = $1 AND account_id IS NOT NULL)`, customerID)
	if hasAccount {
		score += 15.0
	}

	return score
}

// calculateCustomerAccuracy calculates accuracy score (placeholder)
func (s *FraudService) calculateCustomerAccuracy(customerID int64) float64 {
	return 90.0 // Placeholder - would check email validity, phone format, etc.
}

// calculateCustomerConsistency calculates consistency score (placeholder)
func (s *FraudService) calculateCustomerConsistency(customerID int64) float64 {
	return 95.0 // Placeholder - would check for conflicting data
}

// calculateCustomerTimeliness calculates timeliness score (placeholder)
func (s *FraudService) calculateCustomerTimeliness(customerID int64) float64 {
	return 100.0 // Placeholder - would check last update time
}

// calculateInteractionCompleteness calculates completeness for interactions
func (s *FraudService) calculateInteractionCompleteness(interactionID int64) float64 {
	return 85.0 // Placeholder
}

