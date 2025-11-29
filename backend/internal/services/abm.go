package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// ABMService handles Account-Based Marketing features
type ABMService struct {
	db *sqlx.DB
}

// NewABMService creates a new ABM service
func NewABMService(db *sqlx.DB) *ABMService {
	return &ABMService{db: db}
}

// Account represents a B2B account
type Account struct {
	ID              int64           `db:"id" json:"id"`
	TenantID        int64           `db:"tenant_id" json:"tenant_id"`
	Name            string          `db:"name" json:"name"`
	Domain          sql.NullString  `db:"domain" json:"domain"`
	Industry        sql.NullString  `db:"industry" json:"industry"`
	CompanySize     sql.NullString  `db:"company_size" json:"company_size"`
	AnnualRevenue   sql.NullFloat64 `db:"annual_revenue" json:"annual_revenue"`
	Location        sql.NullString  `db:"location" json:"location"`
	TargetAccount   bool            `db:"target_account" json:"target_account"`
	AccountTier     sql.NullString  `db:"account_tier" json:"account_tier"`
	HealthScore     sql.NullFloat64 `db:"health_score" json:"health_score"`
	EngagementScore sql.NullFloat64 `db:"engagement_score" json:"engagement_score"`
	IntentScore     sql.NullFloat64 `db:"intent_score" json:"intent_score"`
	LifecycleStage  sql.NullString  `db:"lifecycle_stage" json:"lifecycle_stage"`
	CRMAccountID    sql.NullString  `db:"crm_account_id" json:"crm_account_id"`
	OwnerID         sql.NullString  `db:"owner_id" json:"owner_id"`
	CreatedAt       time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at" json:"updated_at"`
	Metadata        JSONB           `db:"metadata" json:"metadata,omitempty"`
}

// AccountEngagement represents an account engagement event
type AccountEngagement struct {
	ID             int64           `db:"id" json:"id"`
	TenantID       int64           `db:"tenant_id" json:"tenant_id"`
	AccountID      int64           `db:"account_id" json:"account_id"`
	EngagementType string          `db:"engagement_type" json:"engagement_type"`
	EngagementDate time.Time       `db:"engagement_date" json:"engagement_date"`
	TouchpointID   sql.NullInt64   `db:"touchpoint_id" json:"touchpoint_id"`
	ScoreImpact    sql.NullFloat64 `db:"score_impact" json:"score_impact"`
	CreatedAt      time.Time       `db:"created_at" json:"created_at"`
	Metadata       JSONB           `db:"metadata" json:"metadata,omitempty"`
}

// AccountSummary provides aggregate metrics for an account
type AccountSummary struct {
	Account
	ContactsCount          int       `json:"contacts_count"`
	EngagementEventsCount  int       `json:"engagement_events_count"`
	SessionsCount          int       `json:"sessions_count"`
	TotalRevenue           float64   `json:"total_revenue"`
	LastActivityDate       time.Time `json:"last_activity_date"`
	RecentEngagements      []AccountEngagement `json:"recent_engagements,omitempty"`
}

// CreateAccount creates a new account
func (s *ABMService) CreateAccount(tenantID int64, account *Account) error {
	account.TenantID = tenantID
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	query := `
		INSERT INTO accounts (
			tenant_id, name, domain, industry, company_size, annual_revenue,
			location, target_account, account_tier, health_score, engagement_score,
			intent_score, lifecycle_stage, crm_account_id, owner_id, metadata,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18
		) RETURNING id`

	err := s.db.QueryRow(
		query,
		account.TenantID, account.Name, account.Domain, account.Industry,
		account.CompanySize, account.AnnualRevenue, account.Location,
		account.TargetAccount, account.AccountTier, account.HealthScore,
		account.EngagementScore, account.IntentScore, account.LifecycleStage,
		account.CRMAccountID, account.OwnerID, account.Metadata,
		account.CreatedAt, account.UpdatedAt,
	).Scan(&account.ID)

	return err
}

// GetAccount retrieves an account by ID
func (s *ABMService) GetAccount(tenantID, accountID int64) (*Account, error) {
	var account Account
	query := `SELECT * FROM accounts WHERE id = $1 AND tenant_id = $2`
	err := s.db.Get(&account, query, accountID, tenantID)
	return &account, err
}

// ListAccounts lists accounts with optional filters
func (s *ABMService) ListAccounts(tenantID int64, targetOnly bool, tier, lifecycle string, limit, offset int) ([]Account, error) {
	query := `SELECT * FROM accounts WHERE tenant_id = $1`
	args := []interface{}{tenantID}
	argIdx := 2

	if targetOnly {
		query += fmt.Sprintf(" AND target_account = true")
	}
	if tier != "" {
		query += fmt.Sprintf(" AND account_tier = $%d", argIdx)
		args = append(args, tier)
		argIdx++
	}
	if lifecycle != "" {
		query += fmt.Sprintf(" AND lifecycle_stage = $%d", argIdx)
		args = append(args, lifecycle)
		argIdx++
	}

	query += " ORDER BY engagement_score DESC NULLS LAST"
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, limit)
		argIdx++
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, offset)
	}

	var accounts []Account
	err := s.db.Select(&accounts, query, args...)
	return accounts, err
}

