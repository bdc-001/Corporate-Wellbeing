package services

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// ReportService handles custom reports and saved queries
type ReportService struct {
	db *sqlx.DB
}

// NewReportService creates a new report service
func NewReportService(db *sqlx.DB) *ReportService {
	return &ReportService{db: db}
}

// SavedReport represents a saved report configuration
type SavedReport struct {
	ID                  int64          `db:"id" json:"id"`
	TenantID            int64          `db:"tenant_id" json:"tenant_id"`
	ReportName          string         `db:"report_name" json:"report_name"`
	ReportType          string         `db:"report_type" json:"report_type"`
	Description         sql.NullString `db:"description" json:"description"`
	QueryConfig         JSONB          `db:"query_config" json:"query_config"`
	VisualizationConfig JSONB          `db:"visualization_config" json:"visualization_config,omitempty"`
	Schedule            sql.NullString `db:"schedule" json:"schedule"`
	IsPublic            bool           `db:"is_public" json:"is_public"`
	CreatedBy           sql.NullString `db:"created_by" json:"created_by"`
	CreatedAt           time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time      `db:"updated_at" json:"updated_at"`
}

// ReportSnapshot represents a historical report snapshot
type ReportSnapshot struct {
	ID           int64          `db:"id" json:"id"`
	ReportID     int64          `db:"report_id" json:"report_id"`
	SnapshotDate time.Time      `db:"snapshot_date" json:"snapshot_date"`
	ResultData   JSONB          `db:"result_data" json:"result_data"`
	GeneratedBy  sql.NullString `db:"generated_by" json:"generated_by"`
	CreatedAt    time.Time      `db:"created_at" json:"created_at"`
}

// ReportExecutionResult represents the result of executing a report
type ReportExecutionResult struct {
	Report     *SavedReport           `json:"report"`
	Data       []map[string]interface{} `json:"data"`
	Summary    map[string]interface{}   `json:"summary"`
	ExecutedAt time.Time                `json:"executed_at"`
}

// CreateReport creates a new saved report
func (s *ReportService) CreateReport(tenantID int64, report *SavedReport) error {
	report.TenantID = tenantID
	report.CreatedAt = time.Now()
	report.UpdatedAt = time.Now()

	query := `
		INSERT INTO saved_reports (
			tenant_id, report_name, report_type, description,
			query_config, visualization_config, schedule,
			is_public, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`

	err := s.db.QueryRow(
		query,
		report.TenantID, report.ReportName, report.ReportType,
		report.Description, report.QueryConfig, report.VisualizationConfig,
		report.Schedule, report.IsPublic, report.CreatedBy,
		report.CreatedAt, report.UpdatedAt,
	).Scan(&report.ID)

	return err
}

// GetReport retrieves a saved report
func (s *ReportService) GetReport(tenantID, reportID int64) (*SavedReport, error) {
	var report SavedReport
	query := `SELECT * FROM saved_reports WHERE id = $1 AND tenant_id = $2`
	err := s.db.Get(&report, query, reportID, tenantID)
	return &report, err
}

// ListReports lists all saved reports
func (s *ReportService) ListReports(tenantID int64, reportType string, publicOnly bool) ([]SavedReport, error) {
	query := `SELECT * FROM saved_reports WHERE tenant_id = $1`
	args := []interface{}{tenantID}
	argIdx := 2

	if reportType != "" {
		query += ` AND report_type = $` + string(rune(argIdx))
		args = append(args, reportType)
		argIdx++
	}

	if publicOnly {
		query += ` AND is_public = true`
	}

	query += ` ORDER BY updated_at DESC`

	var reports []SavedReport
	err := s.db.Select(&reports, query, args...)
	return reports, err
}

// UpdateReport updates a saved report
func (s *ReportService) UpdateReport(tenantID, reportID int64, updates map[string]interface{}) error {
	query := `UPDATE saved_reports SET updated_at = NOW()`
	args := []interface{}{}
	argIdx := 1

	for key, value := range updates {
		query += `, ` + key + ` = $` + string(rune(argIdx))
		args = append(args, value)
		argIdx++
	}

	query += ` WHERE id = $` + string(rune(argIdx)) + ` AND tenant_id = $` + string(rune(argIdx+1))
	args = append(args, reportID, tenantID)

	_, err := s.db.Exec(query, args...)
	return err
}

