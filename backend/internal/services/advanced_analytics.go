package services

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// AdvancedAnalyticsService provides Factors.ai-style analytics
type AdvancedAnalyticsService struct {
	db *sqlx.DB
}

func NewAdvancedAnalyticsService(db *sqlx.DB) *AdvancedAnalyticsService {
	return &AdvancedAnalyticsService{db: db}
}

// FunnelStageMetrics represents metrics for each funnel stage
type FunnelStageMetrics struct {
	Stage              string    `json:"stage" db:"stage"`
	TotalAccounts      int       `json:"total_accounts" db:"total_accounts"`
	TotalRevenue       float64   `json:"total_revenue" db:"total_revenue"`
	AvgTimeInStage     float64   `json:"avg_time_in_stage" db:"avg_time_in_stage"` // in hours
	ConversionRate     float64   `json:"conversion_rate" db:"conversion_rate"`
	DropOffRate        float64   `json:"drop_off_rate" db:"drop_off_rate"`
}

// GetFunnelStageMetrics returns metrics for each funnel stage
func (s *AdvancedAnalyticsService) GetFunnelStageMetrics(tenantID int64, from, to *time.Time, segmentID *int) ([]FunnelStageMetrics, error) {
	query := `
		WITH stage_data AS (
			SELECT 
				COALESCE(i.funnel_stage, ce.funnel_stage, 'Unknown') as stage,
				COUNT(DISTINCT i.customer_id) as total_accounts,
				COALESCE(SUM(ar.attributed_amount), 0) as total_revenue,
				AVG(EXTRACT(EPOCH FROM (i.ended_at - i.started_at))/3600) as avg_time_in_stage
			FROM interactions i
			LEFT JOIN attribution_results ar ON i.id = ar.interaction_id AND ar.tenant_id = $1
			LEFT JOIN conversion_events ce ON ar.conversion_event_id = ce.id
			WHERE i.tenant_id = $1
	`
	
	args := []interface{}{tenantID}
	argPos := 2

	if from != nil {
		query += fmt.Sprintf(" AND i.started_at >= $%d", argPos)
		args = append(args, *from)
		argPos++
	}

	if to != nil {
		query += fmt.Sprintf(" AND i.started_at <= $%d", argPos)
		args = append(args, *to)
		argPos++
	}

	if segmentID != nil {
		query += fmt.Sprintf(`
			AND i.customer_id IN (
				SELECT customer_id FROM customer_segments WHERE segment_id = $%d
			)`, argPos)
		args = append(args, *segmentID)
		argPos++
	}

	query += `
			GROUP BY stage
		)
		SELECT 
			stage,
			total_accounts,
			total_revenue,
			COALESCE(avg_time_in_stage, 0) as avg_time_in_stage,
			CASE 
				WHEN LAG(total_accounts) OVER (ORDER BY 
					CASE stage
						WHEN 'MQL' THEN 1
						WHEN 'SQL' THEN 2
						WHEN 'Opportunity' THEN 3
						WHEN 'Closed-Won' THEN 4
						ELSE 5
					END
				) > 0 THEN
					(total_accounts::float / LAG(total_accounts) OVER (ORDER BY 
						CASE stage
							WHEN 'MQL' THEN 1
							WHEN 'SQL' THEN 2
							WHEN 'Opportunity' THEN 3
							WHEN 'Closed-Won' THEN 4
							ELSE 5
						END
					)) * 100
				ELSE 0
			END as conversion_rate,
			CASE 
				WHEN LAG(total_accounts) OVER (ORDER BY 
					CASE stage
						WHEN 'MQL' THEN 1
						WHEN 'SQL' THEN 2
						WHEN 'Opportunity' THEN 3
						WHEN 'Closed-Won' THEN 4
						ELSE 5
					END
				) > 0 THEN
					((LAG(total_accounts) OVER (ORDER BY 
						CASE stage
							WHEN 'MQL' THEN 1
							WHEN 'SQL' THEN 2
							WHEN 'Opportunity' THEN 3
							WHEN 'Closed-Won' THEN 4
							ELSE 5
						END
					) - total_accounts)::float / LAG(total_accounts) OVER (ORDER BY 
						CASE stage
							WHEN 'MQL' THEN 1
							WHEN 'SQL' THEN 2
							WHEN 'Opportunity' THEN 3
							WHEN 'Closed-Won' THEN 4
							ELSE 5
						END
					)) * 100
				ELSE 0
			END as drop_off_rate
		FROM stage_data
		ORDER BY 
			CASE stage
				WHEN 'MQL' THEN 1
				WHEN 'SQL' THEN 2
				WHEN 'Opportunity' THEN 3
				WHEN 'Closed-Won' THEN 4
				ELSE 5
			END
	`

	var results []FunnelStageMetrics
	err := s.db.Select(&results, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get funnel stage metrics: %w", err)
	}

	return results, nil
}

