package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/convin/crae/internal/models"
	"github.com/convin/crae/internal/services"
)

// ConvinWebhookPayload represents a webhook payload from Convin
type ConvinWebhookPayload struct {
	EventType string                 `json:"event_type"`
	CallID    string                 `json:"call_id"`
	TenantID  int64                  `json:"tenant_id"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// HandleConvinWebhook handles webhooks from Convin for live call flow
func (h *Handlers) HandleConvinWebhook(c *gin.Context) {
	var payload ConvinWebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload: " + err.Error()})
		return
	}

	tenantID := payload.TenantID
	if tenantID == 0 {
		// Try to get from header
		tenantID, _ = h.getTenantID(c)
	}

	// Process different event types
	switch payload.EventType {
	case "call.started":
		err := h.processCallStarted(tenantID, payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case "call.ended":
		err := h.processCallEnded(tenantID, payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case "call.transcript.updated":
		err := h.processTranscriptUpdate(tenantID, payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case "call.intent.detected":
		err := h.processIntentDetection(tenantID, payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	default:
		// Unknown event type, but acknowledge receipt
		c.JSON(http.StatusOK, gin.H{"status": "received", "message": "Event type not processed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

func (h *Handlers) processCallStarted(tenantID int64, payload ConvinWebhookPayload) error {
	// Extract call data
	data := payload.Data
	externalID := payload.CallID
	
	// Get channel ID for "call"
	var channelID int
	err := h.ingestionSvc.GetDB().Get(&channelID, `SELECT id FROM channels WHERE name = $1`, "call")
	if err != nil {
		return err
	}

	// Extract vendor ID if available
	var vendorID *int
	var vendorCode *string
	if vc, ok := data["vendor_code"].(string); ok && vc != "" {
		vendorCode = &vc
		err = h.ingestionSvc.GetDB().Get(&vendorID, 
			`SELECT id FROM vendors WHERE tenant_id = $1 AND code = $2`, 
			tenantID, vc)
		if err != nil {
			vendorID = nil
		}
	}

	// Extract customer identifiers
	var customerIdentifiers []models.CustomerIdentifier
	if phone, ok := data["phone_number"].(string); ok {
		customerIdentifiers = append(customerIdentifiers, models.CustomerIdentifier{
			Type:  "phone",
			Value: phone,
		})
	}
	if email, ok := data["email"].(string); ok {
		customerIdentifiers = append(customerIdentifiers, models.CustomerIdentifier{
			Type:  "email",
			Value: email,
		})
	}

	// Parse started_at
	startedAt := payload.Timestamp
	if startedAtStr, ok := data["started_at"].(string); ok {
		if parsed, err := time.Parse(time.RFC3339, startedAtStr); err == nil {
			startedAt = parsed
		}
	}

	// Create interaction request
	req := services.IngestInteractionRequest{
		ExternalInteractionID: externalID,
		Channel:               "call",
		VendorCode:            vendorCode,
		StartedAt:             startedAt,
		Direction:            getString(data, "direction", "inbound"),
		Language:              getString(data, "language", "en"),
		CustomerIdentifiers:   customerIdentifiers,
		RawMetadata:           data,
	}

	// Ingest interaction
	_, err = h.ingestionSvc.IngestInteraction(tenantID, req)
	return err
}

func (h *Handlers) processCallEnded(tenantID int64, payload ConvinWebhookPayload) error {
	data := payload.Data
	externalID := payload.CallID

	// Parse ended_at
	endedAt := payload.Timestamp
	if endedAtStr, ok := data["ended_at"].(string); ok {
		if parsed, err := time.Parse(time.RFC3339, endedAtStr); err == nil {
			endedAt = parsed
		}
	}

	// Update interaction with end time and duration
	query := `
		UPDATE interactions 
		SET ended_at = $1, 
		    duration_seconds = EXTRACT(EPOCH FROM ($1 - started_at))::int,
		    updated_at = NOW()
		WHERE external_interaction_id = $2 AND tenant_id = $3
		RETURNING id`

	var interactionID int64
	err := h.ingestionSvc.GetDB().Get(&interactionID, query, endedAt, externalID, tenantID)
	if err != nil {
		return err
	}

	// Update transcript location if available
	if transcriptURL, ok := data["transcript_url"].(string); ok {
		_, err = h.ingestionSvc.GetDB().Exec(
			`UPDATE interactions SET transcript_location = $1 WHERE id = $2`,
			transcriptURL, interactionID,
		)
	}

	// Update outcome if available
	if outcome, ok := data["outcome"].(string); ok {
		_, err = h.ingestionSvc.GetDB().Exec(
			`UPDATE interactions SET outcome_prediction = $1 WHERE id = $2`,
			outcome, interactionID,
		)
	}

	return nil
}

func (h *Handlers) processTranscriptUpdate(tenantID int64, payload ConvinWebhookPayload) error {
	data := payload.Data
	externalID := payload.CallID

	// Update transcript location
	if transcriptURL, ok := data["transcript_url"].(string); ok {
		_, err := h.ingestionSvc.GetDB().Exec(
			`UPDATE interactions 
			 SET transcript_location = $1, updated_at = NOW()
			 WHERE external_interaction_id = $2 AND tenant_id = $3`,
			transcriptURL, externalID, tenantID,
		)
		return err
	}

	return nil
}

func (h *Handlers) processIntentDetection(tenantID int64, payload ConvinWebhookPayload) error {
	data := payload.Data
	externalID := payload.CallID

	// Update primary intent
	if primaryIntent, ok := data["primary_intent"].(string); ok {
		_, err := h.ingestionSvc.GetDB().Exec(
			`UPDATE interactions 
			 SET primary_intent = $1, updated_at = NOW()
			 WHERE external_interaction_id = $2 AND tenant_id = $3`,
			primaryIntent, externalID, tenantID,
		)
		if err != nil {
			return err
		}
	}

	// Update secondary intents if available
	if secondaryIntents, ok := data["secondary_intents"].([]interface{}); ok && len(secondaryIntents) > 0 {
		// Convert to JSONB format
		// This would need proper JSONB handling - for now, store in raw_metadata
		// In production, implement proper JSONB update
		_ = secondaryIntents // TODO: Implement JSONB update
	}

	// Update purchase probability if available
	if prob, ok := data["purchase_probability"].(float64); ok {
		_, err := h.ingestionSvc.GetDB().Exec(
			`UPDATE interactions 
			 SET purchase_probability = $1, updated_at = NOW()
			 WHERE external_interaction_id = $2 AND tenant_id = $3`,
			prob, externalID, tenantID,
		)
		return err
	}

	return nil
}

// HandleGenericTelephonyWebhook handles webhooks from generic telephony providers
func (h *Handlers) HandleGenericTelephonyWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	tenantID, _ := h.getTenantID(c)

	// Extract common fields
	callID, _ := payload["call_id"].(string)
	eventType, _ := payload["event"].(string)

	// Convert to standard format and process
	convinPayload := ConvinWebhookPayload{
		EventType: eventType,
		CallID:    callID,
		TenantID:  tenantID,
		Data:      payload,
		Timestamp: time.Now(),
	}

	// Process based on event type
	switch eventType {
	case "call.start", "call.started":
		convinPayload.EventType = "call.started"
		h.processCallStarted(tenantID, convinPayload)
	case "call.end", "call.ended":
		convinPayload.EventType = "call.ended"
		h.processCallEnded(tenantID, convinPayload)
	default:
		c.JSON(http.StatusOK, gin.H{"status": "received"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "processed"})
}

// Helper function to safely get string from map
func getString(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return defaultValue
}

