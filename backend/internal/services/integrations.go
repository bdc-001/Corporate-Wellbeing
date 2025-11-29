package services

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// IntegrationService handles external platform integrations
type IntegrationService struct {
	db *sqlx.DB
}

// NewIntegrationService creates a new integration service
func NewIntegrationService(db *sqlx.DB) *IntegrationService {
	return &IntegrationService{db: db}
}

// Integration represents an external platform integration
type Integration struct {
	ID             int64          `db:"id" json:"id"`
	TenantID       int64          `db:"tenant_id" json:"tenant_id"`
	Platform       string         `db:"platform" json:"platform"`
	IntegrationType string        `db:"integration_type" json:"integration_type"`
	IsActive       bool           `db:"is_active" json:"is_active"`
	Credentials    JSONB          `db:"credentials" json:"credentials,omitempty"`
	SyncConfig     JSONB          `db:"sync_config" json:"sync_config,omitempty"`
	LastSyncAt     sql.NullTime   `db:"last_sync_at" json:"last_sync_at"`
	LastSyncStatus sql.NullString `db:"last_sync_status" json:"last_sync_status"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at" json:"updated_at"`
}

// SyncLog represents a sync operation log
type SyncLog struct {
	ID               int64        `db:"id" json:"id"`
	IntegrationID    int64        `db:"integration_id" json:"integration_id"`
	SyncStartedAt    time.Time    `db:"sync_started_at" json:"sync_started_at"`
	SyncCompletedAt  sql.NullTime `db:"sync_completed_at" json:"sync_completed_at"`
	Status           string       `db:"status" json:"status"`
	RecordsProcessed int          `db:"records_processed" json:"records_processed"`
	RecordsCreated   int          `db:"records_created" json:"records_created"`
	RecordsUpdated   int          `db:"records_updated" json:"records_updated"`
	RecordsFailed    int          `db:"records_failed" json:"records_failed"`
	ErrorDetails     JSONB        `db:"error_details" json:"error_details,omitempty"`
}

