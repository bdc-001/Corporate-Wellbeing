package services

import (
	"fmt"
	"math"
	"time"

	"github.com/convin/crae/internal/models"
	"github.com/jmoiron/sqlx"
)

type AttributionService struct {
	db *sqlx.DB
}

func NewAttributionService(db *sqlx.DB) *AttributionService {
	return &AttributionService{db: db}
}

// AttributionConfig represents configuration for an attribution run
type AttributionConfig struct {
	TimeWindowHours    int      `json:"time_window_hours"`
	IncludeChannels    []string `json:"include_channels"`
	EventTypes         []string `json:"event_types"`
	MinPurchaseAmount  float64  `json:"min_purchase_amount"`
}

// CreateAttributionRun creates a new attribution run
func (s *AttributionService) CreateAttributionRun(tenantID int64, modelCode string, name string, config AttributionConfig) (*models.AttributionRun, error) {
	// Get model ID
	var modelID int
	err := s.db.Get(&modelID, `SELECT id FROM attribution_models WHERE code = $1`, modelCode)
	if err != nil {
		return nil, fmt.Errorf("attribution model not found: %s", modelCode)
	}

	// Convert config to JSONB
	configJSON := models.JSONB{
		"time_window_hours":   config.TimeWindowHours,
		"include_channels":    config.IncludeChannels,
		"event_types":         config.EventTypes,
		"min_purchase_amount": config.MinPurchaseAmount,
	}

	var run models.AttributionRun
	err = s.db.QueryRowx(
		`INSERT INTO attribution_runs (tenant_id, model_id, name, config, status)
		 VALUES ($1, $2, $3, $4, 'pending')
		 RETURNING id, tenant_id, model_id, name, description, config, status, started_at, completed_at, created_at, updated_at`,
		tenantID, modelID, name, configJSON,
	).Scan(
		&run.ID, &run.TenantID, &run.ModelID, &run.Name, &run.Description,
		&run.Config, &run.Status, &run.StartedAt, &run.CompletedAt, &run.CreatedAt, &run.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create attribution run: %w", err)
	}

	return &run, nil
}

// ExecuteAttributionRun executes an attribution run
func (s *AttributionService) ExecuteAttributionRun(runID int64) error {
	// Get run details
	var run models.AttributionRun
	err := s.db.Get(&run, `SELECT * FROM attribution_runs WHERE id = $1`, runID)
	if err != nil {
		return fmt.Errorf("attribution run not found: %w", err)
	}

	// Update status to running
	now := time.Now()
	_, err = s.db.Exec(
		`UPDATE attribution_runs SET status = 'running', started_at = $1 WHERE id = $2`,
		now, runID,
	)
	if err != nil {
		return fmt.Errorf("failed to update run status: %w", err)
	}

	// Parse config
	config := AttributionConfig{
		TimeWindowHours:   72, // default
		IncludeChannels:   []string{},
		EventTypes:        []string{},
		MinPurchaseAmount: 0,
	}
	if run.Config != nil {
		if tw, ok := run.Config["time_window_hours"].(float64); ok {
			config.TimeWindowHours = int(tw)
		}
		if ch, ok := run.Config["include_channels"].([]interface{}); ok {
			for _, v := range ch {
				if s, ok := v.(string); ok {
					config.IncludeChannels = append(config.IncludeChannels, s)
				}
			}
		}
		if et, ok := run.Config["event_types"].([]interface{}); ok {
			for _, v := range et {
				if s, ok := v.(string); ok {
					config.EventTypes = append(config.EventTypes, s)
				}
			}
		}
		if mp, ok := run.Config["min_purchase_amount"].(float64); ok {
			config.MinPurchaseAmount = mp
		}
	}

	// Get model code
	var modelCode string
	err = s.db.Get(&modelCode, `SELECT code FROM attribution_models WHERE id = $1`, run.ModelID)
	if err != nil {
		return fmt.Errorf("failed to get model code: %w", err)
	}

	// Get conversion events for this tenant
	query := `
		SELECT ce.*
		FROM conversion_events ce
		WHERE ce.tenant_id = $1
	`
	args := []interface{}{run.TenantID}
	argPos := 2

	if len(config.EventTypes) > 0 {
		query += fmt.Sprintf(" AND ce.event_type = ANY($%d)", argPos)
		args = append(args, config.EventTypes)
		argPos++
	}
	if config.MinPurchaseAmount > 0 {
		query += fmt.Sprintf(" AND ce.amount_decimal >= $%d", argPos)
		args = append(args, config.MinPurchaseAmount)
		argPos++
	}

	var conversions []models.ConversionEvent
	err = s.db.Select(&conversions, query, args...)
	if err != nil {
		return fmt.Errorf("failed to get conversion events: %w", err)
	}

	// Process each conversion event
	for _, conversion := range conversions {
		err = s.attributeConversion(runID, &run, conversion, modelCode, config)
		if err != nil {
			// Log error but continue
			fmt.Printf("Error attributing conversion %d: %v\n", conversion.ID, err)
		}
	}

	// Update status to completed
	_, err = s.db.Exec(
		`UPDATE attribution_runs SET status = 'completed', completed_at = $1 WHERE id = $2`,
		time.Now(), runID,
	)
	if err != nil {
		return fmt.Errorf("failed to update run status: %w", err)
	}

	return nil
}

