package services

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// BehaviorService handles user behavior analytics
type BehaviorService struct {
	db *sqlx.DB
}

// NewBehaviorService creates a new behavior service
func NewBehaviorService(db *sqlx.DB) *BehaviorService {
	return &BehaviorService{db: db}
}

// PageView represents a page view event
type PageView struct {
	ID            int64          `db:"id" json:"id"`
	TenantID      int64          `db:"tenant_id" json:"tenant_id"`
	SessionID     string         `db:"session_id" json:"session_id"`
	CustomerID    sql.NullInt64  `db:"customer_id" json:"customer_id"`
	PageURL       string         `db:"page_url" json:"page_url"`
	PageTitle     sql.NullString `db:"page_title" json:"page_title"`
	Referrer      sql.NullString `db:"referrer" json:"referrer"`
	ViewTimestamp time.Time      `db:"view_timestamp" json:"view_timestamp"`
	TimeOnPage    sql.NullInt64  `db:"time_on_page" json:"time_on_page"`
	ScrollDepth   sql.NullInt64  `db:"scroll_depth" json:"scroll_depth"`
	ExitPage      bool           `db:"exit_page" json:"exit_page"`
	DeviceType    sql.NullString `db:"device_type" json:"device_type"`
	Browser       sql.NullString `db:"browser" json:"browser"`
	OS            sql.NullString `db:"os" json:"os"`
	Location      JSONB          `db:"location" json:"location,omitempty"`
	CreatedAt     time.Time      `db:"created_at" json:"created_at"`
}

// UserAction represents a user action/event
type UserAction struct {
	ID              int64          `db:"id" json:"id"`
	TenantID        int64          `db:"tenant_id" json:"tenant_id"`
	SessionID       string         `db:"session_id" json:"session_id"`
	CustomerID      sql.NullInt64  `db:"customer_id" json:"customer_id"`
	ActionType      string         `db:"action_type" json:"action_type"`
	ActionTarget    sql.NullString `db:"action_target" json:"action_target"`
	ActionTimestamp time.Time      `db:"action_timestamp" json:"action_timestamp"`
	PageURL         sql.NullString `db:"page_url" json:"page_url"`
	ActionData      JSONB          `db:"action_data" json:"action_data,omitempty"`
	CreatedAt       time.Time      `db:"created_at" json:"created_at"`
}

// Session represents a user session
type Session struct {
	ID               int64           `db:"id" json:"id"`
	TenantID         int64           `db:"tenant_id" json:"tenant_id"`
	SessionID        string          `db:"session_id" json:"session_id"`
	CustomerID       sql.NullInt64   `db:"customer_id" json:"customer_id"`
	AccountID        sql.NullInt64   `db:"account_id" json:"account_id"`
	SessionStart     time.Time       `db:"session_start" json:"session_start"`
	SessionEnd       sql.NullTime    `db:"session_end" json:"session_end"`
	Duration         sql.NullInt64   `db:"duration" json:"duration"`
	PageViewsCount   int             `db:"page_views_count" json:"page_views_count"`
	ActionsCount     int             `db:"actions_count" json:"actions_count"`
	EntryPage        sql.NullString  `db:"entry_page" json:"entry_page"`
	ExitPage         sql.NullString  `db:"exit_page" json:"exit_page"`
	UTMSource        sql.NullString  `db:"utm_source" json:"utm_source"`
	UTMMedium        sql.NullString  `db:"utm_medium" json:"utm_medium"`
	UTMCampaign      sql.NullString  `db:"utm_campaign" json:"utm_campaign"`
	UTMTerm          sql.NullString  `db:"utm_term" json:"utm_term"`
	UTMContent       sql.NullString  `db:"utm_content" json:"utm_content"`
	DeviceType       sql.NullString  `db:"device_type" json:"device_type"`
	Browser          sql.NullString  `db:"browser" json:"browser"`
	OS               sql.NullString  `db:"os" json:"os"`
	Location         JSONB           `db:"location" json:"location,omitempty"`
	Converted        bool            `db:"converted" json:"converted"`
	ConversionValue  sql.NullFloat64 `db:"conversion_value" json:"conversion_value"`
	CreatedAt        time.Time       `db:"created_at" json:"created_at"`
}

