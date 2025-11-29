package api

import (
	"time"

	"github.com/convin/crae/internal/api/handlers"
	"github.com/convin/crae/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	Router *gin.Engine
	db     *sqlx.DB
}

func NewServer(db *sqlx.DB, cfg interface{}) *Server {
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-Tenant-ID", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize all services
	identitySvc := services.NewIdentityService(db)
	ingestionSvc := services.NewIngestionService(db, identitySvc)
	attributionSvc := services.NewAttributionService(db)
	analyticsSvc := services.NewAnalyticsService(db)
	advancedAnalyticsSvc := services.NewAdvancedAnalyticsService(db)
	abmSvc := services.NewABMService(db)
	leadScoringSvc := services.NewLeadScoringService(db)
	cohortSvc := services.NewCohortService(db)
	realtimeSvc := services.NewRealtimeService(db)
	fraudSvc := services.NewFraudService(db)
	behaviorSvc := services.NewBehaviorService(db)
	integrationSvc := services.NewIntegrationService(db)
	reportSvc := services.NewReportService(db)
	experimentSvc := services.NewExperimentService(db)
	mmmSvc := services.NewMMMService(db)
	userMgmtSvc := services.NewUserManagementService(db)
	roleMgmtSvc := services.NewRoleManagementService(db)
	teamMgmtSvc := services.NewTeamManagementService(db)

	// Initialize handlers with all services
	h := handlers.NewHandlers(
		identitySvc,
		ingestionSvc,
		attributionSvc,
		analyticsSvc,
		advancedAnalyticsSvc,
		abmSvc,
		leadScoringSvc,
		cohortSvc,
		realtimeSvc,
		fraudSvc,
		behaviorSvc,
		integrationSvc,
		reportSvc,
		experimentSvc,
		mmmSvc,
		userMgmtSvc,
		roleMgmtSvc,
		teamMgmtSvc,
	)

	// ========================================================================
	// API Routes - Comprehensive Feature Set
	// ========================================================================

	v1 := router.Group("/v1")
	{
		// ====================================================================
		// Data Ingestion APIs
		// ====================================================================
		v1.POST("/interactions", h.IngestInteraction)
		v1.POST("/conversions", h.IngestConversion)
		v1.POST("/events", h.IngestEvent)
		v1.POST("/page-views", h.TrackPageView)

		// ====================================================================
		// Customer Identity & Journey
		// ====================================================================
		v1.GET("/customers/:customer_id/journey", h.GetCustomerJourney)

		// ====================================================================
		// Attribution Engine
		// ====================================================================
		v1.POST("/attribution/runs", h.CreateAttributionRun)
		v1.GET("/attribution/runs/:run_id", h.GetAttributionRun)
		v1.POST("/attribution/runs/:run_id/execute", h.ExecuteAttributionRun)

		// ====================================================================
		// Core Analytics
		// ====================================================================
		analytics := v1.Group("/analytics")
		{
			// Traditional analytics
			analytics.GET("/agents/revenue", h.GetAgentRevenueSummary)
			analytics.GET("/vendors/comparison", h.GetVendorComparison)
			analytics.GET("/intents/revenue", h.GetIntentProfitability)

			// Advanced analytics (Factors.ai-style)
			analytics.GET("/funnel/stages", h.GetFunnelStageMetrics)
			analytics.GET("/content/engagement", h.GetContentEngagementMetrics)
			analytics.GET("/channels/roi", h.GetMultiChannelROI)
			analytics.GET("/journey/velocity", h.GetJourneyVelocity)
			analytics.POST("/reports/custom", h.GetCustomReport)

			// Real-time metrics
			analytics.GET("/realtime/metrics", h.GetRealtimeMetrics)
		}

		// ====================================================================
		// Vendors
		// ====================================================================
		v1.GET("/vendors", h.ListVendors)

		// ====================================================================
		// Account-Based Marketing (ABM)
		// ====================================================================
		abm := v1.Group("/abm")
		{
			abm.POST("/accounts", h.CreateAccount)
			abm.GET("/accounts", h.ListAccounts)
			abm.GET("/accounts/:id", h.GetAccount)
			abm.GET("/accounts/:id/summary", h.GetAccountSummary)
			abm.POST("/accounts/engagements", h.TrackAccountEngagement)
			abm.GET("/insights/target-accounts", h.GetTargetAccountInsights)
		}

		// ====================================================================
		// Lead Scoring & Predictive Analytics
		// ====================================================================
		leads := v1.Group("/leads")
		{
			leads.POST("/customers/:customer_id/score", h.CalculateLeadScore)
			leads.GET("/high-value", h.GetHighValueLeads)
			leads.POST("/predictions", h.CreatePrediction)
		}

		// ====================================================================
		// Cohort Analysis & Segmentation
		// ====================================================================
		cohorts := v1.Group("/cohorts")
		{
			cohorts.POST("/compute", h.ComputeCohortMetrics)
			cohorts.GET("/segments/:segment_id/retention", h.GetRetentionCurve)
		}

		// ====================================================================
		// Real-Time Data & Alerts
		// ====================================================================
		realtime := v1.Group("/realtime")
		{
			realtime.GET("/alerts", h.GetAlerts)
			realtime.POST("/alerts/:id/acknowledge", h.AcknowledgeAlert)
		}

		// ====================================================================
		// Fraud Detection & Data Quality
		// ====================================================================
		fraud := v1.Group("/fraud")
		{
			fraud.GET("/detect", h.DetectFraud)
			fraud.GET("/incidents", h.GetFraudIncidents)
		}

		quality := v1.Group("/quality")
		{
			quality.GET("/scores", h.CalculateDataQuality)
		}

		// ====================================================================
		// User Behavior Analytics
		// ====================================================================
		behavior := v1.Group("/behavior")
		{
			behavior.GET("/sessions/:session_id", h.GetSessionDetails)
			behavior.GET("/pages/top", h.GetTopPages)
		}

		// ====================================================================
		// Integrations
		// ====================================================================
		integrations := v1.Group("/integrations")
		{
			integrations.POST("", h.CreateIntegration)
			integrations.GET("", h.ListIntegrations)
			integrations.POST("/:id/sync", h.SyncIntegration)
		}

		// ====================================================================
		// Custom Reports & Saved Queries
		// ====================================================================
		reports := v1.Group("/reports")
		{
			reports.POST("", h.CreateReport)
			reports.GET("", h.ListReports)
			reports.POST("/:id/execute", h.ExecuteReport)
		}

		// ====================================================================
		// A/B Testing & Experiments
		// ====================================================================
		experiments := v1.Group("/experiments")
		{
			experiments.POST("", h.CreateExperiment)
			experiments.GET("/:id/results", h.GetExperimentResults)
		}

		// ====================================================================
		// Feature Flags
		// ====================================================================
		features := v1.Group("/features")
		{
			features.POST("/flags", h.CreateFeatureFlag)
			features.GET("/flags", h.ListFeatureFlags)
		}

		// ====================================================================
		// Marketing Mix Modeling (MMM)
		// ====================================================================
		mmm := v1.Group("/mmm")
		{
			mmm.POST("/run", h.RunMMMAnalysis)
			mmm.GET("/models", h.GetMMMModels)
			mmm.GET("/models/:model_id/results", h.GetMMMResults)
		}

		// ====================================================================
		// User Management
		// ====================================================================
		users := v1.Group("/users")
		{
			users.POST("", h.CreateUser)
			users.POST("/bulk", h.BulkCreateUsers)
			users.POST("/login", h.Login)
			users.GET("", h.ListUsers)
			users.GET("/:id", h.GetUser)
			users.PUT("/:id", h.UpdateUser)
			users.DELETE("/:id", h.DeleteUser)
		}

		// ====================================================================
		// Role Management
		// ====================================================================
		roles := v1.Group("/roles")
		{
			roles.POST("", h.CreateRole)
			roles.GET("", h.ListRoles)
			roles.GET("/:id", h.GetRole)
			roles.PUT("/:id", h.UpdateRole)
			roles.DELETE("/:id", h.DeleteRole)
		}

		// ====================================================================
		// Permissions
		// ====================================================================
		v1.GET("/permissions", h.ListPermissions)

		// ====================================================================
		// Team Management
		// ====================================================================
		teams := v1.Group("/teams")
		{
			teams.POST("", h.CreateTeam)
			teams.GET("", h.ListTeams)
			teams.GET("/:id", h.GetTeam)
			teams.PUT("/:id", h.UpdateTeam)
			teams.DELETE("/:id", h.DeleteTeam)
			teams.POST("/:id/members", h.AddTeamMembers)
		}

		// ====================================================================
		// Use Cases
		// ====================================================================
		useCases := v1.Group("/use-cases")
		{
			useCases.POST("", h.CreateUseCase)
			useCases.GET("", h.ListUseCases)
			useCases.DELETE("/:id", h.DeleteUseCase)
		}
	}

	// ========================================================================
	// Health Check & System Info
	// ========================================================================
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "Convin Revenue Attribution Engine (CRAE)",
			"version": "2.0.0-comprehensive",
			"features": []string{
				"Multi-Touch Attribution",
				"Account-Based Marketing (ABM)",
				"Lead Scoring & Predictive Analytics",
				"Cohort Analysis",
				"Real-Time Data Streaming",
				"Fraud Detection",
				"User Behavior Analytics",
				"CRM & Ad Platform Integrations",
				"Custom Reports & Dashboards",
				"A/B Testing & Experiments",
				"Feature Flags",
				"Marketing Mix Modeling (MMM)",
			},
		})
	})

	return &Server{
		Router: router,
		db:     db,
	}
}
