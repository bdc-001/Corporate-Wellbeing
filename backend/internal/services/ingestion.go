package services

import (
	"fmt"
	"time"

	"github.com/convin/crae/internal/models"
	"github.com/jmoiron/sqlx"
)

type IngestionService struct {
	db            *sqlx.DB
	identitySvc   *IdentityService
}

func NewIngestionService(db *sqlx.DB, identitySvc *IdentityService) *IngestionService {
	return &IngestionService{
		db:          db,
		identitySvc: identitySvc,
	}
}

// GetDB returns the database connection (for use in handlers)
func (s *IngestionService) GetDB() *sqlx.DB {
	return s.db
}

// IngestInteractionRequest represents an interaction ingestion request
type IngestInteractionRequest struct {
	ExternalInteractionID string                            `json:"external_interaction_id"`
	Channel               string                            `json:"channel"`
	VendorCode            *string                           `json:"vendor_code"`
	CustomerIdentifiers  []models.CustomerIdentifier        `json:"customer_identifiers"`
	StartedAt             time.Time                         `json:"started_at"`
	EndedAt               *time.Time                        `json:"ended_at"`
	Direction             string                            `json:"direction"`
	Language              string                            `json:"language"`
	Participants          []InteractionParticipantRequest  `json:"participants"`
	TranscriptURL         string                            `json:"transcript_url"`
	PrimaryIntent         string                            `json:"primary_intent"`
	SecondaryIntents      []string                          `json:"secondary_intents"`
	OutcomePrediction     string                            `json:"outcome_prediction"`
	PurchaseProbability   *float64                          `json:"purchase_probability"`
	RawMetadata           map[string]interface{}            `json:"raw_metadata"`
}