// TrackPageView tracks a page view
func (s *BehaviorService) TrackPageView(tenantID int64, pageView *PageView) error {
	pageView.TenantID = tenantID
	pageView.CreatedAt = time.Now()

	query := `
		INSERT INTO page_views (
			tenant_id, session_id, customer_id, page_url, page_title,
			referrer, view_timestamp, time_on_page, scroll_depth,
			exit_page, device_type, browser, os, location, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id`

	err := s.db.QueryRow(
		query,
		pageView.TenantID, pageView.SessionID, pageView.CustomerID,
		pageView.PageURL, pageView.PageTitle, pageView.Referrer,
		pageView.ViewTimestamp, pageView.TimeOnPage, pageView.ScrollDepth,
		pageView.ExitPage, pageView.DeviceType, pageView.Browser,
		pageView.OS, pageView.Location, pageView.CreatedAt,
	).Scan(&pageView.ID)

	return err
}

// TrackUserAction tracks a user action
func (s *BehaviorService) TrackUserAction(tenantID int64, action *UserAction) error {
	action.TenantID = tenantID
	action.CreatedAt = time.Now()

	query := `
		INSERT INTO user_actions (
			tenant_id, session_id, customer_id, action_type,
			action_target, action_timestamp, page_url, action_data, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	err := s.db.QueryRow(
		query,
		action.TenantID, action.SessionID, action.CustomerID,
		action.ActionType, action.ActionTarget, action.ActionTimestamp,
		action.PageURL, action.ActionData, action.CreatedAt,
	).Scan(&action.ID)

	return err
}

// CreateSession creates a new session
func (s *BehaviorService) CreateSession(tenantID int64, session *Session) error {
	session.TenantID = tenantID
	session.CreatedAt = time.Now()

	query := `
		INSERT INTO sessions (
			tenant_id, session_id, customer_id, account_id, session_start,
			entry_page, utm_source, utm_medium, utm_campaign, utm_term,
			utm_content, device_type, browser, os, location, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id`

	err := s.db.QueryRow(
		query,
		session.TenantID, session.SessionID, session.CustomerID,
		session.AccountID, session.SessionStart, session.EntryPage,
		session.UTMSource, session.UTMMedium, session.UTMCampaign,
		session.UTMTerm, session.UTMContent, session.DeviceType,
		session.Browser, session.OS, session.Location, session.CreatedAt,
	).Scan(&session.ID)

	return err
}

// UpdateSession updates session end time and metrics
func (s *BehaviorService) UpdateSession(tenantID int64, sessionID string, sessionEnd time.Time, exitPage string, converted bool, conversionValue float64) error {
	duration := sessionEnd.Unix() - time.Now().Unix() // Simplified

	query := `
		UPDATE sessions
		SET session_end = $1, duration = $2, exit_page = $3,
		    converted = $4, conversion_value = $5,
		    page_views_count = (SELECT COUNT(*) FROM page_views WHERE session_id = $6),
		    actions_count = (SELECT COUNT(*) FROM user_actions WHERE session_id = $6)
		WHERE session_id = $6 AND tenant_id = $7`

	_, err := s.db.Exec(
		query, sessionEnd, duration, exitPage, converted,
		conversionValue, sessionID, tenantID,
	)

	return err
}

// GetSession retrieves a session by ID
func (s *BehaviorService) GetSession(tenantID int64, sessionID string) (*Session, error) {
	var session Session
	query := `SELECT * FROM sessions WHERE session_id = $1 AND tenant_id = $2`
	err := s.db.Get(&session, query, sessionID, tenantID)
	return &session, err
}

// GetSessionDetails retrieves full session with page views and actions
func (s *BehaviorService) GetSessionDetails(tenantID int64, sessionID string) (map[string]interface{}, error) {
	session, err := s.GetSession(tenantID, sessionID)
	if err != nil {
		return nil, err
	}

	// Get page views
	var pageViews []PageView
	pvQuery := `SELECT * FROM page_views WHERE session_id = $1 AND tenant_id = $2 ORDER BY view_timestamp`
	s.db.Select(&pageViews, pvQuery, sessionID, tenantID)

	// Get actions
	var actions []UserAction
	actionQuery := `SELECT * FROM user_actions WHERE session_id = $1 AND tenant_id = $2 ORDER BY action_timestamp`
	s.db.Select(&actions, actionQuery, sessionID, tenantID)

	return map[string]interface{}{
		"session":    session,
		"page_views": pageViews,
		"actions":    actions,
	}, nil
}

// GetTopPages retrieves most visited pages
func (s *BehaviorService) GetTopPages(tenantID int64, startDate, endDate time.Time, limit int) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			page_url,
			COUNT(*) as views,
			COUNT(DISTINCT session_id) as unique_sessions,
			AVG(time_on_page) as avg_time_on_page,
			AVG(scroll_depth) as avg_scroll_depth
		FROM page_views
		WHERE tenant_id = $1 
		AND view_timestamp >= $2 
		AND view_timestamp < $3
		GROUP BY page_url
		ORDER BY views DESC
		LIMIT $4`

	rows, err := s.db.Query(query, tenantID, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		var pageURL string
		var views, uniqueSessions int
		var avgTime, avgScroll sql.NullFloat64

		rows.Scan(&pageURL, &views, &uniqueSessions, &avgTime, &avgScroll)

		results = append(results, map[string]interface{}{
			"page_url":        pageURL,
			"views":           views,
			"unique_sessions": uniqueSessions,
			"avg_time":        avgTime.Float64,
			"avg_scroll":      avgScroll.Float64,
		})
	}

	return results, nil
}

// GetUserFlowAnalysis analyzes user flow from entry to exit
func (s *BehaviorService) GetUserFlowAnalysis(tenantID int64, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			entry_page,
			exit_page,
			COUNT(*) as session_count,
			AVG(duration) as avg_duration,
			COUNT(CASE WHEN converted = true THEN 1 END) as conversions,
			AVG(CASE WHEN converted = true THEN conversion_value END) as avg_conversion_value
		FROM sessions
		WHERE tenant_id = $1 
		AND session_start >= $2 
		AND session_start < $3
		AND entry_page IS NOT NULL
		AND exit_page IS NOT NULL
		GROUP BY entry_page, exit_page
		ORDER BY session_count DESC
		LIMIT 50`

	rows, err := s.db.Query(query, tenantID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		var entryPage, exitPage string
		var sessionCount, conversions int
		var avgDuration sql.NullFloat64
		var avgConversionValue sql.NullFloat64

		rows.Scan(&entryPage, &exitPage, &sessionCount, &avgDuration, &conversions, &avgConversionValue)

		results = append(results, map[string]interface{}{
			"entry_page":           entryPage,
			"exit_page":            exitPage,
			"session_count":        sessionCount,
			"avg_duration":         avgDuration.Float64,
			"conversions":          conversions,
			"avg_conversion_value": avgConversionValue.Float64,
		})
	}

	return results, nil
}

// GetDeviceBreakdown retrieves session breakdown by device type
func (s *BehaviorService) GetDeviceBreakdown(tenantID int64, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			device_type,
			COUNT(*) as session_count,
			AVG(duration) as avg_duration,
			AVG(page_views_count) as avg_page_views,
			COUNT(CASE WHEN converted = true THEN 1 END) as conversions
		FROM sessions
		WHERE tenant_id = $1 
		AND session_start >= $2 
		AND session_start < $3
		AND device_type IS NOT NULL
		GROUP BY device_type
		ORDER BY session_count DESC`

	rows, err := s.db.Query(query, tenantID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]map[string]interface{}, 0)
	for rows.Next() {
		var deviceType string
		var sessionCount, conversions int
		var avgDuration, avgPageViews sql.NullFloat64

		rows.Scan(&deviceType, &sessionCount, &avgDuration, &avgPageViews, &conversions)

		results = append(results, map[string]interface{}{
			"device_type":    deviceType,
			"session_count":  sessionCount,
			"avg_duration":   avgDuration.Float64,
			"avg_page_views": avgPageViews.Float64,
			"conversions":    conversions,
		})
	}

	return results, nil
}

// GetConversionFunnel analyzes conversion funnel drop-offs
func (s *BehaviorService) GetConversionFunnel(tenantID int64, funnelSteps []string, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	// Simplified funnel analysis
	results := make([]map[string]interface{}, 0)

	for i, step := range funnelSteps {
		var count int
		query := `
			SELECT COUNT(DISTINCT session_id)
			FROM page_views
			WHERE tenant_id = $1 
			AND page_url LIKE $2
			AND view_timestamp >= $3 
			AND view_timestamp < $4`

		s.db.Get(&count, query, tenantID, "%"+step+"%", startDate, endDate)

		dropoff := 0.0
		if i > 0 && len(results) > 0 {
			prevCount := results[i-1]["count"].(int)
			if prevCount > 0 {
				dropoff = float64(prevCount-count) / float64(prevCount) * 100
			}
		}

		results = append(results, map[string]interface{}{
			"step":     step,
			"count":    count,
			"dropoff_%": dropoff,
		})
	}

	return results, nil
}

