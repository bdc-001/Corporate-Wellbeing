package services

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// RealtimeService handles real-time data streaming and processing
type RealtimeService struct {
	db *sqlx.DB
}

// NewRealtimeService creates a new realtime service
func NewRealtimeService(db *sqlx.DB) *RealtimeService {
	return &RealtimeService{db: db}
}

// Event represents a real-time event
type Event struct {
	ID             int64          `db:"id" json:"id"`
	TenantID       int64          `db:"tenant_id" json:"tenant_id"`
	EventType      string         `db:"event_type" json:"event_type"`
	EventTimestamp time.Time      `db:"event_timestamp" json:"event_timestamp"`
	CustomerID     sql.NullInt64  `db:"customer_id" json:"customer_id"`
	AccountID      sql.NullInt64  `db:"account_id" json:"account_id"`
	SessionID      sql.NullString `db:"session_id" json:"session_id"`
	EventData      JSONB          `db:"event_data" json:"event_data"`
	Processed      bool           `db:"processed" json:"processed"`
	ProcessedAt    sql.NullTime   `db:"processed_at" json:"processed_at"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
}

// Alert represents a system alert
type Alert struct {
	ID              int64          `db:"id" json:"id"`
	TenantID        int64          `db:"tenant_id" json:"tenant_id"`
	AlertType       string         `db:"alert_type" json:"alert_type"`
	Severity        string         `db:"severity" json:"severity"`
	Title           string         `db:"title" json:"title"`
	Description     sql.NullString `db:"description" json:"description"`
	EntityType      sql.NullString `db:"entity_type" json:"entity_type"`
	EntityID        sql.NullInt64  `db:"entity_id" json:"entity_id"`
	TriggeredAt     time.Time      `db:"triggered_at" json:"triggered_at"`
	Acknowledged    bool           `db:"acknowledged" json:"acknowledged"`
	AcknowledgedBy  sql.NullString `db:"acknowledged_by" json:"acknowledged_by"`
	AcknowledgedAt  sql.NullTime   `db:"acknowledged_at" json:"acknowledged_at"`
	Resolved        bool           `db:"resolved" json:"resolved"`
	ResolvedAt      sql.NullTime   `db:"resolved_at" json:"resolved_at"`
	Metadata        JSONB          `db:"metadata" json:"metadata,omitempty"`
}

// IngestEvent ingests a new event into the stream
func (s *RealtimeService) IngestEvent(tenantID int64, event *Event) error {
	event.TenantID = tenantID
	event.CreatedAt = time.Now()
	event.Processed = false

	query := `
		INSERT INTO event_stream (
			tenant_id, event_type, event_timestamp, customer_id,
			account_id, session_id, event_data, processed, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	err := s.db.QueryRow(
		query,
		event.TenantID, event.EventType, event.EventTimestamp,
		event.CustomerID, event.AccountID, event.SessionID,
		event.EventData, event.Processed, event.CreatedAt,
	).Scan(&event.ID)

	return err
}

// GetUnprocessedEvents retrieves events that haven't been processed
func (s *RealtimeService) GetUnprocessedEvents(tenantID int64, limit int) ([]Event, error) {
	query := `
		SELECT * FROM event_stream
		WHERE tenant_id = $1 AND processed = false
		ORDER BY event_timestamp ASC
		LIMIT $2`

	var events []Event
	err := s.db.Select(&events, query, tenantID, limit)
	return events, err
}

// MarkEventProcessed marks an event as processed
func (s *RealtimeService) MarkEventProcessed(eventID int64) error {
	query := `
		UPDATE event_stream
		SET processed = true, processed_at = NOW()
		WHERE id = $1`

	_, err := s.db.Exec(query, eventID)
	return err
}

// ProcessEvent processes a single event (business logic)
func (s *RealtimeService) ProcessEvent(event *Event) error {
	// Example processing logic
	switch event.EventType {
	case "page_view":
		// Track page view
		return s.processPageView(event)
	case "form_submit":
		// Track form submission
		return s.processFormSubmit(event)
	case "purchase":
		// Track purchase
		return s.processPurchase(event)
	}

	// Mark as processed even if no specific handler
	return s.MarkEventProcessed(event.ID)
}

// processPageView handles page view events
func (s *RealtimeService) processPageView(event *Event) error {
	// Insert into page_views table if exists
	// For now, just mark as processed
	return s.MarkEventProcessed(event.ID)
}

// processFormSubmit handles form submission events
func (s *RealtimeService) processFormSubmit(event *Event) error {
	// Create interaction or lead
	return s.MarkEventProcessed(event.ID)
}

// processPurchase handles purchase events
func (s *RealtimeService) processPurchase(event *Event) error {
	// Create conversion event
	return s.MarkEventProcessed(event.ID)
}

// CreateAlert creates a new alert
func (s *RealtimeService) CreateAlert(tenantID int64, alert *Alert) error {
	alert.TenantID = tenantID
	alert.TriggeredAt = time.Now()

	query := `
		INSERT INTO alerts (
			tenant_id, alert_type, severity, title, description,
			entity_type, entity_id, triggered_at, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	err := s.db.QueryRow(
		query,
		alert.TenantID, alert.AlertType, alert.Severity, alert.Title,
		alert.Description, alert.EntityType, alert.EntityID,
		alert.TriggeredAt, alert.Metadata,
	).Scan(&alert.ID)

	return err
}

// GetAlerts retrieves alerts with optional filters
func (s *RealtimeService) GetAlerts(tenantID int64, severity, alertType string, acknowledgedOnly, unresolvedOnly bool, limit int) ([]Alert, error) {
	query := `SELECT * FROM alerts WHERE tenant_id = $1`
	args := []interface{}{tenantID}
	argIdx := 2

	if severity != "" {
		query += ` AND severity = $` + string(rune(argIdx))
		args = append(args, severity)
		argIdx++
	}

	if alertType != "" {
		query += ` AND alert_type = $` + string(rune(argIdx))
		args = append(args, alertType)
		argIdx++
	}

	if acknowledgedOnly {
		query += ` AND acknowledged = true`
	}

	if unresolvedOnly {
		query += ` AND resolved = false`
	}

	query += ` ORDER BY triggered_at DESC`
	
	if limit > 0 {
		query += ` LIMIT $` + string(rune(argIdx))
		args = append(args, limit)
	}

	var alerts []Alert
	err := s.db.Select(&alerts, query, args...)
	return alerts, err
}

// AcknowledgeAlert marks an alert as acknowledged
func (s *RealtimeService) AcknowledgeAlert(tenantID, alertID int64, acknowledgedBy string) error {
	query := `
		UPDATE alerts
		SET acknowledged = true, acknowledged_by = $1, acknowledged_at = NOW()
		WHERE id = $2 AND tenant_id = $3`

	_, err := s.db.Exec(query, acknowledgedBy, alertID, tenantID)
	return err
}

// ResolveAlert marks an alert as resolved
func (s *RealtimeService) ResolveAlert(tenantID, alertID int64) error {
	query := `
		UPDATE alerts
		SET resolved = true, resolved_at = NOW()
		WHERE id = $1 AND tenant_id = $2`

	_, err := s.db.Exec(query, alertID, tenantID)
	return err
}

// GetRealtimeMetrics provides real-time dashboard metrics
func (s *RealtimeService) GetRealtimeMetrics(tenantID int64, timeWindow int) (map[string]interface{}, error) {
	since := time.Now().Add(-time.Duration(timeWindow) * time.Minute)

	metrics := make(map[string]interface{})

	// Count recent events
	var eventCount int
	err := s.db.Get(&eventCount, `
		SELECT COUNT(*) FROM event_stream 
		WHERE tenant_id = $1 AND event_timestamp >= $2`,
		tenantID, since)
	if err != nil {
		return nil, err
	}
	metrics["event_count"] = eventCount

	// Count active sessions
	var sessionCount int
	err = s.db.Get(&sessionCount, `
		SELECT COUNT(DISTINCT session_id) FROM event_stream 
		WHERE tenant_id = $1 AND event_timestamp >= $2`,
		tenantID, since)
	if err == nil {
		metrics["active_sessions"] = sessionCount
	}

	// Count unacknowledged critical alerts
	var criticalAlerts int
	err = s.db.Get(&criticalAlerts, `
		SELECT COUNT(*) FROM alerts 
		WHERE tenant_id = $1 AND severity = 'critical' AND acknowledged = false`,
		tenantID)
	if err == nil {
		metrics["critical_alerts"] = criticalAlerts
	}

	return metrics, nil
}