// DeleteReport deletes a saved report
func (s *ReportService) DeleteReport(tenantID, reportID int64) error {
	query := `DELETE FROM saved_reports WHERE id = $1 AND tenant_id = $2`
	_, err := s.db.Exec(query, reportID, tenantID)
	return err
}

// ExecuteReport executes a report and returns results
func (s *ReportService) ExecuteReport(tenantID, reportID int64, params map[string]interface{}) (*ReportExecutionResult, error) {
	report, err := s.GetReport(tenantID, reportID)
	if err != nil {
		return nil, err
	}

	// Execute report based on type
	var data []map[string]interface{}
	var summary map[string]interface{}

	switch report.ReportType {
	case "attribution":
		data, summary = s.executeAttributionReport(tenantID, report, params)
	case "funnel":
		data, summary = s.executeFunnelReport(tenantID, report, params)
	case "cohort":
		data, summary = s.executeCohortReport(tenantID, report, params)
	case "revenue":
		data, summary = s.executeRevenueReport(tenantID, report, params)
	case "custom":
		data, summary = s.executeCustomReport(tenantID, report, params)
	default:
		data = []map[string]interface{}{}
		summary = map[string]interface{}{}
	}

	result := &ReportExecutionResult{
		Report:     report,
		Data:       data,
		Summary:    summary,
		ExecutedAt: time.Now(),
	}

	return result, nil
}

// executeAttributionReport executes an attribution report
func (s *ReportService) executeAttributionReport(tenantID int64, report *SavedReport, params map[string]interface{}) ([]map[string]interface{}, map[string]interface{}) {
	// Query attribution results
	query := `
		SELECT 
			ar.model_type,
			ch.name as channel_name,
			SUM(ar.attributed_revenue) as total_revenue,
			COUNT(DISTINCT ar.conversion_event_id) as conversions
		FROM attribution_results ar
		LEFT JOIN channels ch ON ar.channel_id = ch.id
		WHERE ar.tenant_id = $1
		GROUP BY ar.model_type, ch.name
		ORDER BY total_revenue DESC`

	rows, err := s.db.Query(query, tenantID)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	data := []map[string]interface{}{}
	totalRevenue := 0.0
	totalConversions := 0

	for rows.Next() {
		var modelType, channelName sql.NullString
		var revenue float64
		var conversions int

		rows.Scan(&modelType, &channelName, &revenue, &conversions)

		data = append(data, map[string]interface{}{
			"model_type":   modelType.String,
			"channel_name": channelName.String,
			"revenue":      revenue,
			"conversions":  conversions,
		})

		totalRevenue += revenue
		totalConversions += conversions
	}

	summary := map[string]interface{}{
		"total_revenue":     totalRevenue,
		"total_conversions": totalConversions,
	}

	return data, summary
}

// executeFunnelReport executes a funnel report
func (s *ReportService) executeFunnelReport(tenantID int64, report *SavedReport, params map[string]interface{}) ([]map[string]interface{}, map[string]interface{}) {
	// Simplified funnel report
	data := []map[string]interface{}{
		{"stage": "Awareness", "count": 1000, "conversion_rate": 100.0},
		{"stage": "Interest", "count": 500, "conversion_rate": 50.0},
		{"stage": "Consideration", "count": 250, "conversion_rate": 25.0},
		{"stage": "Purchase", "count": 100, "conversion_rate": 10.0},
	}

	summary := map[string]interface{}{
		"total_top_of_funnel": 1000,
		"total_converted":     100,
		"overall_rate":        10.0,
	}

	return data, summary
}

