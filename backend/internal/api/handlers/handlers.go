package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/convin/crae/internal/services"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	identitySvc          *services.IdentityService
	ingestionSvc         *services.IngestionService
	attributionSvc       *services.AttributionService
	analyticsSvc         *services.AnalyticsService
	advancedAnalyticsSvc *services.AdvancedAnalyticsService
	abmSvc               *services.ABMService
	leadScoringSvc       *services.LeadScoringService
	cohortSvc            *services.CohortService
	realtimeSvc          *services.RealtimeService
	fraudSvc             *services.FraudService
	behaviorSvc          *services.BehaviorService
	integrationSvc       *services.IntegrationService
	reportSvc            *services.ReportService
	experimentSvc        *services.ExperimentService
	mmmSvc               *services.MMMService
	userMgmtSvc          *services.UserManagementService
	roleMgmtSvc          *services.RoleManagementService
	teamMgmtSvc          *services.TeamManagementService
}

func NewHandlers(
	identitySvc *services.IdentityService,
	ingestionSvc *services.IngestionService,
	attributionSvc *services.AttributionService,
	analyticsSvc *services.AnalyticsService,
	advancedAnalyticsSvc *services.AdvancedAnalyticsService,
	abmSvc *services.ABMService,
	leadScoringSvc *services.LeadScoringService,
	cohortSvc *services.CohortService,
	realtimeSvc *services.RealtimeService,
	fraudSvc *services.FraudService,
	behaviorSvc *services.BehaviorService,
	integrationSvc *services.IntegrationService,
	reportSvc *services.ReportService,
	experimentSvc *services.ExperimentService,
	mmmSvc *services.MMMService,
	userMgmtSvc *services.UserManagementService,
	roleMgmtSvc *services.RoleManagementService,
	teamMgmtSvc *services.TeamManagementService,
) *Handlers {
	return &Handlers{
		identitySvc:          identitySvc,
		ingestionSvc:         ingestionSvc,
		attributionSvc:       attributionSvc,
		analyticsSvc:         analyticsSvc,
		advancedAnalyticsSvc: advancedAnalyticsSvc,
		abmSvc:               abmSvc,
		leadScoringSvc:       leadScoringSvc,
		cohortSvc:            cohortSvc,
		realtimeSvc:          realtimeSvc,
		fraudSvc:             fraudSvc,
		behaviorSvc:          behaviorSvc,
		integrationSvc:       integrationSvc,
		reportSvc:            reportSvc,
		experimentSvc:        experimentSvc,
		mmmSvc:               mmmSvc,
		userMgmtSvc:          userMgmtSvc,
		roleMgmtSvc:          roleMgmtSvc,
		teamMgmtSvc:          teamMgmtSvc,
	}
}