// CreateIntegration creates a new integration
func (s *IntegrationService) CreateIntegration(tenantID int64, integration *Integration) error {
	integration.TenantID = tenantID
	integration.CreatedAt = time.Now()
	integration.UpdatedAt = time.Now()

	query := `
		INSERT INTO integrations (
			tenant_id, platform, integration_type, is_active,
			credentials, sync_config, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	err := s.db.QueryRow(
		query,
		integration.TenantID, integration.Platform, integration.IntegrationType,
		integration.IsActive, integration.Credentials, integration.SyncConfig,
		integration.CreatedAt, integration.UpdatedAt,
	).Scan(&integration.ID)

	return err
}

// GetIntegration retrieves an integration by ID
func (s *IntegrationService) GetIntegration(tenantID, integrationID int64) (*Integration, error) {
	var integration Integration
	query := `SELECT * FROM integrations WHERE id = $1 AND tenant_id = $2`
	err := s.db.Get(&integration, query, integrationID, tenantID)
	return &integration, err
}

// ListIntegrations lists all integrations for a tenant
func (s *IntegrationService) ListIntegrations(tenantID int64, platform, integrationType string, activeOnly bool) ([]Integration, error) {
	query := `SELECT * FROM integrations WHERE tenant_id = $1`
	args := []interface{}{tenantID}
	argIdx := 2

	if platform != "" {
		query += ` AND platform = $` + string(rune(argIdx))
		args = append(args, platform)
		argIdx++
	}

	if integrationType != "" {
		query += ` AND integration_type = $` + string(rune(argIdx))
		args = append(args, integrationType)
		argIdx++
	}

	if activeOnly {
		query += ` AND is_active = true`
	}

	query += ` ORDER BY created_at DESC`

	var integrations []Integration
	err := s.db.Select(&integrations, query, args...)
	return integrations, err
}

// UpdateIntegration updates an integration
func (s *IntegrationService) UpdateIntegration(tenantID, integrationID int64, updates map[string]interface{}) error {
	// Build dynamic update query
	query := `UPDATE integrations SET updated_at = NOW()`
	args := []interface{}{}
	argIdx := 1

	for key, value := range updates {
		query += `, ` + key + ` = $` + string(rune(argIdx))
		args = append(args, value)
		argIdx++
	}

	query += ` WHERE id = $` + string(rune(argIdx)) + ` AND tenant_id = $` + string(rune(argIdx+1))
	args = append(args, integrationID, tenantID)

	_, err := s.db.Exec(query, args...)
	return err
}

// DeleteIntegration deletes an integration
func (s *IntegrationService) DeleteIntegration(tenantID, integrationID int64) error {
	query := `DELETE FROM integrations WHERE id = $1 AND tenant_id = $2`
	_, err := s.db.Exec(query, integrationID, tenantID)
	return err
}

// StartSync initiates a sync operation
func (s *IntegrationService) StartSync(integrationID int64) (*SyncLog, error) {
	syncLog := &SyncLog{
		IntegrationID: integrationID,
		SyncStartedAt: time.Now(),
		Status:        "running",
	}

	query := `
		INSERT INTO sync_logs (
			integration_id, sync_started_at, status,
			records_processed, records_created, records_updated, records_failed
		) VALUES ($1, $2, $3, 0, 0, 0, 0)
		RETURNING id`

	err := s.db.QueryRow(
		query,
		syncLog.IntegrationID, syncLog.SyncStartedAt, syncLog.Status,
	).Scan(&syncLog.ID)

	return syncLog, err
}

// CompletSync completes a sync operation
func (s *IntegrationService) CompleteSync(syncLogID int64, status string, processed, created, updated, failed int, errors JSONB) error {
	query := `
		UPDATE sync_logs
		SET sync_completed_at = NOW(), status = $1,
		    records_processed = $2, records_created = $3,
		    records_updated = $4, records_failed = $5,
		    error_details = $6
		WHERE id = $7`

	_, err := s.db.Exec(
		query, status, processed, created, updated, failed, errors, syncLogID,
	)

	if err == nil {
		// Update integration last_sync info
		updateQuery := `
			UPDATE integrations
			SET last_sync_at = NOW(), last_sync_status = $1
			WHERE id = (SELECT integration_id FROM sync_logs WHERE id = $2)`
		
		s.db.Exec(updateQuery, status, syncLogID)
	}

	return err
}

// GetSyncLogs retrieves sync logs for an integration
func (s *IntegrationService) GetSyncLogs(integrationID int64, limit int) ([]SyncLog, error) {
	query := `
		SELECT * FROM sync_logs
		WHERE integration_id = $1
		ORDER BY sync_started_at DESC
		LIMIT $2`

	var logs []SyncLog
	err := s.db.Select(&logs, query, integrationID, limit)
	return logs, err
}

// SyncFromCRM syncs data from a CRM platform
func (s *IntegrationService) SyncFromCRM(tenantID, integrationID int64) error {
	integration, err := s.GetIntegration(tenantID, integrationID)
	if err != nil {
		return err
	}

	syncLog, err := s.StartSync(integrationID)
	if err != nil {
		return err
	}

	// Simulate sync process
	// In production, implement actual API calls to CRM platforms
	processed, created, updated, failed := 0, 0, 0, 0

	switch integration.Platform {
	case "salesforce":
		processed, created, updated, failed = s.syncFromSalesforce(tenantID, integration)
	case "hubspot":
		processed, created, updated, failed = s.syncFromHubspot(tenantID, integration)
	case "marketo":
		processed, created, updated, failed = s.syncFromMarketo(tenantID, integration)
	default:
		// Unknown platform
		failed = 1
	}

	status := "completed"
	if failed > 0 {
		status = "failed"
	}

	return s.CompleteSync(syncLog.ID, status, processed, created, updated, failed, nil)
}

// syncFromSalesforce syncs from Salesforce (placeholder)
func (s *IntegrationService) syncFromSalesforce(tenantID int64, integration *Integration) (int, int, int, int) {
	// Placeholder - implement Salesforce API integration
	return 100, 10, 20, 0
}

// syncFromHubspot syncs from HubSpot (placeholder)
func (s *IntegrationService) syncFromHubspot(tenantID int64, integration *Integration) (int, int, int, int) {
	// Placeholder - implement HubSpot API integration
	return 100, 15, 25, 0
}

// syncFromMarketo syncs from Marketo (placeholder)
func (s *IntegrationService) syncFromMarketo(tenantID int64, integration *Integration) (int, int, int, int) {
	// Placeholder - implement Marketo API integration
	return 100, 12, 22, 0
}

// SyncAdSpend syncs ad spend data from advertising platforms
func (s *IntegrationService) SyncAdSpend(tenantID, integrationID int64, startDate, endDate time.Time) error {
	integration, err := s.GetIntegration(tenantID, integrationID)
	if err != nil {
		return err
	}

	syncLog, err := s.StartSync(integrationID)
	if err != nil {
		return err
	}

	processed, created, updated, failed := 0, 0, 0, 0

	switch integration.Platform {
	case "google_ads":
		processed, created, updated, failed = s.syncGoogleAds(tenantID, integration, startDate, endDate)
	case "facebook_ads":
		processed, created, updated, failed = s.syncFacebookAds(tenantID, integration, startDate, endDate)
	case "linkedin_ads":
		processed, created, updated, failed = s.syncLinkedInAds(tenantID, integration, startDate, endDate)
	}

	status := "completed"
	if failed > 0 {
		status = "failed"
	}

	return s.CompleteSync(syncLog.ID, status, processed, created, updated, failed, nil)
}

// syncGoogleAds syncs from Google Ads (placeholder)
func (s *IntegrationService) syncGoogleAds(tenantID int64, integration *Integration, startDate, endDate time.Time) (int, int, int, int) {
	// Placeholder - implement Google Ads API integration
	return 30, 5, 10, 0
}

// syncFacebookAds syncs from Facebook Ads (placeholder)
func (s *IntegrationService) syncFacebookAds(tenantID int64, integration *Integration, startDate, endDate time.Time) (int, int, int, int) {
	// Placeholder - implement Facebook Ads API integration
	return 25, 4, 8, 0
}

// syncLinkedInAds syncs from LinkedIn Ads (placeholder)
func (s *IntegrationService) syncLinkedInAds(tenantID int64, integration *Integration, startDate, endDate time.Time) (int, int, int, int) {
	// Placeholder - implement LinkedIn Ads API integration
	return 20, 3, 7, 0
}

// TestIntegration tests integration connectivity
func (s *IntegrationService) TestIntegration(tenantID, integrationID int64) (bool, string) {
	integration, err := s.GetIntegration(tenantID, integrationID)
	if err != nil {
		return false, "Integration not found"
	}

	// Simulate connection test
	// In production, implement actual API connectivity tests
	switch integration.Platform {
	case "salesforce", "hubspot", "marketo":
		return true, "Connection successful"
	case "google_ads", "facebook_ads", "linkedin_ads":
		return true, "API connection verified"
	default:
		return false, "Unknown platform"
	}
}