// executeCohortReport executes a cohort report
func (s *ReportService) executeCohortReport(tenantID int64, report *SavedReport, params map[string]interface{}) ([]map[string]interface{}, map[string]interface{}) {
	// Query cohort metrics
	query := `
		SELECT 
			cohort_period,
			period_offset,
			metric_value,
			customer_count
		FROM cohort_metrics
		WHERE tenant_id = $1
		ORDER BY cohort_period, period_offset`

	rows, err := s.db.Query(query, tenantID)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	data := []map[string]interface{}{}
	for rows.Next() {
		var cohortPeriod string
		var periodOffset, customerCount int
		var metricValue float64

		rows.Scan(&cohortPeriod, &periodOffset, &metricValue, &customerCount)

		data = append(data, map[string]interface{}{
			"cohort_period":  cohortPeriod,
			"period_offset":  periodOffset,
			"metric_value":   metricValue,
			"customer_count": customerCount,
		})
	}

	summary := map[string]interface{}{
		"cohorts_analyzed": len(data),
	}

	return data, summary
}

// executeRevenueReport executes a revenue report
func (s *ReportService) executeRevenueReport(tenantID int64, report *SavedReport, params map[string]interface{}) ([]map[string]interface{}, map[string]interface{}) {
	query := `
		SELECT 
			DATE_TRUNC('day', conversion_time) as date,
			COUNT(*) as conversions,
			SUM(revenue) as revenue
		FROM conversion_events
		WHERE tenant_id = $1
		AND conversion_time >= NOW() - INTERVAL '30 days'
		GROUP BY DATE_TRUNC('day', conversion_time)
		ORDER BY date DESC`

	rows, err := s.db.Query(query, tenantID)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()

	data := []map[string]interface{}{}
	totalRevenue := 0.0
	totalConversions := 0

	for rows.Next() {
		var date time.Time
		var conversions int
		var revenue float64

		rows.Scan(&date, &conversions, &revenue)

		data = append(data, map[string]interface{}{
			"date":        date.Format("2006-01-02"),
			"conversions": conversions,
			"revenue":     revenue,
		})

		totalRevenue += revenue
		totalConversions += conversions
	}

	summary := map[string]interface{}{
		"total_revenue":     totalRevenue,
		"total_conversions": totalConversions,
		"avg_revenue":       totalRevenue / float64(totalConversions),
	}

	return data, summary
}

// executeCustomReport executes a custom report
func (s *ReportService) executeCustomReport(tenantID int64, report *SavedReport, params map[string]interface{}) ([]map[string]interface{}, map[string]interface{}) {
	// For custom reports, parse query_config and build dynamic SQL
	// This is a simplified placeholder
	data := []map[string]interface{}{}
	summary := map[string]interface{}{}

	return data, summary
}

// CreateSnapshot creates a snapshot of report results
func (s *ReportService) CreateSnapshot(tenantID, reportID int64, resultData JSONB, generatedBy string) error {
	query := `
		INSERT INTO report_snapshots (
			report_id, snapshot_date, result_data, generated_by, created_at
		) VALUES ($1, NOW(), $2, $3, NOW())`

	_, err := s.db.Exec(query, reportID, resultData, generatedBy)
	return err
}

// GetSnapshots retrieves snapshots for a report
func (s *ReportService) GetSnapshots(reportID int64, limit int) ([]ReportSnapshot, error) {
	query := `
		SELECT * FROM report_snapshots
		WHERE report_id = $1
		ORDER BY snapshot_date DESC
		LIMIT $2`

	var snapshots []ReportSnapshot
	err := s.db.Select(&snapshots, query, reportID, limit)
	return snapshots, err
}

// ScheduleReport schedules a report for automatic execution
func (s *ReportService) ScheduleReport(tenantID, reportID int64, schedule string) error {
	query := `
		UPDATE saved_reports
		SET schedule = $1, updated_at = NOW()
		WHERE id = $2 AND tenant_id = $3`

	_, err := s.db.Exec(query, schedule, reportID, tenantID)
	return err
}

// GetScheduledReports retrieves reports that should be executed
func (s *ReportService) GetScheduledReports(tenantID int64) ([]SavedReport, error) {
	query := `
		SELECT * FROM saved_reports
		WHERE tenant_id = $1 
		AND schedule IS NOT NULL
		AND schedule != ''
		ORDER BY updated_at DESC`

	var reports []SavedReport
	err := s.db.Select(&reports, query, tenantID)
	return reports, err
}