// UpdateAccountScores updates health, engagement, and intent scores
func (s *ABMService) UpdateAccountScores(tenantID, accountID int64, healthScore, engagementScore, intentScore float64) error {
	query := `
		UPDATE accounts 
		SET health_score = $1, engagement_score = $2, intent_score = $3, updated_at = NOW()
		WHERE id = $4 AND tenant_id = $5`
	
	_, err := s.db.Exec(query, healthScore, engagementScore, intentScore, accountID, tenantID)
	return err
}

// TrackAccountEngagement records an engagement event for an account
func (s *ABMService) TrackAccountEngagement(tenantID int64, engagement *AccountEngagement) error {
	engagement.TenantID = tenantID
	engagement.CreatedAt = time.Now()

	query := `
		INSERT INTO account_engagements (
			tenant_id, account_id, engagement_type, engagement_date,
			touchpoint_id, score_impact, metadata, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	err := s.db.QueryRow(
		query,
		engagement.TenantID, engagement.AccountID, engagement.EngagementType,
		engagement.EngagementDate, engagement.TouchpointID, engagement.ScoreImpact,
		engagement.Metadata, engagement.CreatedAt,
	).Scan(&engagement.ID)

	return err
}

// GetAccountSummary retrieves comprehensive account summary with metrics
func (s *ABMService) GetAccountSummary(tenantID, accountID int64) (*AccountSummary, error) {
	query := `
		SELECT 
			a.*,
			COUNT(DISTINCT c.id) as contacts_count,
			COUNT(DISTINCT ae.id) as engagement_events_count,
			COUNT(DISTINCT s.session_id) as sessions_count,
			COALESCE(SUM(ce.revenue), 0) as total_revenue,
			MAX(s.session_start) as last_activity_date
		FROM accounts a
		LEFT JOIN customers c ON a.id = c.account_id
		LEFT JOIN account_engagements ae ON a.id = ae.account_id
		LEFT JOIN sessions s ON a.id = s.account_id
		LEFT JOIN conversion_events ce ON c.id = ce.customer_id
		WHERE a.id = $1 AND a.tenant_id = $2
		GROUP BY a.id`

	var summary AccountSummary
	err := s.db.Get(&summary, query, accountID, tenantID)
	if err != nil {
		return nil, err
	}

	// Get recent engagements
	engQuery := `
		SELECT * FROM account_engagements 
		WHERE account_id = $1 AND tenant_id = $2 
		ORDER BY engagement_date DESC 
		LIMIT 10`
	
	err = s.db.Select(&summary.RecentEngagements, engQuery, accountID, tenantID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &summary, nil
}

// GetTargetAccountInsights provides insights for target accounts
func (s *ABMService) GetTargetAccountInsights(tenantID int64, tier string) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_accounts,
			COUNT(CASE WHEN lifecycle_stage = 'Target' THEN 1 END) as target_stage,
			COUNT(CASE WHEN lifecycle_stage = 'Engaged' THEN 1 END) as engaged_stage,
			COUNT(CASE WHEN lifecycle_stage = 'MQL' THEN 1 END) as mql_stage,
			COUNT(CASE WHEN lifecycle_stage = 'SQL' THEN 1 END) as sql_stage,
			COUNT(CASE WHEN lifecycle_stage = 'Opportunity' THEN 1 END) as opportunity_stage,
			COUNT(CASE WHEN lifecycle_stage = 'Customer' THEN 1 END) as customer_stage,
			AVG(health_score) as avg_health_score,
			AVG(engagement_score) as avg_engagement_score,
			AVG(intent_score) as avg_intent_score
		FROM accounts
		WHERE tenant_id = $1 AND target_account = true`

	args := []interface{}{tenantID}
	if tier != "" {
		query += " AND account_tier = $2"
		args = append(args, tier)
	}

	var result struct {
		TotalAccounts    int     `db:"total_accounts"`
		TargetStage      int     `db:"target_stage"`
		EngagedStage     int     `db:"engaged_stage"`
		MQLStage         int     `db:"mql_stage"`
		SQLStage         int     `db:"sql_stage"`
		OpportunityStage int     `db:"opportunity_stage"`
		CustomerStage    int     `db:"customer_stage"`
		AvgHealthScore   float64 `db:"avg_health_score"`
		AvgEngagementScore float64 `db:"avg_engagement_score"`
		AvgIntentScore   float64 `db:"avg_intent_score"`
	}

	err := s.db.Get(&result, query, args...)
	if err != nil {
		return nil, err
	}

	insights := map[string]interface{}{
		"total_accounts": result.TotalAccounts,
		"stage_breakdown": map[string]int{
			"target":      result.TargetStage,
			"engaged":     result.EngagedStage,
			"mql":         result.MQLStage,
			"sql":         result.SQLStage,
			"opportunity": result.OpportunityStage,
			"customer":    result.CustomerStage,
		},
		"average_scores": map[string]float64{
			"health":     result.AvgHealthScore,
			"engagement": result.AvgEngagementScore,
			"intent":     result.AvgIntentScore,
		},
	}

	return insights, nil
}

