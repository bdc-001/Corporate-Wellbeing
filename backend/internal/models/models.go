package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// JSONB is a helper type for PostgreSQL JSONB columns
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), j)
	}
	return json.Unmarshal(bytes, j)
}

// Tenant represents a multi-tenant organization
type Tenant struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Code      string    `db:"code" json:"code"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Currency represents a currency
type Currency struct {
	ID        int       `db:"id" json:"id"`
	Code      string    `db:"code" json:"code"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Channel represents a communication channel
type Channel struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

// EventSource represents a source system (CRM, billing, etc.)
type EventSource struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Type        string    `db:"type" json:"type"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

// AttributionModel represents an attribution model
type AttributionModel struct {
	ID          int       `db:"id" json:"id"`
	Code        string    `db:"code" json:"code"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Params      JSONB     `db:"params" json:"params"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

// Product represents a product or plan
type Product struct {
	ID               int       `db:"id" json:"id"`
	ExternalProductID string   `db:"external_product_id" json:"external_product_id"`
	Name             string    `db:"name" json:"name"`
	Category         string    `db:"category" json:"category"`
	Metadata         JSONB     `db:"metadata" json:"metadata"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
}

// Vendor represents a BPO vendor
type Vendor struct {
	ID          int       `db:"id" json:"id"`
	TenantID    int64     `db:"tenant_id" json:"tenant_id"`
	Name        string    `db:"name" json:"name"`
	Code        string    `db:"code" json:"code"`
	Description string    `db:"description" json:"description"`
	IsActive    bool      `db:"is_active" json:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// Team represents a team with hierarchical structure
type Team struct {
	ID          int       `db:"id" json:"id"`
	TenantID    int64     `db:"tenant_id" json:"tenant_id"`
	VendorID    *int      `db:"vendor_id" json:"vendor_id"` // Nullable
	GroupID     *int      `db:"group_id" json:"group_id"` // Parent team ID
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description"` // Nullable
	About       *string   `db:"about" json:"about"` // Nullable
	ManagerID   *int64    `db:"manager_id" json:"manager_id"` // References users.id
	UseCaseID   *int      `db:"use_case_id" json:"use_case_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// UseCase represents a team use case (Sales, Collection, Support, etc.)
type UseCase struct {
	ID          int       `db:"id" json:"id"`
	TenantID    int64     `db:"tenant_id" json:"tenant_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// Agent represents an agent
type Agent struct {
	ID              int       `db:"id" json:"id"`
	VendorID        int       `db:"vendor_id" json:"vendor_id"`
	TeamID          *int       `db:"team_id" json:"team_id"`
	Name            string    `db:"name" json:"name"`
	Email           string    `db:"email" json:"email"`
	ExternalAgentID string    `db:"external_agent_id" json:"external_agent_id"`
	IsActive        bool      `db:"is_active" json:"is_active"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// Customer represents a unified customer entity
type Customer struct {
	ID        int64     `db:"id" json:"id"`
	TenantID  int64     `db:"tenant_id" json:"tenant_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// CustomerIdentifier represents an identifier for a customer
type CustomerIdentifier struct {
	ID          int64     `db:"id" json:"id"`
	CustomerID  int64     `db:"customer_id" json:"customer_id"`
	Type        string    `db:"type" json:"type"`
	Value       string    `db:"value" json:"value"`
	SourceSystem string   `db:"source_system" json:"source_system"`
	IsPrimary   bool      `db:"is_primary" json:"is_primary"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

// Interaction represents a call, chat, or other interaction
type Interaction struct {
	ID                  int64     `db:"id" json:"id"`
	TenantID            int64     `db:"tenant_id" json:"tenant_id"`
	CustomerID          *int64    `db:"customer_id" json:"customer_id"`
	ExternalInteractionID string  `db:"external_interaction_id" json:"external_interaction_id"`
	ChannelID           int       `db:"channel_id" json:"channel_id"`
	VendorID            *int       `db:"vendor_id" json:"vendor_id"`
	StartedAt           time.Time `db:"started_at" json:"started_at"`
	EndedAt             *time.Time `db:"ended_at" json:"ended_at"`
	DurationSeconds     *int       `db:"duration_seconds" json:"duration_seconds"`
	Direction           string    `db:"direction" json:"direction"`
	Language            string    `db:"language" json:"language"`
	TranscriptLocation  string    `db:"transcript_location" json:"transcript_location"`
	PrimaryIntent       string    `db:"primary_intent" json:"primary_intent"`
	SecondaryIntents    JSONB     `db:"secondary_intents" json:"secondary_intents"`
	OutcomePrediction   string    `db:"outcome_prediction" json:"outcome_prediction"`
	PurchaseProbability *float64  `db:"purchase_probability" json:"purchase_probability"`
	RawMetadata         JSONB     `db:"raw_metadata" json:"raw_metadata"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
	UpdatedAt           time.Time `db:"updated_at" json:"updated_at"`
}

// InteractionParticipant represents a participant in an interaction
type InteractionParticipant struct {
	ID              int64     `db:"id" json:"id"`
	InteractionID   int64     `db:"interaction_id" json:"interaction_id"`
	ParticipantType string    `db:"participant_type" json:"participant_type"`
	AgentID         *int       `db:"agent_id" json:"agent_id"`
	Role            string    `db:"role" json:"role"`
	Metadata        JSONB     `db:"metadata" json:"metadata"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

// ConversionEvent represents a purchase, renewal, or other conversion
type ConversionEvent struct {
	ID              int64     `db:"id" json:"id"`
	TenantID        int64     `db:"tenant_id" json:"tenant_id"`
	CustomerID      int64     `db:"customer_id" json:"customer_id"`
	EventSourceID   int       `db:"event_source_id" json:"event_source_id"`
	ExternalEventID string    `db:"external_event_id" json:"external_event_id"`
	EventType       string    `db:"event_type" json:"event_type"`
	ProductID       *int       `db:"product_id" json:"product_id"`
	CurrencyID      int       `db:"currency_id" json:"currency_id"`
	AmountDecimal   float64   `db:"amount_decimal" json:"amount_decimal"`
	OccurredAt      time.Time `db:"occurred_at" json:"occurred_at"`
	RawPayload      JSONB     `db:"raw_payload" json:"raw_payload"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

// AttributionRun represents a batch attribution calculation
type AttributionRun struct {
	ID          int64     `db:"id" json:"id"`
	TenantID    int64     `db:"tenant_id" json:"tenant_id"`
	ModelID     int       `db:"model_id" json:"model_id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Config      JSONB     `db:"config" json:"config"`
	Status      string    `db:"status" json:"status"`
	StartedAt   *time.Time `db:"started_at" json:"started_at"`
	CompletedAt *time.Time `db:"completed_at" json:"completed_at"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// AttributionResult represents an attribution result
type AttributionResult struct {
	ID                 int64     `db:"id" json:"id"`
	TenantID           int64     `db:"tenant_id" json:"tenant_id"`
	AttributionRunID   int64     `db:"attribution_run_id" json:"attribution_run_id"`
	ConversionEventID  int64     `db:"conversion_event_id" json:"conversion_event_id"`
	InteractionID      int64     `db:"interaction_id" json:"interaction_id"`
	CustomerID         int64     `db:"customer_id" json:"customer_id"`
	AgentID            *int       `db:"agent_id" json:"agent_id"`
	TeamID             *int       `db:"team_id" json:"team_id"`
	VendorID           *int       `db:"vendor_id" json:"vendor_id"`
	ModelID            int       `db:"model_id" json:"model_id"`
	AttributionWeight  float64   `db:"attribution_weight" json:"attribution_weight"`
	AttributedAmount   float64   `db:"attributed_amount" json:"attributed_amount"`
	IsPrimaryTouch     bool      `db:"is_primary_touch" json:"is_primary_touch"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
}

// CustomerJourney represents a customer's complete journey
type CustomerJourney struct {
	CustomerID       int64                      `json:"customer_id"`
	Identifiers      []CustomerIdentifier       `json:"identifiers"`
	Interactions     []InteractionWithChannel   `json:"interactions"`
	ConversionEvents []ConversionEventWithSource `json:"conversion_events"`
}

type InteractionWithChannel struct {
	Interaction
	ChannelName  string                      `db:"channel_name" json:"channel_name"`
	Participants []InteractionParticipant     `json:"participants"`
}

type ConversionEventWithSource struct {
	ConversionEvent
	EventSourceName string `db:"event_source_name" json:"event_source_name"`
	CurrencyCode    string `db:"currency_code" json:"currency_code"`
}

// User represents a platform user
type User struct {
	ID                      int64      `db:"id" json:"id"`
	TenantID                int64      `db:"tenant_id" json:"tenant_id"`
	Email                   string     `db:"email" json:"email"`
	Name                    string     `db:"name" json:"name"` // Kept for backward compatibility
	FirstName               *string    `db:"first_name" json:"first_name"`
	LastName                *string    `db:"last_name" json:"last_name"`
	PasswordHash            string     `db:"password_hash" json:"-"` // Never return in JSON
	Phone                   *string    `db:"phone" json:"phone"`
	RoleID                  *int       `db:"role_id" json:"role_id"`
	ManagerID               *int64     `db:"manager_id" json:"manager_id"`
	AuditorID               *int64     `db:"auditor_id" json:"auditor_id"`
	TeamID                  *int       `db:"team_id" json:"team_id"`
	UserType                string     `db:"user_type" json:"user_type"`
	Location                *string    `db:"location" json:"location"`
	Timezone                *string    `db:"timezone" json:"timezone"`
	IsActive                bool       `db:"is_active" json:"is_active"`
	LastLoginAt             *time.Time `db:"last_login_at" json:"last_login_at"`
	PasswordResetToken      *string    `db:"password_reset_token" json:"-"`
	PasswordResetExpiresAt  *time.Time `db:"password_reset_expires_at" json:"-"`
	CreatedAt               time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt               time.Time  `db:"updated_at" json:"updated_at"`
	// Joined fields
	RoleName                *string    `json:"role_name,omitempty"`
	ManagerName             *string    `json:"manager_name,omitempty"`
	AuditorName             *string    `json:"auditor_name,omitempty"`
	TeamName                *string    `json:"team_name,omitempty"`
}

// Role represents a user role with permissions
type Role struct {
	ID          int       `db:"id" json:"id"`
	TenantID    int64     `db:"tenant_id" json:"tenant_id"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description"`
	CodeNames   []string  `db:"code_names" json:"code_names"` // Array of permission code names
	CanBeEdited bool      `db:"can_be_edited" json:"can_be_edited"`
	IsDefault   bool      `db:"is_default" json:"is_default"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	// Joined fields
	AllowedTeamIDs []int  `json:"allowed_team_ids,omitempty"`
	UserCount      *int   `json:"user_count,omitempty"`
}

// Permission represents a permission code name
type Permission struct {
	ID          int       `db:"id" json:"id"`
	CodeName    string    `db:"code_name" json:"code_name"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description"`
	GroupID     *int      `db:"group_id" json:"group_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	// Joined fields
	GroupName   *string   `json:"group_name,omitempty"`
}

// PermissionGroup represents a group of permissions
type PermissionGroup struct {
	ID           int       `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Description  *string   `db:"description" json:"description"`
	DisplayOrder int       `db:"display_order" json:"display_order"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	// Joined fields
	Permissions  []Permission `json:"permissions,omitempty"`
}