// getTenantID extracts tenant ID from request (from token or header)
func (h *Handlers) getTenantID(c *gin.Context) (int64, error) {
	// In production, extract from JWT token
	// For now, use header or default to 1
	tenantIDStr := c.GetHeader("X-Tenant-ID")
	if tenantIDStr == "" {
		tenantIDStr = "1" // Default for development
	}

	tenantID, err := strconv.ParseInt(tenantIDStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return tenantID, nil
}

// IngestInteraction handles interaction ingestion
func (h *Handlers) IngestInteraction(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var req services.IngestInteractionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.ingestionSvc.IngestInteraction(tenantID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// IngestConversion handles conversion event ingestion
func (h *Handlers) IngestConversion(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var req services.IngestConversionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.ingestionSvc.IngestConversion(tenantID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetCustomerJourney returns customer journey
func (h *Handlers) GetCustomerJourney(c *gin.Context) {
	customerIDStr := c.Param("customer_id")
	customerID, err := strconv.ParseInt(customerIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	from := c.Query("from")
	to := c.Query("to")

	var fromPtr, toPtr *string
	if from != "" {
		fromPtr = &from
	}
	if to != "" {
		toPtr = &to
	}

	journey, err := h.identitySvc.GetCustomerJourney(customerID, fromPtr, toPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, journey)
}

// CreateAttributionRun creates a new attribution run
func (h *Handlers) CreateAttributionRun(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var req struct {
		ModelCode string                     `json:"model_code"`
		Name      string                     `json:"name"`
		Config    services.AttributionConfig `json:"config"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	run, err := h.attributionSvc.CreateAttributionRun(tenantID, req.ModelCode, req.Name, req.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, run)
}

// GetAttributionRun returns attribution run details
func (h *Handlers) GetAttributionRun(c *gin.Context) {
	runIDStr := c.Param("run_id")
	_, err := strconv.ParseInt(runIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid run ID"})
		return
	}

	// This would need to be implemented in the service
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// ExecuteAttributionRun executes an attribution run
func (h *Handlers) ExecuteAttributionRun(c *gin.Context) {
	runIDStr := c.Param("run_id")
	runID, err := strconv.ParseInt(runIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid run ID"})
		return
	}

	err = h.attributionSvc.ExecuteAttributionRun(runID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "completed"})
}

// GetAgentRevenueSummary returns agent revenue summary
func (h *Handlers) GetAgentRevenueSummary(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var from, to *time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err == nil {
			from = &t
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err == nil {
			to = &t
		}
	}

	var vendorID *int
	if vendorIDStr := c.Query("vendor_id"); vendorIDStr != "" {
		id, err := strconv.Atoi(vendorIDStr)
		if err == nil {
			vendorID = &id
		}
	}

	modelCode := c.DefaultQuery("model_code", "AI_WEIGHTED")

	results, err := h.analyticsSvc.GetAgentRevenueSummary(tenantID, from, to, vendorID, modelCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"from":       from,
		"to":         to,
		"model_code": modelCode,
		"agents":     results,
	})
}

// GetVendorComparison returns vendor comparison
func (h *Handlers) GetVendorComparison(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var from, to *time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err == nil {
			from = &t
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err == nil {
			to = &t
		}
	}

	modelCode := c.DefaultQuery("model_code", "AI_WEIGHTED")

	results, err := h.analyticsSvc.GetVendorComparison(tenantID, from, to, modelCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"vendors": results,
	})
}

// GetIntentProfitability returns intent-level profitability
func (h *Handlers) GetIntentProfitability(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var from, to *time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err == nil {
			from = &t
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err == nil {
			to = &t
		}
	}

	modelCode := c.DefaultQuery("model_code", "AI_WEIGHTED")

	results, err := h.analyticsSvc.GetIntentProfitability(tenantID, from, to, modelCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"intents": results,
	})
}

// GetFunnelStageMetrics returns funnel stage metrics
func (h *Handlers) GetFunnelStageMetrics(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var from, to *time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err == nil {
			from = &t
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err == nil {
			to = &t
		}
	}

	var segmentID *int
	if segmentIDStr := c.Query("segment_id"); segmentIDStr != "" {
		id, err := strconv.Atoi(segmentIDStr)
		if err == nil {
			segmentID = &id
		}
	}

	results, err := h.advancedAnalyticsSvc.GetFunnelStageMetrics(tenantID, from, to, segmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stages": results})
}

// GetContentEngagementMetrics returns content engagement metrics
func (h *Handlers) GetContentEngagementMetrics(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var from, to *time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err == nil {
			from = &t
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err == nil {
			to = &t
		}
	}

	results, err := h.advancedAnalyticsSvc.GetContentEngagementMetrics(tenantID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": results})
}

// GetMultiChannelROI returns multi-channel ROI metrics
func (h *Handlers) GetMultiChannelROI(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var from, to *time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err == nil {
			from = &t
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err == nil {
			to = &t
		}
	}

	results, err := h.advancedAnalyticsSvc.GetMultiChannelROI(tenantID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"channels": results})
}

// GetJourneyVelocity returns journey velocity metrics
func (h *Handlers) GetJourneyVelocity(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var from, to *time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err == nil {
			from = &t
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err == nil {
			to = &t
		}
	}

	results, err := h.advancedAnalyticsSvc.GetJourneyVelocity(tenantID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"velocity": results})
}

// GetCustomReport generates a custom report
func (h *Handlers) GetCustomReport(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var req services.CustomReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var from, to *time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err == nil {
			from = &t
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err == nil {
			to = &t
		}
	}

	results, err := h.advancedAnalyticsSvc.GetCustomReport(tenantID, req, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}

// ============================================================================
// ABM (Account-Based Marketing) Handlers
// ============================================================================

// CreateAccount creates a new ABM account
func (h *Handlers) CreateAccount(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var account services.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.abmSvc.CreateAccount(tenantID, &account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// GetAccount retrieves an account
func (h *Handlers) GetAccount(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	accountID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	account, err := h.abmSvc.GetAccount(tenantID, accountID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, account)
}

// ListAccounts lists accounts
func (h *Handlers) ListAccounts(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	targetOnly := c.Query("target_only") == "true"
	tier := c.Query("tier")
	lifecycle := c.Query("lifecycle")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	accounts, err := h.abmSvc.ListAccounts(tenantID, targetOnly, tier, lifecycle, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

// GetAccountSummary retrieves account summary with metrics
func (h *Handlers) GetAccountSummary(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	accountID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	summary, err := h.abmSvc.GetAccountSummary(tenantID, accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// TrackAccountEngagement tracks an account engagement event
func (h *Handlers) TrackAccountEngagement(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var engagement services.AccountEngagement
	if err := c.ShouldBindJSON(&engagement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.abmSvc.TrackAccountEngagement(tenantID, &engagement)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, engagement)
}

// GetTargetAccountInsights retrieves insights for target accounts
func (h *Handlers) GetTargetAccountInsights(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	tier := c.Query("tier")
	insights, err := h.abmSvc.GetTargetAccountInsights(tenantID, tier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, insights)
}

// ============================================================================
// Lead Scoring & Predictive Analytics Handlers
// ============================================================================

// CalculateLeadScore calculates lead score for a customer
func (h *Handlers) CalculateLeadScore(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	customerID, err := strconv.ParseInt(c.Param("customer_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	score, err := h.leadScoringSvc.CalculateLeadScore(tenantID, customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, score)
}

// CreatePrediction creates a new prediction
func (h *Handlers) CreatePrediction(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var prediction services.Prediction
	if err := c.ShouldBindJSON(&prediction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.leadScoringSvc.CreatePrediction(tenantID, &prediction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, prediction)
}

// GetHighValueLeads retrieves high-value leads
func (h *Handlers) GetHighValueLeads(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	minScore, _ := strconv.ParseFloat(c.DefaultQuery("min_score", "70"), 64)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	leads, err := h.leadScoringSvc.GetHighValueLeads(tenantID, minScore, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"leads": leads})
}

// ============================================================================
// Cohort Analysis Handlers
// ============================================================================

// ComputeCohortMetrics computes cohort metrics
func (h *Handlers) ComputeCohortMetrics(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var req services.CohortAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.TenantID = tenantID

	results, err := h.cohortSvc.ComputeCohortMetrics(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}

// GetRetentionCurve retrieves retention curve data
func (h *Handlers) GetRetentionCurve(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	segmentID, err := strconv.ParseInt(c.Param("segment_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid segment ID"})
		return
	}

	curve, err := h.cohortSvc.GetRetentionCurve(tenantID, segmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"retention_curve": curve})
}

// ============================================================================
// Real-Time & Alerts Handlers
// ============================================================================

// IngestEvent ingests a real-time event
func (h *Handlers) IngestEvent(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var event services.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.realtimeSvc.IngestEvent(tenantID, &event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, event)
}

// CreateAlert creates a new alert (for testing/admin purposes)
func (h *Handlers) CreateAlert(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var alert services.Alert
	if err := c.ShouldBindJSON(&alert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if alert.AlertType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "alert_type is required"})
		return
	}
	if alert.Severity == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "severity is required"})
		return
	}
	if alert.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	err = h.realtimeSvc.CreateAlert(tenantID, &alert)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, alert)
}

// GetAlerts retrieves alerts
func (h *Handlers) GetAlerts(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	severity := c.Query("severity")
	alertType := c.Query("type")
	acknowledgedOnly := c.Query("acknowledged") == "true"
	unresolvedOnly := c.Query("unresolved") == "true"
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	alerts, err := h.realtimeSvc.GetAlerts(tenantID, severity, alertType, acknowledgedOnly, unresolvedOnly, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"alerts": alerts})
}

// AcknowledgeAlert acknowledges an alert
func (h *Handlers) AcknowledgeAlert(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	alertID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
		return
	}

	// Try to get from body first, then query param
	var req struct {
		By string `json:"by"`
	}
	if err := c.ShouldBindJSON(&req); err == nil && req.By != "" {
		err = h.realtimeSvc.AcknowledgeAlert(tenantID, alertID, req.By)
	} else {
		acknowledgedBy := c.Query("by")
		if acknowledgedBy == "" {
			acknowledgedBy = "system" // Default if not provided
		}
		err = h.realtimeSvc.AcknowledgeAlert(tenantID, alertID, acknowledgedBy)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Alert acknowledged"})
}

// GetRealtimeMetrics retrieves real-time metrics
func (h *Handlers) GetRealtimeMetrics(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	timeWindow, _ := strconv.Atoi(c.DefaultQuery("window", "15"))
	metrics, err := h.realtimeSvc.GetRealtimeMetrics(tenantID, timeWindow)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// ============================================================================
// Fraud Detection Handlers
// ============================================================================

// DetectFraud runs fraud detection on an entity
func (h *Handlers) DetectFraud(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	entityType := c.Query("entity_type")
	entityID, err := strconv.ParseInt(c.Query("entity_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	incidents, err := h.fraudSvc.DetectFraud(tenantID, entityType, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"incidents": incidents})
}

// GetFraudIncidents retrieves fraud incidents
func (h *Handlers) GetFraudIncidents(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	status := c.Query("status")
	severity := c.Query("severity")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	incidents, err := h.fraudSvc.GetFraudIncidents(tenantID, status, severity, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"incidents": incidents})
}

// CalculateDataQuality calculates data quality scores
func (h *Handlers) CalculateDataQuality(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	entityType := c.Query("entity_type")
	entityID, err := strconv.ParseInt(c.Query("entity_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	score, err := h.fraudSvc.CalculateDataQuality(tenantID, entityType, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, score)
}

// ============================================================================
// User Behavior Analytics Handlers
// ============================================================================

// TrackPageView tracks a page view
func (h *Handlers) TrackPageView(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var pageView services.PageView
	if err := c.ShouldBindJSON(&pageView); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.behaviorSvc.TrackPageView(tenantID, &pageView)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, pageView)
}

// GetSessionDetails retrieves session details
func (h *Handlers) GetSessionDetails(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	sessionID := c.Param("session_id")
	details, err := h.behaviorSvc.GetSessionDetails(tenantID, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, details)
}

// GetTopPages retrieves top pages
func (h *Handlers) GetTopPages(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	startDate, _ := time.Parse(time.RFC3339, c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -7).Format(time.RFC3339)))
	endDate, _ := time.Parse(time.RFC3339, c.DefaultQuery("end_date", time.Now().Format(time.RFC3339)))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	pages, err := h.behaviorSvc.GetTopPages(tenantID, startDate, endDate, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pages": pages})
}

// ============================================================================
// Integration Handlers
// ============================================================================

// CreateIntegration creates a new integration
func (h *Handlers) CreateIntegration(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var integration services.Integration
	if err := c.ShouldBindJSON(&integration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.integrationSvc.CreateIntegration(tenantID, &integration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, integration)
}

// ListIntegrations lists integrations
func (h *Handlers) ListIntegrations(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	platform := c.Query("platform")
	integrationType := c.Query("type")
	activeOnly := c.Query("active_only") == "true"

	integrations, err := h.integrationSvc.ListIntegrations(tenantID, platform, integrationType, activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"integrations": integrations})
}

// SyncIntegration syncs an integration
func (h *Handlers) SyncIntegration(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	integrationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid integration ID"})
		return
	}

	err = h.integrationSvc.SyncFromCRM(tenantID, integrationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sync started"})
}

// ============================================================================
// Custom Reports Handlers
// ============================================================================

// CreateReport creates a saved report
func (h *Handlers) CreateReport(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var report services.SavedReport
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.reportSvc.CreateReport(tenantID, &report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, report)
}

// ListReports lists saved reports
func (h *Handlers) ListReports(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	reportType := c.Query("type")
	publicOnly := c.Query("public_only") == "true"

	reports, err := h.reportSvc.ListReports(tenantID, reportType, publicOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reports": reports})
}

// ExecuteReport executes a report
func (h *Handlers) ExecuteReport(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	reportID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	params := make(map[string]interface{})
	c.ShouldBindJSON(&params)

	result, err := h.reportSvc.ExecuteReport(tenantID, reportID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ============================================================================
// Experiment Handlers
// ============================================================================

// CreateExperiment creates a new A/B test
func (h *Handlers) CreateExperiment(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var experiment services.Experiment
	if err := c.ShouldBindJSON(&experiment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.experimentSvc.CreateExperiment(tenantID, &experiment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, experiment)
}

// GetExperimentResults retrieves experiment results
func (h *Handlers) GetExperimentResults(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	experimentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid experiment ID"})
		return
	}

	results, err := h.experimentSvc.GetExperimentResults(tenantID, experimentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// CreateFeatureFlag creates a feature flag
func (h *Handlers) CreateFeatureFlag(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var flag services.FeatureFlag
	if err := c.ShouldBindJSON(&flag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.experimentSvc.CreateFeatureFlag(tenantID, &flag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, flag)
}

// ListFeatureFlags lists feature flags
func (h *Handlers) ListFeatureFlags(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	flags, err := h.experimentSvc.ListFeatureFlags(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"flags": flags})
}

// ============================================================================
// Marketing Mix Modeling (MMM) Handlers
// ============================================================================

// RunMMMAnalysis runs marketing mix modeling analysis
func (h *Handlers) RunMMMAnalysis(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var req services.MMMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.TenantID = tenantID

	results, err := h.mmmSvc.RunMMMAnalysis(tenantID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetMMMModels lists MMM models
func (h *Handlers) GetMMMModels(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	status := c.Query("status")
	models, err := h.mmmSvc.ListMMMModels(tenantID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"models": models})
}

// GetMMMResults retrieves MMM results
func (h *Handlers) GetMMMResults(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	modelID, err := strconv.ParseInt(c.Param("model_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	results, err := h.mmmSvc.GetMMMResults(tenantID, modelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// ========================================================================
// USER MANAGEMENT HANDLERS
// ========================================================================

// CreateUser creates a new user
func (h *Handlers) CreateUser(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var req services.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	user, password, err := h.userMgmtSvc.CreateUser(tenantID, req)
	if err != nil {
		// Check if it's a duplicate email error
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":     user,
		"password": password, // Return password for email notification
	})
}

// BulkCreateUsers creates multiple users from CSV/Excel
func (h *Handlers) BulkCreateUsers(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var req struct {
		Users []services.CreateUserRequest `json:"users" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, passwords, errors := h.userMgmtSvc.BulkCreateUsers(tenantID, req.Users)

	response := gin.H{
		"users":         users,
		"passwords":     passwords,
		"success_count": len(users),
		"error_count":   len(errors),
	}
	if len(errors) > 0 {
		response["errors"] = errors
	}

	c.JSON(http.StatusOK, response)
}

// UpdateUser updates an existing user
func (h *Handlers) UpdateUser(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req services.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userMgmtSvc.UpdateUser(tenantID, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUser retrieves a user by ID
func (h *Handlers) GetUser(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.userMgmtSvc.GetUser(tenantID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ListUsers retrieves all users for a tenant
func (h *Handlers) ListUsers(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var filters services.UserFilters
	if c.Query("role_id") != "" {
		roleID, _ := strconv.Atoi(c.Query("role_id"))
		filters.RoleID = &roleID
	}
	if c.Query("team_id") != "" {
		teamID, _ := strconv.Atoi(c.Query("team_id"))
		filters.TeamID = &teamID
	}
	if c.Query("is_active") != "" {
		isActive := c.Query("is_active") == "true"
		filters.IsActive = &isActive
	}
	if c.Query("search") != "" {
		search := c.Query("search")
		filters.Search = &search
	}

	users, err := h.userMgmtSvc.ListUsers(tenantID, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// DeleteUser deletes a user
func (h *Handlers) DeleteUser(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.userMgmtSvc.DeleteUser(tenantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// ========================================================================
// ROLE MANAGEMENT HANDLERS
// ========================================================================

// CreateRole creates a new role
func (h *Handlers) CreateRole(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var req services.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.roleMgmtSvc.CreateRole(tenantID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, role)
}

// UpdateRole updates an existing role
func (h *Handlers) UpdateRole(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var req services.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.roleMgmtSvc.UpdateRole(tenantID, roleID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

// GetRole retrieves a role by ID
func (h *Handlers) GetRole(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	role, err := h.roleMgmtSvc.GetRole(tenantID, roleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

// ListRoles retrieves all roles for a tenant
func (h *Handlers) ListRoles(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	roles, err := h.roleMgmtSvc.ListRoles(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": roles})
}

// DeleteRole deletes a role
func (h *Handlers) DeleteRole(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	err = h.roleMgmtSvc.DeleteRole(tenantID, roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

// ListPermissions retrieves all available permissions
func (h *Handlers) ListPermissions(c *gin.Context) {
	groups, err := h.roleMgmtSvc.ListPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permission_groups": groups})
}

// ========================================================================
// TEAM MANAGEMENT HANDLERS
// ========================================================================

// CreateTeam creates a new team with optional subteams and members
func (h *Handlers) CreateTeam(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var req services.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := h.teamMgmtSvc.CreateTeam(tenantID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, team)
}

// GetTeam retrieves a team with subteams and members
func (h *Handlers) GetTeam(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	teamID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	team, subteams, members, err := h.teamMgmtSvc.GetTeam(tenantID, teamID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"team":     team,
		"subteams": subteams,
		"members":  members,
	})
}

// Login authenticates a user
func (h *Handlers) Login(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userMgmtSvc.AuthenticateUser(tenantID, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// ListVendors lists all vendors for a tenant
func (h *Handlers) ListVendors(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var vendors []struct {
		ID          int    `db:"id" json:"id"`
		TenantID    int64  `db:"tenant_id" json:"tenant_id"`
		Name        string `db:"name" json:"name"`
		Code        string `db:"code" json:"code"`
		Description string `db:"description" json:"description"`
		IsActive    bool   `db:"is_active" json:"is_active"`
	}

	db := h.userMgmtSvc.GetDB()
	err = db.Select(&vendors,
		`SELECT id, tenant_id, name, code, description, is_active
		FROM vendors WHERE tenant_id = $1 AND is_active = true
		ORDER BY name ASC`,
		tenantID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"vendors": vendors})
}

// ListTeams lists all teams for a tenant
func (h *Handlers) ListTeams(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	includeSubteams := c.Query("include_subteams") == "true"
	teams, err := h.teamMgmtSvc.ListTeams(tenantID, includeSubteams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"teams": teams})
}

// UpdateTeam updates a team
func (h *Handlers) UpdateTeam(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	teamID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var req services.UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := h.teamMgmtSvc.UpdateTeam(tenantID, teamID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

// DeleteTeam deletes a team with safe member migration
func (h *Handlers) DeleteTeam(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	teamID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var req struct {
		TransferToTeamID *int `json:"transfer_to_team_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// TransferToTeamID is optional, only required if team has members
	}

	err = h.teamMgmtSvc.DeleteTeam(tenantID, teamID, req.TransferToTeamID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team deleted successfully"})
}

// AddTeamMembers adds members to a team
func (h *Handlers) AddTeamMembers(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	teamID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return
	}

	var req struct {
		Members []int64 `json:"members" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.teamMgmtSvc.AddTeamMembers(tenantID, teamID, req.Members)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Members added successfully"})
}

// Use Case Management

// CreateUseCase creates a new use case
func (h *Handlers) CreateUseCase(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	useCase, err := h.teamMgmtSvc.CreateUseCase(tenantID, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, useCase)
}

// ListUseCases lists all use cases for a tenant
func (h *Handlers) ListUseCases(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	useCases, err := h.teamMgmtSvc.ListUseCases(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"use_cases": useCases})
}

// DeleteUseCase deletes a use case
func (h *Handlers) DeleteUseCase(c *gin.Context) {
	tenantID, err := h.getTenantID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
		return
	}

	useCaseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid use case ID"})
		return
	}

	err = h.teamMgmtSvc.DeleteUseCase(tenantID, useCaseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Use case deleted successfully"})
}
