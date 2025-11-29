package services

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type AnalyticsService struct {
	db *sqlx.DB
}

func NewAnalyticsService(db *sqlx.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

// AgentRevenueSummary represents revenue summary for an agent
type AgentRevenueSummary struct {
	AgentID                        int     `json:"agent_id" db:"agent_id"`
	Name                           string  `json:"name" db:"name"`
	Email                          string  `json:"email" db:"email"`
	VendorID                       *int    `json:"vendor_id" db:"vendor_id"`
	VendorName                     string  `json:"vendor_name" db:"vendor_name"`
	TeamID                         *int    `json:"team_id" db:"team_id"`
	TeamName                       string  `json:"team_name" db:"team_name"`
	TotalAttributedAmount          float64 `json:"total_attributed_amount" db:"total_attributed_amount"`
	TotalConversions               int     `json:"total_conversions" db:"total_conversions"`
	AvgAttributedAmountPerInteraction float64 `json:"avg_attributed_amount_per_interaction" db:"avg_attributed_amount_per_interaction"`
}

// GetAgentRevenueSummary returns revenue summary for agents
func (s *AnalyticsService) GetAgentRevenueSummary(tenantID int64, from, to *time.Time, vendorID *int, modelCode string) ([]AgentRevenueSummary, error) {
	// Check if attribution_results exist, otherwise use fallback
	var hasAttribution int
	_ = s.db.Get(&hasAttribution, `SELECT COUNT(*) FROM attribution_results WHERE tenant_id = $1 LIMIT 1`, tenantID)

	var query string
	args := []interface{}{tenantID}
	argPos := 2

	if hasAttribution > 0 {
		// Use attribution_results
		query = `
			SELECT 
				a.id as agent_id,
				a.name,
				a.email,
				a.vendor_id,
				v.name as vendor_name,
				a.team_id,
				t.name as team_name,
				COALESCE(SUM(ar.attributed_amount), 0) as total_attributed_amount,
				COUNT(DISTINCT ar.conversion_event_id) as total_conversions,
				COALESCE(SUM(ar.attributed_amount) / NULLIF(COUNT(DISTINCT ar.interaction_id), 0), 0) as avg_attributed_amount_per_interaction
			FROM agents a
			LEFT JOIN vendors v ON a.vendor_id = v.id
			LEFT JOIN teams t ON a.team_id = t.id
			LEFT JOIN attribution_results ar ON a.id = ar.agent_id AND ar.tenant_id = $1
			LEFT JOIN attribution_runs run ON ar.attribution_run_id = run.id
			LEFT JOIN attribution_models am ON run.model_id = am.id
			LEFT JOIN conversion_events ce ON ar.conversion_event_id = ce.id
			WHERE a.vendor_id IN (SELECT id FROM vendors WHERE tenant_id = $1)
		`

		if modelCode != "" {
			query += fmt.Sprintf(" AND am.code = $%d", argPos)
			args = append(args, modelCode)
			argPos++
		}

		if from != nil {
			query += fmt.Sprintf(" AND ce.occurred_at >= $%d", argPos)
			args = append(args, *from)
			argPos++
		}

		if to != nil {
			query += fmt.Sprintf(" AND ce.occurred_at <= $%d", argPos)
			args = append(args, *to)
			argPos++
		}

		query += `
			GROUP BY a.id, a.name, a.email, a.vendor_id, v.name, a.team_id, t.name
			HAVING COUNT(DISTINCT ar.conversion_event_id) > 0
			ORDER BY total_attributed_amount DESC
		`
	} else {
		// Fallback: Link interactions -> agents -> customers -> conversions
		query = `
			SELECT 
				a.id as agent_id,
				a.name,
				a.email,
				a.vendor_id,
				COALESCE(v.name, '') as vendor_name,
				a.team_id,
				COALESCE(t.name, '') as team_name,
				COALESCE(SUM(ce.amount_decimal), 0) as total_attributed_amount,
				COUNT(DISTINCT ce.id) as total_conversions,
				COALESCE(SUM(ce.amount_decimal) / NULLIF(COUNT(DISTINCT i.id), 0), 0) as avg_attributed_amount_per_interaction
			FROM agents a
			LEFT JOIN vendors v ON a.vendor_id = v.id
			LEFT JOIN teams t ON a.team_id = t.id
			LEFT JOIN interaction_participants ip ON a.id = ip.agent_id AND ip.participant_type = 'agent'
			LEFT JOIN interactions i ON ip.interaction_id = i.id AND i.tenant_id = $1
			LEFT JOIN customers c ON i.customer_id = c.id
			LEFT JOIN conversion_events ce ON c.id = ce.customer_id AND ce.tenant_id = $1
			WHERE a.vendor_id IN (SELECT id FROM vendors WHERE tenant_id = $1)
		`

		if vendorID != nil {
			query += fmt.Sprintf(" AND a.vendor_id = $%d", argPos)
			args = append(args, *vendorID)
			argPos++
		}

		if from != nil {
			query += fmt.Sprintf(" AND ce.occurred_at >= $%d", argPos)
			args = append(args, *from)
			argPos++
		}

		if to != nil {
			query += fmt.Sprintf(" AND ce.occurred_at <= $%d", argPos)
			args = append(args, *to)
			argPos++
		}

		query += `
			GROUP BY a.id, a.name, a.email, a.vendor_id, COALESCE(v.name, ''), a.team_id, COALESCE(t.name, '')
			HAVING COUNT(DISTINCT ce.id) > 0
			ORDER BY total_attributed_amount DESC
		`
	}

	var results []AgentRevenueSummary
	err := s.db.Select(&results, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent revenue summary: %w", err)
	}

	// Ensure we return an empty slice, not nil
	if results == nil {
		results = []AgentRevenueSummary{}
	}

	return results, nil
}

// VendorComparison represents vendor comparison data
type VendorComparison struct {
	VendorID            int     `json:"vendor_id" db:"vendor_id"`
	Name                string  `json:"name" db:"name"`
	TotalAttributedAmount float64 `json:"total_attributed_amount" db:"total_attributed_amount"`
	TotalConversions    int     `json:"total_conversions" db:"total_conversions"`
	AvgConversionValue  float64 `json:"avg_conversion_value" db:"avg_conversion_value"`
}

// GetVendorComparison returns comparison data for vendors
func (s *AnalyticsService) GetVendorComparison(tenantID int64, from, to *time.Time, modelCode string) ([]VendorComparison, error) {
	query := `
		WITH vendor_conversions AS (
			-- Use attribution_results if available
			SELECT 
				ar.vendor_id,
				SUM(ar.attributed_amount) as total_amount,
				COUNT(DISTINCT ar.conversion_event_id) as conversions
			FROM attribution_results ar
			LEFT JOIN attribution_runs run ON ar.attribution_run_id = run.id
			LEFT JOIN attribution_models am ON run.model_id = am.id
			LEFT JOIN conversion_events ce ON ar.conversion_event_id = ce.id
			WHERE ar.tenant_id = $1 AND ar.vendor_id IS NOT NULL
	`
	
	args := []interface{}{tenantID}
	argPos := 2

	if modelCode != "" {
		query += fmt.Sprintf(" AND am.code = $%d", argPos)
		args = append(args, modelCode)
		argPos++
	}

	if from != nil {
		query += fmt.Sprintf(" AND ce.occurred_at >= $%d", argPos)
		args = append(args, *from)
		argPos++
	}

	if to != nil {
		query += fmt.Sprintf(" AND ce.occurred_at <= $%d", argPos)
		args = append(args, *to)
		argPos++
	}

	query += `
			GROUP BY ar.vendor_id
		),
		vendor_conversions_fallback AS (
			-- Fallback: Link interactions -> agents -> vendors -> customers -> conversions
			SELECT 
				v.id as vendor_id,
				SUM(ce.amount_decimal) as total_amount,
				COUNT(DISTINCT ce.id) as conversions
			FROM interactions i
			INNER JOIN interaction_participants ip ON i.id = ip.interaction_id AND ip.participant_type = 'agent'
			INNER JOIN agents a ON ip.agent_id = a.id
			INNER JOIN vendors v ON a.vendor_id = v.id
			INNER JOIN customers c ON i.customer_id = c.id
			INNER JOIN conversion_events ce ON c.id = ce.customer_id AND ce.tenant_id = $1
			WHERE i.tenant_id = $1
	`

	if from != nil {
		query += fmt.Sprintf(" AND ce.occurred_at >= $%d", argPos)
		args = append(args, *from)
		argPos++
	}

	if to != nil {
		query += fmt.Sprintf(" AND ce.occurred_at <= $%d", argPos)
		args = append(args, *to)
		argPos++
	}

	query += `
			GROUP BY v.id
		),
		combined_vendor_data AS (
			SELECT 
				COALESCE(vc.vendor_id, vcf.vendor_id) as vendor_id,
				COALESCE(vc.total_amount, vcf.total_amount, 0) as total_amount,
				COALESCE(vc.conversions, vcf.conversions, 0) as conversions
			FROM vendor_conversions vc
			FULL OUTER JOIN vendor_conversions_fallback vcf ON vc.vendor_id = vcf.vendor_id
		)
		SELECT 
			v.id as vendor_id,
			v.name,
			COALESCE(cvd.total_amount, 0) as total_attributed_amount,
			COALESCE(cvd.conversions, 0) as total_conversions,
			COALESCE(cvd.total_amount / NULLIF(cvd.conversions, 0), 0) as avg_conversion_value
		FROM vendors v
		LEFT JOIN combined_vendor_data cvd ON v.id = cvd.vendor_id
		WHERE v.tenant_id = $1
		GROUP BY v.id, v.name, cvd.total_amount, cvd.conversions
		HAVING COALESCE(cvd.conversions, 0) > 0
		ORDER BY total_attributed_amount DESC
	`

	var results []VendorComparison
	err := s.db.Select(&results, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get vendor comparison: %w", err)
	}

	return results, nil
}

// IntentProfitability represents intent-level profitability data
type IntentProfitability struct {
	IntentCode              string  `json:"intent_code" db:"intent_code"`
	TotalAttributedAmount   float64 `json:"total_attributed_amount" db:"total_attributed_amount"`
	TotalConversions        int     `json:"total_conversions" db:"total_conversions"`
	AvgHandleTimeSeconds    float64 `json:"avg_handle_time_seconds" db:"avg_handle_time_seconds"`
	ProfitabilityScore      float64 `json:"profitability_score" db:"profitability_score"`
}

// GetIntentProfitability returns intent-level profitability data
func (s *AnalyticsService) GetIntentProfitability(tenantID int64, from, to *time.Time, modelCode string) ([]IntentProfitability, error) {
	query := `
		WITH intent_conversions AS (
			-- Use attribution_results if available
			SELECT 
				i.primary_intent,
				SUM(ar.attributed_amount) as total_amount,
				COUNT(DISTINCT ar.conversion_event_id) as conversions,
				AVG(i.duration_seconds) as avg_duration
			FROM interactions i
			INNER JOIN attribution_results ar ON i.id = ar.interaction_id AND ar.tenant_id = $1
			LEFT JOIN attribution_runs run ON ar.attribution_run_id = run.id
			LEFT JOIN attribution_models am ON run.model_id = am.id
			LEFT JOIN conversion_events ce ON ar.conversion_event_id = ce.id
			WHERE i.tenant_id = $1
			  AND i.primary_intent IS NOT NULL
			  AND i.primary_intent != ''
	`
	
	args := []interface{}{tenantID}
	argPos := 2

	if modelCode != "" {
		query += fmt.Sprintf(" AND am.code = $%d", argPos)
		args = append(args, modelCode)
		argPos++
	}

	if from != nil {
		query += fmt.Sprintf(" AND ce.occurred_at >= $%d", argPos)
		args = append(args, *from)
		argPos++
	}

	if to != nil {
		query += fmt.Sprintf(" AND ce.occurred_at <= $%d", argPos)
		args = append(args, *to)
		argPos++
	}

	query += `
			GROUP BY i.primary_intent
		),
		intent_conversions_fallback AS (
			-- Fallback: Link interactions -> customers -> conversions
			SELECT 
				i.primary_intent,
				SUM(ce.amount_decimal) as total_amount,
				COUNT(DISTINCT ce.id) as conversions,
				AVG(i.duration_seconds) as avg_duration
			FROM interactions i
			INNER JOIN customers c ON i.customer_id = c.id
			INNER JOIN conversion_events ce ON c.id = ce.customer_id AND ce.tenant_id = $1
			WHERE i.tenant_id = $1
			  AND i.primary_intent IS NOT NULL
			  AND i.primary_intent != ''
	`

	if from != nil {
		query += fmt.Sprintf(" AND ce.occurred_at >= $%d", argPos)
		args = append(args, *from)
		argPos++
	}

	if to != nil {
		query += fmt.Sprintf(" AND ce.occurred_at <= $%d", argPos)
		args = append(args, *to)
		argPos++
	}

	query += `
			GROUP BY i.primary_intent
		),
		combined_intent_data AS (
			SELECT 
				COALESCE(ic.primary_intent, icf.primary_intent) as primary_intent,
				COALESCE(ic.total_amount, icf.total_amount, 0) as total_amount,
				COALESCE(ic.conversions, icf.conversions, 0) as conversions,
				COALESCE(ic.avg_duration, icf.avg_duration, 0) as avg_duration
			FROM intent_conversions ic
			FULL OUTER JOIN intent_conversions_fallback icf ON ic.primary_intent = icf.primary_intent
		)
		SELECT 
			cid.primary_intent as intent_code,
			COALESCE(cid.total_amount, 0) as total_attributed_amount,
			COALESCE(cid.conversions, 0) as total_conversions,
			COALESCE(cid.avg_duration, 0) as avg_handle_time_seconds,
			COALESCE(
				cid.total_amount / NULLIF(cid.avg_duration, 0) * 1000,
				0
			) as profitability_score
		FROM combined_intent_data cid
		WHERE cid.primary_intent IS NOT NULL
		  AND cid.primary_intent != ''
		  AND COALESCE(cid.conversions, 0) > 0
		ORDER BY total_attributed_amount DESC
	`

	var results []IntentProfitability
	err := s.db.Select(&results, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get intent profitability: %w", err)
	}

	return results, nil
}