// ContentEngagementMetrics represents content engagement analytics
type ContentEngagementMetrics struct {
	ContentID          int     `json:"content_id" db:"content_id"`
	ContentName        string  `json:"content_name" db:"content_name"`
	ContentType        string  `json:"content_type" db:"content_type"`
	TotalEngagements   int     `json:"total_engagements" db:"total_engagements"`
	UniqueAccounts     int     `json:"unique_accounts" db:"unique_accounts"`
	AttributedRevenue  float64 `json:"attributed_revenue" db:"attributed_revenue"`
	ConversionRate     float64 `json:"conversion_rate" db:"conversion_rate"`
}

// GetContentEngagementMetrics returns content engagement analytics
func (s *AdvancedAnalyticsService) GetContentEngagementMetrics(tenantID int64, from, to *time.Time) ([]ContentEngagementMetrics, error) {
	query := `
		SELECT 
			ca.id as content_id,
			ca.name as content_name,
			ca.content_type,
			COUNT(ce.id) as total_engagements,
			COUNT(DISTINCT ce.customer_id) as unique_accounts,
			COALESCE(SUM(ar.attributed_amount), 0) as attributed_revenue,
			CASE 
				WHEN COUNT(DISTINCT ce.customer_id) > 0 THEN
					(COUNT(DISTINCT ar.conversion_event_id)::float / COUNT(DISTINCT ce.customer_id)) * 100
				ELSE 0
			END as conversion_rate
		FROM content_assets ca
		LEFT JOIN content_engagements ce ON ca.id = ce.content_id AND ce.tenant_id = $1
		LEFT JOIN attribution_results ar ON ce.interaction_id = ar.interaction_id AND ar.tenant_id = $1
		WHERE ca.tenant_id = $1
	`

	args := []interface{}{tenantID}
	argPos := 2

	if from != nil {
		query += fmt.Sprintf(" AND ce.engaged_at >= $%d", argPos)
		args = append(args, *from)
		argPos++
	}

	if to != nil {
		query += fmt.Sprintf(" AND ce.engaged_at <= $%d", argPos)
		args = append(args, *to)
		argPos++
	}

	query += `
		GROUP BY ca.id, ca.name, ca.content_type
		HAVING COUNT(ce.id) > 0
		ORDER BY attributed_revenue DESC
	`

	var results []ContentEngagementMetrics
	err := s.db.Select(&results, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get content engagement metrics: %w", err)
	}

	return results, nil
}

// MultiChannelROI represents ROI across different ad platforms
type MultiChannelROI struct {
	Platform           string  `json:"platform" db:"platform"`
	TotalSpend         float64 `json:"total_spend" db:"total_spend"`
	TotalImpressions   int64   `json:"total_impressions" db:"total_impressions"`
	TotalClicks        int64   `json:"total_clicks" db:"total_clicks"`
	AttributedRevenue  float64 `json:"attributed_revenue" db:"attributed_revenue"`
	ROI                float64 `json:"roi" db:"roi"`
	CPC                float64 `json:"cpc" db:"cpc"`
	CPM                float64 `json:"cpm" db:"cpm"`
}

