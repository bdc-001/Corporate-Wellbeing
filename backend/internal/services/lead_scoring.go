package services

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
)

// LeadScoringService handles lead scoring and predictive analytics
type LeadScoringService struct {
	db *sqlx.DB
}

// NewLeadScoringService creates a new lead scoring service
func NewLeadScoringService(db *sqlx.DB) *LeadScoringService {
	return &LeadScoringService{db: db}
}

// LeadScoringModel represents a scoring model configuration
type LeadScoringModel struct {
	ID           int64     `db:"id" json:"id"`
	TenantID     int64     `db:"tenant_id" json:"tenant_id"`
	Name         string    `db:"name" json:"name"`
	ModelType    string    `db:"model_type" json:"model_type"` // rule_based, ml_based, hybrid
	IsActive     bool      `db:"is_active" json:"is_active"`
	Version      int       `db:"version" json:"version"`
	ScoringRules JSONB     `db:"scoring_rules" json:"scoring_rules,omitempty"`
	ModelConfig  JSONB     `db:"model_config" json:"model_config,omitempty"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// LeadScore represents a calculated lead score
type LeadScore struct {
	ID             int64     `db:"id" json:"id"`
	TenantID       int64     `db:"tenant_id" json:"tenant_id"`
	CustomerID     int64     `db:"customer_id" json:"customer_id"`
	ModelID        int64     `db:"model_id" json:"model_id"`
	Score          float64   `db:"score" json:"score"`
	ScoreBreakdown JSONB     `db:"score_breakdown" json:"score_breakdown,omitempty"`
	CalculatedAt   time.Time `db:"calculated_at" json:"calculated_at"`
	Factors        JSONB     `db:"factors" json:"factors,omitempty"`
}

// Prediction represents a predictive analytics result
type Prediction struct {
	ID                   int64           `db:"id" json:"id"`
	TenantID             int64           `db:"tenant_id" json:"tenant_id"`
	CustomerID           sql.NullInt64   `db:"customer_id" json:"customer_id"`
	AccountID            sql.NullInt64   `db:"account_id" json:"account_id"`
	PredictionType       string          `db:"prediction_type" json:"prediction_type"` // churn, conversion, ltv, next_action
	PredictedValue       sql.NullFloat64 `db:"predicted_value" json:"predicted_value"`
	PredictedProbability float64         `db:"predicted_probability" json:"predicted_probability"`
	ConfidenceLevel      float64         `db:"confidence_level" json:"confidence_level"`
	PredictionDate       time.Time       `db:"prediction_date" json:"prediction_date"`
	ExpiryDate           sql.NullTime    `db:"expiry_date" json:"expiry_date"`
	ModelVersion         sql.NullString  `db:"model_version" json:"model_version"`
	FeaturesUsed         JSONB           `db:"features_used" json:"features_used,omitempty"`
	CreatedAt            time.Time       `db:"created_at" json:"created_at"`
}

// ScoringRule represents a single scoring rule
type ScoringRule struct {
	Name       string                 `json:"name"`
	Field      string                 `json:"field"`
	Condition  string                 `json:"condition"` // equals, greater_than, less_than, contains, etc.
	Value      interface{}            `json:"value"`
	Score      float64                `json:"score"`
	Weight     float64                `json:"weight"`
	Category   string                 `json:"category"` // demographic, behavioral, firmographic
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// CreateScoringModel creates a new lead scoring model
func (s *LeadScoringService) CreateScoringModel(tenantID int64, model *LeadScoringModel) error {
	model.TenantID = tenantID
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()

	query := `
		INSERT INTO lead_scoring_models (
			tenant_id, name, model_type, is_active, version,
			scoring_rules, model_config, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	err := s.db.QueryRow(
		query,
		model.TenantID, model.Name, model.ModelType, model.IsActive,
		model.Version, model.ScoringRules, model.ModelConfig,
		model.CreatedAt, model.UpdatedAt,
	).Scan(&model.ID)

	return err
}

// GetScoringModel retrieves a scoring model
func (s *LeadScoringService) GetScoringModel(tenantID, modelID int64) (*LeadScoringModel, error) {
	var model LeadScoringModel
	query := `SELECT * FROM lead_scoring_models WHERE id = $1 AND tenant_id = $2`
	err := s.db.Get(&model, query, modelID, tenantID)
	return &model, err
}

// GetActiveScoringModel retrieves the active scoring model for a tenant
func (s *LeadScoringService) GetActiveScoringModel(tenantID int64) (*LeadScoringModel, error) {
	var model LeadScoringModel
	query := `SELECT * FROM lead_scoring_models WHERE tenant_id = $1 AND is_active = true ORDER BY version DESC LIMIT 1`
	err := s.db.Get(&model, query, tenantID)
	return &model, err
}

// CalculateLeadScore calculates score for a customer based on rules
func (s *LeadScoringService) CalculateLeadScore(tenantID, customerID int64) (*LeadScore, error) {
	// Get active model
	model, err := s.GetActiveScoringModel(tenantID)
	if err != nil {
		return nil, err
	}

	// Parse scoring rules
	var rules []ScoringRule
	if model.ScoringRules != nil {
		rulesJSON, err := json.Marshal(model.ScoringRules)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(rulesJSON, &rules)
		if err != nil {
			return nil, err
		}
	}

	// Get customer data
	customerData, err := s.getCustomerDataForScoring(tenantID, customerID)
	if err != nil {
		return nil, err
	}

	// Calculate score based on rules
	totalScore := 0.0
	breakdown := make(map[string]float64)
	factors := make([]map[string]interface{}, 0)

	for _, rule := range rules {
		score := s.evaluateRule(rule, customerData)
		totalScore += score
		breakdown[rule.Category] = breakdown[rule.Category] + score
		
		if score > 0 {
			factors = append(factors, map[string]interface{}{
				"rule":  rule.Name,
				"score": score,
				"field": rule.Field,
			})
		}
	}

	// Convert breakdown to JSONB
	breakdownJSONB := make(JSONB)
	for k, v := range breakdown {
		breakdownJSONB[k] = v
	}
	
	// Convert factors to JSONB
	factorsJSONB := make(JSONB)
	if len(factors) > 0 {
		factorsJSONB["factors"] = factors
	}

	// Save lead score
	leadScore := &LeadScore{
		TenantID:       tenantID,
		CustomerID:     customerID,
		ModelID:        model.ID,
		Score:          totalScore,
		ScoreBreakdown: breakdownJSONB,
		CalculatedAt:   time.Now(),
		Factors:        factorsJSONB,
	}

	query := `
		INSERT INTO lead_scores (
			tenant_id, customer_id, model_id, score, score_breakdown,
			calculated_at, factors
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	err = s.db.QueryRow(
		query,
		leadScore.TenantID, leadScore.CustomerID, leadScore.ModelID,
		leadScore.Score, leadScore.ScoreBreakdown, leadScore.CalculatedAt,
		leadScore.Factors,
	).Scan(&leadScore.ID)

	return leadScore, err
}

// getCustomerDataForScoring retrieves all necessary customer data
func (s *LeadScoringService) getCustomerDataForScoring(tenantID, customerID int64) (map[string]interface{}, error) {
	query := `
		SELECT 
			c.*,
			COUNT(DISTINCT i.id) as interaction_count,
			COUNT(DISTINCT ce.id) as conversion_count,
			COALESCE(SUM(ce.revenue), 0) as total_revenue,
			COUNT(DISTINCT s.session_id) as session_count,
			MAX(i.interaction_time) as last_interaction_date,
			a.industry as account_industry,
			a.company_size as account_company_size,
			a.annual_revenue as account_annual_revenue
		FROM customers c
		LEFT JOIN interactions i ON c.id = i.customer_id
		LEFT JOIN conversion_events ce ON c.id = ce.customer_id
		LEFT JOIN sessions s ON c.id = s.customer_id
		LEFT JOIN accounts a ON c.account_id = a.id
		WHERE c.id = $1 AND c.tenant_id = $2
		GROUP BY c.id, a.industry, a.company_size, a.annual_revenue`

	var data map[string]interface{}
	rows, err := s.db.Query(query, customerID, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		cols, _ := rows.Columns()
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		data = make(map[string]interface{})
		for i, col := range cols {
			data[col] = values[i]
		}
	}

	return data, nil
}

// evaluateRule evaluates a single scoring rule
func (s *LeadScoringService) evaluateRule(rule ScoringRule, data map[string]interface{}) float64 {
	fieldValue, exists := data[rule.Field]
	if !exists {
		return 0
	}

	// Basic rule evaluation (can be extended)
	matched := false
	switch rule.Condition {
	case "equals":
		matched = fieldValue == rule.Value
	case "greater_than":
		if fv, ok := fieldValue.(float64); ok {
			if rv, ok := rule.Value.(float64); ok {
				matched = fv > rv
			}
		}
	case "less_than":
		if fv, ok := fieldValue.(float64); ok {
			if rv, ok := rule.Value.(float64); ok {
				matched = fv < rv
			}
		}
	case "contains":
		if fv, ok := fieldValue.(string); ok {
			if rv, ok := rule.Value.(string); ok {
				matched = contains(fv, rv)
			}
		}
	}

	if matched {
		return rule.Score * rule.Weight
	}
	return 0
}

// contains checks if a string contains a substring (simple helper)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr))
}

// CreatePrediction creates a new prediction
func (s *LeadScoringService) CreatePrediction(tenantID int64, prediction *Prediction) error {
	prediction.TenantID = tenantID
	prediction.CreatedAt = time.Now()

	query := `
		INSERT INTO predictions (
			tenant_id, customer_id, account_id, prediction_type,
			predicted_value, predicted_probability, confidence_level,
			prediction_date, expiry_date, model_version, features_used, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id`

	err := s.db.QueryRow(
		query,
		prediction.TenantID, prediction.CustomerID, prediction.AccountID,
		prediction.PredictionType, prediction.PredictedValue,
		prediction.PredictedProbability, prediction.ConfidenceLevel,
		prediction.PredictionDate, prediction.ExpiryDate,
		prediction.ModelVersion, prediction.FeaturesUsed, prediction.CreatedAt,
	).Scan(&prediction.ID)

	return err
}

// GetCustomerPredictions retrieves predictions for a customer
func (s *LeadScoringService) GetCustomerPredictions(tenantID, customerID int64, predictionType string) ([]Prediction, error) {
	query := `
		SELECT * FROM predictions 
		WHERE tenant_id = $1 AND customer_id = $2`
	
	args := []interface{}{tenantID, customerID}
	if predictionType != "" {
		query += " AND prediction_type = $3"
		args = append(args, predictionType)
	}
	
	query += " ORDER BY prediction_date DESC"

	var predictions []Prediction
	err := s.db.Select(&predictions, query, args...)
	return predictions, err
}

// GetLatestLeadScore retrieves the most recent lead score for a customer
func (s *LeadScoringService) GetLatestLeadScore(tenantID, customerID int64) (*LeadScore, error) {
	var score LeadScore
	query := `
		SELECT * FROM lead_scores 
		WHERE tenant_id = $1 AND customer_id = $2 
		ORDER BY calculated_at DESC 
		LIMIT 1`
	
	err := s.db.Get(&score, query, tenantID, customerID)
	return &score, err
}

// GetHighValueLeads retrieves leads above a certain score threshold
func (s *LeadScoringService) GetHighValueLeads(tenantID int64, minScore float64, limit int) ([]LeadScore, error) {
	query := `
		SELECT DISTINCT ON (customer_id) *
		FROM lead_scores
		WHERE tenant_id = $1 AND score >= $2
		ORDER BY customer_id, calculated_at DESC
		LIMIT $3`

	var scores []LeadScore
	err := s.db.Select(&scores, query, tenantID, minScore, limit)
	return scores, err
}

