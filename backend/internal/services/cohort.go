package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// CohortService handles cohort analysis and segmentation
type CohortService struct {
	db *sqlx.DB
}

// NewCohortService creates a new cohort service
func NewCohortService(db *sqlx.DB) *CohortService {
	return &CohortService{db: db}
}

// CohortMetric represents computed cohort metrics
type CohortMetric struct {
	ID            int64     `db:"id" json:"id"`
	TenantID      int64     `db:"tenant_id" json:"tenant_id"`
	SegmentID     int64     `db:"segment_id" json:"segment_id"`
	CohortPeriod  string    `db:"cohort_period" json:"cohort_period"`
	PeriodOffset  int       `db:"period_offset" json:"period_offset"`
	MetricName    string    `db:"metric_name" json:"metric_name"`
	MetricValue   float64   `db:"metric_value" json:"metric_value"`
	CustomerCount int       `db:"customer_count" json:"customer_count"`
	ComputedAt    time.Time `db:"computed_at" json:"computed_at"`
}

// CohortAnalysisRequest represents a request for cohort analysis
type CohortAnalysisRequest struct {
	TenantID        int64     `json:"tenant_id"`
	SegmentID       int64     `json:"segment_id"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	CohortInterval  string    `json:"cohort_interval"` // daily, weekly, monthly
	MetricType      string    `json:"metric_type"`     // retention, revenue, conversions
	PeriodOffsets   int       `json:"period_offsets"`  // Number of periods to track
}

// CohortAnalysisResult represents the result of cohort analysis
type CohortAnalysisResult struct {
	CohortPeriod  string             `json:"cohort_period"`
	InitialSize   int                `json:"initial_size"`
	Metrics       map[int]float64    `json:"metrics"` // period_offset -> metric_value
}

// ComputeCohortMetrics computes cohort metrics for a segment
func (s *CohortService) ComputeCohortMetrics(req *CohortAnalysisRequest) ([]CohortAnalysisResult, error) {
	// Get segment configuration
	segment, err := s.getSegment(req.TenantID, req.SegmentID)
	if err != nil {
		return nil, err
	}

	// Determine cohort periods
	cohortPeriods := s.generateCohortPeriods(req.StartDate, req.EndDate, req.CohortInterval)
	
	results := make([]CohortAnalysisResult, 0)

	for _, cohortPeriod := range cohortPeriods {
		// Get customers in this cohort
		cohortStart, cohortEnd := s.getCohortDateRange(cohortPeriod, req.CohortInterval)
		customerIDs, err := s.getCustomersInCohort(req.TenantID, req.SegmentID, cohortStart, cohortEnd, segment.CohortDateField.String)
		if err != nil {
			continue
		}

		result := CohortAnalysisResult{
			CohortPeriod: cohortPeriod,
			InitialSize:  len(customerIDs),
			Metrics:      make(map[int]float64),
		}

		// Calculate metrics for each period offset
		for offset := 0; offset <= req.PeriodOffsets; offset++ {
			metricValue, customerCount := s.calculateCohortMetric(
				req.TenantID,
				customerIDs,
				req.MetricType,
				cohortStart,
				offset,
				req.CohortInterval,
			)

			result.Metrics[offset] = metricValue

			// Save to database
			s.saveCohortMetric(req.TenantID, req.SegmentID, cohortPeriod, offset, req.MetricType, metricValue, customerCount)
		}

		results = append(results, result)
	}

	// Update segment last_computed_at
	s.updateSegmentComputedAt(req.TenantID, req.SegmentID)

	return results, nil
}

// getSegment retrieves segment information
func (s *CohortService) getSegment(tenantID, segmentID int64) (*struct {
	ID              int64          `db:"id"`
	Name            string         `db:"name"`
	SegmentType     string         `db:"segment_type"`
	CohortDateField sql.NullString `db:"cohort_date_field"`
	CohortInterval  sql.NullString `db:"cohort_interval"`
}, error) {
	var segment struct {
		ID              int64          `db:"id"`
		Name            string         `db:"name"`
		SegmentType     string         `db:"segment_type"`
		CohortDateField sql.NullString `db:"cohort_date_field"`
		CohortInterval  sql.NullString `db:"cohort_interval"`
	}
	
	query := `SELECT id, name, segment_type, cohort_date_field, cohort_interval FROM segments WHERE id = $1 AND tenant_id = $2`
	err := s.db.Get(&segment, query, segmentID, tenantID)
	return &segment, err
}

// generateCohortPeriods generates cohort period strings
func (s *CohortService) generateCohortPeriods(start, end time.Time, interval string) []string {
	periods := make([]string, 0)
	current := start

	for current.Before(end) || current.Equal(end) {
		var period string
		switch interval {
		case "daily":
			period = current.Format("2006-01-02")
			current = current.AddDate(0, 0, 1)
		case "weekly":
			year, week := current.ISOWeek()
			period = fmt.Sprintf("%d-W%02d", year, week)
			current = current.AddDate(0, 0, 7)
		case "monthly":
			period = current.Format("2006-01")
			current = current.AddDate(0, 1, 0)
		default:
			period = current.Format("2006-01")
			current = current.AddDate(0, 1, 0)
		}
		periods = append(periods, period)
	}

	return periods
}

// getCohortDateRange converts a cohort period string to date range
func (s *CohortService) getCohortDateRange(period, interval string) (time.Time, time.Time) {
	var start, end time.Time

	switch interval {
	case "daily":
		start, _ = time.Parse("2006-01-02", period)
		end = start.AddDate(0, 0, 1)
	case "weekly":
		// Parse format: 2024-W01
		start, _ = time.Parse("2006-W02", period)
		end = start.AddDate(0, 0, 7)
	case "monthly":
		start, _ = time.Parse("2006-01", period)
		end = start.AddDate(0, 1, 0)
	}

	return start, end
}

// getCustomersInCohort retrieves customers who joined in a specific cohort period
func (s *CohortService) getCustomersInCohort(tenantID, segmentID int64, start, end time.Time, dateField string) ([]int64, error) {
	// If dateField is not specified, use created_at
	if dateField == "" {
		dateField = "created_at"
	}

	query := fmt.Sprintf(`
		SELECT DISTINCT c.id
		FROM customers c
		INNER JOIN customer_segments cs ON c.id = cs.customer_id
		WHERE cs.segment_id = $1 
		AND c.tenant_id = $2
		AND c.%s >= $3 
		AND c.%s < $4`, dateField, dateField)

	var customerIDs []int64
	err := s.db.Select(&customerIDs, query, segmentID, tenantID, start, end)
	return customerIDs, err
}

// calculateCohortMetric calculates a specific metric for a cohort at a given offset
func (s *CohortService) calculateCohortMetric(tenantID int64, customerIDs []int64, metricType string, cohortStart time.Time, offset int, interval string) (float64, int) {
	if len(customerIDs) == 0 {
		return 0, 0
	}

	// Calculate the period start based on offset
	var periodStart, periodEnd time.Time
	switch interval {
	case "daily":
		periodStart = cohortStart.AddDate(0, 0, offset)
		periodEnd = periodStart.AddDate(0, 0, 1)
	case "weekly":
		periodStart = cohortStart.AddDate(0, 0, offset*7)
		periodEnd = periodStart.AddDate(0, 0, 7)
	case "monthly":
		periodStart = cohortStart.AddDate(0, offset, 0)
		periodEnd = periodStart.AddDate(0, 1, 0)
	}

	var metricValue float64
	var customerCount int

	switch metricType {
	case "retention":
		// Count customers who had activity in this period
		query := `
			SELECT COUNT(DISTINCT customer_id)
			FROM interactions
			WHERE tenant_id = $1 
			AND customer_id = ANY($2)
			AND interaction_time >= $3 
			AND interaction_time < $4`
		
		s.db.Get(&customerCount, query, tenantID, customerIDs, periodStart, periodEnd)
		metricValue = float64(customerCount) / float64(len(customerIDs)) * 100

	case "revenue":
		// Sum revenue generated in this period
		query := `
			SELECT COALESCE(SUM(revenue), 0), COUNT(DISTINCT customer_id)
			FROM conversion_events
			WHERE tenant_id = $1 
			AND customer_id = ANY($2)
			AND conversion_time >= $3 
			AND conversion_time < $4`
		
		s.db.QueryRow(query, tenantID, customerIDs, periodStart, periodEnd).Scan(&metricValue, &customerCount)

	case "conversions":
		// Count conversions in this period
		query := `
			SELECT COUNT(*), COUNT(DISTINCT customer_id)
			FROM conversion_events
			WHERE tenant_id = $1 
			AND customer_id = ANY($2)
			AND conversion_time >= $3 
			AND conversion_time < $4`
		
		s.db.QueryRow(query, tenantID, customerIDs, periodStart, periodEnd).Scan(&customerCount, &metricValue)
	}

	return metricValue, customerCount
}

// saveCohortMetric saves computed cohort metric to database
func (s *CohortService) saveCohortMetric(tenantID, segmentID int64, cohortPeriod string, offset int, metricName string, metricValue float64, customerCount int) error {
	query := `
		INSERT INTO cohort_metrics (
			tenant_id, segment_id, cohort_period, period_offset,
			metric_name, metric_value, customer_count, computed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		ON CONFLICT (tenant_id, segment_id, cohort_period, period_offset, metric_name) 
		DO UPDATE SET 
			metric_value = EXCLUDED.metric_value,
			customer_count = EXCLUDED.customer_count,
			computed_at = NOW()`

	_, err := s.db.Exec(query, tenantID, segmentID, cohortPeriod, offset, metricName, metricValue, customerCount)
	return err
}

// updateSegmentComputedAt updates the last computed timestamp
func (s *CohortService) updateSegmentComputedAt(tenantID, segmentID int64) error {
	query := `UPDATE segments SET last_computed_at = NOW() WHERE id = $1 AND tenant_id = $2`
	_, err := s.db.Exec(query, segmentID, tenantID)
	return err
}

// GetCohortMetrics retrieves saved cohort metrics
func (s *CohortService) GetCohortMetrics(tenantID, segmentID int64, metricName string) ([]CohortMetric, error) {
	query := `
		SELECT * FROM cohort_metrics
		WHERE tenant_id = $1 AND segment_id = $2 AND metric_name = $3
		ORDER BY cohort_period, period_offset`

	var metrics []CohortMetric
	err := s.db.Select(&metrics, query, tenantID, segmentID, metricName)
	return metrics, err
}

// GetRetentionCurve retrieves retention curve data for visualization
func (s *CohortService) GetRetentionCurve(tenantID, segmentID int64) (map[string][]float64, error) {
	metrics, err := s.GetCohortMetrics(tenantID, segmentID, "retention")
	if err != nil {
		return nil, err
	}

	curve := make(map[string][]float64)
	for _, metric := range metrics {
		if _, exists := curve[metric.CohortPeriod]; !exists {
			curve[metric.CohortPeriod] = make([]float64, 0)
		}
		curve[metric.CohortPeriod] = append(curve[metric.CohortPeriod], metric.MetricValue)
	}

	return curve, nil
}

// AddCustomersToSegment adds customers to a segment
func (s *CohortService) AddCustomersToSegment(tenantID, segmentID int64, customerIDs []int64) error {
	if len(customerIDs) == 0 {
		return nil
	}

	query := `
		INSERT INTO customer_segments (segment_id, customer_id, added_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (segment_id, customer_id) DO NOTHING`

	for _, customerID := range customerIDs {
		_, err := s.db.Exec(query, segmentID, customerID)
		if err != nil {
			return err
		}
	}

	return nil
}

// RemoveCustomersFromSegment removes customers from a segment
func (s *CohortService) RemoveCustomersFromSegment(tenantID, segmentID int64, customerIDs []int64) error {
	if len(customerIDs) == 0 {
		return nil
	}

	query := `DELETE FROM customer_segments WHERE segment_id = $1 AND customer_id = ANY($2)`
	_, err := s.db.Exec(query, segmentID, customerIDs)
	return err
}

// RefreshDynamicSegment refreshes membership of a dynamic segment based on criteria
func (s *CohortService) RefreshDynamicSegment(tenantID, segmentID int64) error {
	// Get segment criteria
	var criteria JSONB
	query := `SELECT criteria FROM segments WHERE id = $1 AND tenant_id = $2 AND segment_type = 'dynamic'`
	err := s.db.Get(&criteria, query, segmentID, tenantID)
	if err != nil {
		return err
	}

	// Build dynamic query based on criteria (simplified example)
	// In production, you'd parse criteria and build appropriate WHERE clause
	customerQuery := `
		SELECT DISTINCT c.id
		FROM customers c
		WHERE c.tenant_id = $1
		-- Add dynamic WHERE conditions based on criteria
	`

	var customerIDs []int64
	err = s.db.Select(&customerIDs, customerQuery, tenantID)
	if err != nil {
		return err
	}

	// Clear existing memberships
	deleteQuery := `DELETE FROM customer_segments WHERE segment_id = $1`
	_, err = s.db.Exec(deleteQuery, segmentID)
	if err != nil {
		return err
	}

	// Add new memberships
	return s.AddCustomersToSegment(tenantID, segmentID, customerIDs)
}