// GetMultiChannelROI returns ROI metrics across ad platforms
func (s *AdvancedAnalyticsService) GetMultiChannelROI(tenantID int64, from, to *time.Time) ([]MultiChannelROI, error) {
	query := `
		SELECT 
			COALESCE(as_spend.platform, 'Unknown') as platform,
			COALESCE(SUM(as_spend.spend_amount), 0) as total_spend,
			COALESCE(SUM(as_spend.impressions), 0) as total_impressions,
			COALESCE(SUM(as_spend.clicks), 0) as total_clicks,
			COALESCE(SUM(ar.attributed_amount), 0) as attributed_revenue,
			CASE 
				WHEN SUM(as_spend.spend_amount) > 0 THEN
					((SUM(ar.attributed_amount) - SUM(as_spend.spend_amount)) / SUM(as_spend.spend_amount)) * 100
				ELSE 0
			END as roi,
			CASE 
				WHEN SUM(as_spend.clicks) > 0 THEN
					SUM(as_spend.spend_amount) / SUM(as_spend.clicks)
				ELSE 0
			END as cpc,
			CASE 
				WHEN SUM(as_spend.impressions) > 0 THEN
					(SUM(as_spend.spend_amount) / SUM(as_spend.impressions)) * 1000
				ELSE 0
			END as cpm
		FROM ad_spend as_spend
		LEFT JOIN campaigns c ON as_spend.campaign_id = c.id
		LEFT JOIN interactions i ON i.campaign_id = c.id AND i.tenant_id = $1
		LEFT JOIN attribution_results ar ON i.id = ar.interaction_id AND ar.tenant_id = $1
		WHERE as_spend.tenant_id = $1
	`

	args := []interface{}{tenantID}
	argPos := 2

	if from != nil {
		query += fmt.Sprintf(" AND as_spend.date >= $%d", argPos)
		args = append(args, *from)
		argPos++
	}

	if to != nil {
		query += fmt.Sprintf(" AND as_spend.date <= $%d", argPos)
		args = append(args, *to)
		argPos++
	}

	query += `
		GROUP BY as_spend.platform
		ORDER BY attributed_revenue DESC
	`

	var results []MultiChannelROI
	err := s.db.Select(&results, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get multi-channel ROI: %w", err)
	}

	return results, nil
}

// JourneyVelocity represents time between funnel stages
type JourneyVelocity struct {
	FromStage          string  `json:"from_stage" db:"from_stage"`
	ToStage            string  `json:"to_stage" db:"to_stage"`
	AvgDays            float64 `json:"avg_days" db:"avg_days"`
	MedianDays         float64 `json:"median_days" db:"median_days"`
	Accounts           int     `json:"accounts" db:"accounts"`
}

// GetJourneyVelocity returns velocity metrics between stages
func (s *AdvancedAnalyticsService) GetJourneyVelocity(tenantID int64, from, to *time.Time) ([]JourneyVelocity, error) {
	query := `
		WITH stage_transitions AS (
			SELECT 
				i1.customer_id,
				i1.funnel_stage as from_stage,
				i2.funnel_stage as to_stage,
				EXTRACT(EPOCH FROM (i2.started_at - i1.started_at)) / 86400 as days_between
			FROM interactions i1
			INNER JOIN interactions i2 ON i1.customer_id = i2.customer_id 
				AND i2.started_at > i1.started_at
				AND i2.funnel_stage != i1.funnel_stage
			WHERE i1.tenant_id = $1
				AND i1.funnel_stage IS NOT NULL
				AND i2.funnel_stage IS NOT NULL
	`
	
	args := []interface{}{tenantID}
	argPos := 2

	if from != nil {
		query += fmt.Sprintf(" AND i1.started_at >= $%d", argPos)
		args = append(args, *from)
		argPos++
	}

	if to != nil {
		query += fmt.Sprintf(" AND i1.started_at <= $%d", argPos)
		args = append(args, *to)
		argPos++
	}

	query += `
		)
		SELECT 
			from_stage,
			to_stage,
			AVG(days_between) as avg_days,
			PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY days_between) as median_days,
			COUNT(DISTINCT customer_id) as accounts
		FROM stage_transitions
		GROUP BY from_stage, to_stage
		ORDER BY from_stage, to_stage
	`

	var results []JourneyVelocity
	err := s.db.Select(&results, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get journey velocity: %w", err)
	}

	return results, nil
}