type InteractionParticipantRequest struct {
	ParticipantType string                 `json:"participant_type"`
	ExternalAgentID *string                `json:"external_agent_id"`
	Role            string                 `json:"role"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// IngestInteraction ingests an interaction and returns the created interaction ID
func (s *IngestionService) IngestInteraction(tenantID int64, req IngestInteractionRequest) (*IngestInteractionResponse, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Find or create customer
	var customerID *int64
	if len(req.CustomerIdentifiers) > 0 {
		customer, err := s.identitySvc.FindOrCreateCustomer(tenantID, req.CustomerIdentifiers)
		if err != nil {
			return nil, fmt.Errorf("failed to find/create customer: %w", err)
		}
		customerID = &customer.ID
	}

	// Get channel ID
	var channelID int
	err = tx.Get(&channelID, `SELECT id FROM channels WHERE name = $1`, req.Channel)
	if err != nil {
		return nil, fmt.Errorf("channel not found: %s", req.Channel)
	}

	// Get vendor ID if vendor code provided
	var vendorID *int
	if req.VendorCode != nil {
		err = tx.Get(&vendorID, `SELECT id FROM vendors WHERE tenant_id = $1 AND code = $2`, tenantID, *req.VendorCode)
		if err != nil {
			// Vendor not found, but continue without vendor
			vendorID = nil
		}
	}

	// Calculate duration
	var durationSeconds *int
	if req.EndedAt != nil {
		dur := int(req.EndedAt.Sub(req.StartedAt).Seconds())
		durationSeconds = &dur
	}

	// Convert secondary intents to JSONB
	var secondaryIntentsJSON models.JSONB
	if len(req.SecondaryIntents) > 0 {
		secondaryIntentsJSON = make(models.JSONB)
		secondaryIntentsJSON["intents"] = req.SecondaryIntents
	}

	// Convert raw metadata to JSONB
	var rawMetadataJSON models.JSONB
	if req.RawMetadata != nil {
		rawMetadataJSON = models.JSONB(req.RawMetadata)
	}

	// Insert interaction
	var interaction models.Interaction
	err = tx.QueryRowx(
		`INSERT INTO interactions (
			tenant_id, customer_id, external_interaction_id, channel_id, vendor_id,
			started_at, ended_at, duration_seconds, direction, language,
			transcript_location, primary_intent, secondary_intents,
			outcome_prediction, purchase_probability, raw_metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, created_at, updated_at`,
		tenantID, customerID, req.ExternalInteractionID, channelID, vendorID,
		req.StartedAt, req.EndedAt, durationSeconds, req.Direction, req.Language,
		req.TranscriptURL, req.PrimaryIntent, secondaryIntentsJSON,
		req.OutcomePrediction, req.PurchaseProbability, rawMetadataJSON,
	).Scan(&interaction.ID, &interaction.CreatedAt, &interaction.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert interaction: %w", err)
	}

	// Insert participants
	for _, partReq := range req.Participants {
		var agentID *int
		if partReq.ExternalAgentID != nil {
			err = tx.Get(&agentID, `SELECT id FROM agents WHERE external_agent_id = $1`, *partReq.ExternalAgentID)
			if err != nil {
				// Agent not found, continue without agent ID
				agentID = nil
			}
		}

		var metadataJSON models.JSONB
		if partReq.Metadata != nil {
			metadataJSON = models.JSONB(partReq.Metadata)
		}

		_, err = tx.Exec(
			`INSERT INTO interaction_participants (
				interaction_id, participant_type, agent_id, role, metadata
			) VALUES ($1, $2, $3, $4, $5)`,
			interaction.ID, partReq.ParticipantType, agentID, partReq.Role, metadataJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to insert participant: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &IngestInteractionResponse{
		InteractionID: interaction.ID,
		CustomerID:    customerID,
	}, nil
}

type IngestInteractionResponse struct {
	InteractionID int64   `json:"interaction_id"`
	CustomerID    *int64  `json:"customer_id"`
}

// IngestConversionRequest represents a conversion event ingestion request
type IngestConversionRequest struct {
	EventSource         string                     `json:"event_source"`
	ExternalEventID     string                     `json:"external_event_id"`
	CustomerIdentifiers []models.CustomerIdentifier `json:"customer_identifiers"`
	EventType           string                     `json:"event_type"`
	ProductExternalID   *string                    `json:"product_external_id"`
	Currency            string                     `json:"currency"`
	AmountDecimal       float64                    `json:"amount_decimal"`
	OccurredAt          time.Time                  `json:"occurred_at"`
	RawPayload          map[string]interface{}     `json:"raw_payload"`
}

// IngestConversion ingests a conversion event
func (s *IngestionService) IngestConversion(tenantID int64, req IngestConversionRequest) (*IngestConversionResponse, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Find or create customer
	customer, err := s.identitySvc.FindOrCreateCustomer(tenantID, req.CustomerIdentifiers)
	if err != nil {
		return nil, fmt.Errorf("failed to find/create customer: %w", err)
	}

	// Get event source ID
	var eventSourceID int
	err = tx.Get(&eventSourceID, `SELECT id FROM event_sources WHERE name = $1`, req.EventSource)
	if err != nil {
		return nil, fmt.Errorf("event source not found: %s", req.EventSource)
	}

	// Get currency ID
	var currencyID int
	err = tx.Get(&currencyID, `SELECT id FROM currencies WHERE code = $1`, req.Currency)
	if err != nil {
		return nil, fmt.Errorf("currency not found: %s", req.Currency)
	}

	// Get product ID if provided
	var productID *int
	if req.ProductExternalID != nil {
		err = tx.Get(&productID, `SELECT id FROM products WHERE external_product_id = $1`, *req.ProductExternalID)
		if err != nil {
			// Product not found, continue without product
			productID = nil
		}
	}

	// Convert raw payload to JSONB
	var rawPayloadJSON models.JSONB
	if req.RawPayload != nil {
		rawPayloadJSON = models.JSONB(req.RawPayload)
	}

	// Insert conversion event
	var conversion models.ConversionEvent
	err = tx.QueryRowx(
		`INSERT INTO conversion_events (
			tenant_id, customer_id, event_source_id, external_event_id,
			event_type, product_id, currency_id, amount_decimal, occurred_at, raw_payload
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at`,
		tenantID, customer.ID, eventSourceID, req.ExternalEventID,
		req.EventType, productID, currencyID, req.AmountDecimal, req.OccurredAt, rawPayloadJSON,
	).Scan(&conversion.ID, &conversion.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert conversion event: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &IngestConversionResponse{
		ConversionEventID: conversion.ID,
		CustomerID:        customer.ID,
	}, nil
}

type IngestConversionResponse struct {
	ConversionEventID int64 `json:"conversion_event_id"`
	CustomerID        int64 `json:"customer_id"`
}