// attributeConversion attributes a single conversion event
func (s *AttributionService) attributeConversion(runID int64, run *models.AttributionRun, conversion models.ConversionEvent, modelCode string, config AttributionConfig) error {
	// Get interactions within the time window
	windowStart := conversion.OccurredAt.Add(-time.Duration(config.TimeWindowHours) * time.Hour)

	query := `
		SELECT i.*, 
		       COALESCE(MAX(CASE WHEN ip.participant_type = 'agent' THEN ip.agent_id END), NULL) as agent_id,
		       COALESCE(MAX(CASE WHEN ip.participant_type = 'agent' THEN t.id END), NULL) as team_id,
		       COALESCE(MAX(CASE WHEN ip.participant_type = 'agent' THEN v.id END), NULL) as vendor_id
		FROM interactions i
		LEFT JOIN interaction_participants ip ON i.id = ip.interaction_id AND ip.participant_type = 'agent'
		LEFT JOIN agents a ON ip.agent_id = a.id
		LEFT JOIN teams t ON a.team_id = t.id
		LEFT JOIN vendors v ON a.vendor_id = v.id
		WHERE i.customer_id = $1
		  AND i.started_at >= $2
		  AND i.started_at <= $3
		GROUP BY i.id
	`
	args := []interface{}{conversion.CustomerID, windowStart, conversion.OccurredAt}

	if len(config.IncludeChannels) > 0 {
		query += ` AND i.channel_id IN (SELECT id FROM channels WHERE name = ANY($4))`
		args = append(args, config.IncludeChannels)
	}

	query += ` ORDER BY i.started_at ASC`

	type InteractionWithAgent struct {
		models.Interaction
		AgentID  *int `db:"agent_id"`
		TeamID   *int `db:"team_id"`
		VendorID *int `db:"vendor_id"`
	}

	var interactions []InteractionWithAgent
	err := s.db.Select(&interactions, query, args...)
	if err != nil {
		return fmt.Errorf("failed to get interactions: %w", err)
	}

	if len(interactions) == 0 {
		return nil // No interactions to attribute
	}

	// Calculate attribution weights based on model
	weights := s.calculateWeights(len(interactions), modelCode, conversion.OccurredAt)

	// Insert attribution results
	tx, err := s.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for i, interaction := range interactions {
		weight := weights[i]
		attributedAmount := conversion.AmountDecimal * weight

		// Determine if this is primary touch
		isPrimaryTouch := false
		if modelCode == "FIRST_TOUCH" && i == 0 {
			isPrimaryTouch = true
		} else if modelCode == "LAST_TOUCH" && i == len(interactions)-1 {
			isPrimaryTouch = true
		}

		_, err = tx.Exec(
			`INSERT INTO attribution_results (
				tenant_id, attribution_run_id, conversion_event_id, interaction_id,
				customer_id, agent_id, team_id, vendor_id, model_id,
				attribution_weight, attributed_amount, is_primary_touch
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			conversion.TenantID, runID, conversion.ID, interaction.ID,
			conversion.CustomerID, interaction.AgentID, interaction.TeamID, interaction.VendorID, run.ModelID,
			weight, attributedAmount, isPrimaryTouch,
		)
		if err != nil {
			return fmt.Errorf("failed to insert attribution result: %w", err)
		}
	}

	return tx.Commit()
}

// calculateWeights calculates attribution weights based on the model
func (s *AttributionService) calculateWeights(n int, modelCode string, conversionTime time.Time) []float64 {
	weights := make([]float64, n)

	switch modelCode {
	case "FIRST_TOUCH":
		weights[0] = 1.0
		for i := 1; i < n; i++ {
			weights[i] = 0.0
		}

	case "LAST_TOUCH":
		for i := 0; i < n-1; i++ {
			weights[i] = 0.0
		}
		weights[n-1] = 1.0

	case "LINEAR":
		weight := 1.0 / float64(n)
		for i := range weights {
			weights[i] = weight
		}

	case "TIME_DECAY":
		// Exponential decay: more recent interactions get more weight
		totalWeight := 0.0
		decayWeights := make([]float64, n)
		for i := 0; i < n; i++ {
			// Get interaction time (simplified - would need actual interaction struct)
			// For now, use position-based decay
			decayWeights[i] = math.Exp(-float64(n-i-1) * 0.5)
			totalWeight += decayWeights[i]
		}
		for i := range weights {
			weights[i] = decayWeights[i] / totalWeight
		}

	case "AI_WEIGHTED":
		// Simplified AI-weighted: use purchase probability if available
		// In production, this would use a trained model
		totalWeight := 0.0
		aiWeights := make([]float64, n)
		for i := 0; i < n; i++ {
			// Default weight based on position
			aiWeights[i] = 1.0 / float64(n)
			// Could enhance with purchase_probability from interaction
			totalWeight += aiWeights[i]
		}
		for i := range weights {
			weights[i] = aiWeights[i] / totalWeight
		}

	default:
		// Default to linear
		weight := 1.0 / float64(n)
		for i := range weights {
			weights[i] = weight
		}
	}

	return weights
}