// CustomReportRequest represents a custom report with filters
type CustomReportRequest struct {
	Filters struct {
		Channels      []string `json:"channels"`
		Campaigns     []int    `json:"campaigns"`
		FunnelStages  []string `json:"funnel_stages"`
		Segments      []int    `json:"segments"`
		Vendors       []int    `json:"vendors"`
		Intents       []string `json:"intents"`
	} `json:"filters"`
	GroupBy []string `json:"group_by"` // channel, campaign, funnel_stage, segment, vendor, intent
	Metrics []string `json:"metrics"`  // revenue, conversions, velocity, roi
}

// CustomReportResult represents custom report results
type CustomReportResult struct {
	Groups  map[string]interface{} `json:"groups"`
	Metrics map[string]float64     `json:"metrics"`
}

// GetCustomReport generates a custom report based on filters
func (s *AdvancedAnalyticsService) GetCustomReport(tenantID int64, req CustomReportRequest, from, to *time.Time) ([]map[string]interface{}, error) {
	// This is a simplified version - in production, build dynamic SQL based on filters
	query := `
		SELECT 
			c.name as channel_name,
			camp.name as campaign_name,
			COALESCE(i.funnel_stage, 'Unknown') as funnel_stage,
			COALESCE(SUM(ar.attributed_amount), 0) as revenue,
			COUNT(DISTINCT ar.conversion_event_id) as conversions
		FROM interactions i
		LEFT JOIN channels c ON i.channel_id = c.id
		LEFT JOIN campaigns camp ON i.campaign_id = camp.id
		LEFT JOIN attribution_results ar ON i.id = ar.interaction_id AND ar.tenant_id = $1
		WHERE i.tenant_id = $1
	`

	args := []interface{}{tenantID}
	argPos := 2

	if from != nil {
		query += fmt.Sprintf(" AND i.started_at >= $%d", argPos)
		args = append(args, *from)
		argPos++
	}

	if to != nil {
		query += fmt.Sprintf(" AND i.started_at <= $%d", argPos)
		args = append(args, *to)
		argPos++
	}

	// Apply filters
	if len(req.Filters.Channels) > 0 {
		query += fmt.Sprintf(" AND c.name = ANY($%d)", argPos)
		args = append(args, req.Filters.Channels)
		argPos++
	}

	if len(req.Filters.Campaigns) > 0 {
		query += fmt.Sprintf(" AND camp.id = ANY($%d)", argPos)
		args = append(args, req.Filters.Campaigns)
		argPos++
	}

	if len(req.Filters.FunnelStages) > 0 {
		query += fmt.Sprintf(" AND i.funnel_stage = ANY($%d)", argPos)
		args = append(args, req.Filters.FunnelStages)
		argPos++
	}

	query += `
		GROUP BY c.name, camp.name, i.funnel_stage
		ORDER BY revenue DESC
	`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute custom report: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var channelName, campaignName, funnelStage string
		var revenue float64
		var conversions int

		err := rows.Scan(&channelName, &campaignName, &funnelStage, &revenue, &conversions)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		results = append(results, map[string]interface{}{
			"channel":     channelName,
			"campaign":    campaignName,
			"funnel_stage": funnelStage,
			"revenue":     revenue,
			"conversions": conversions,
		})
	}

	return results, nil
}

