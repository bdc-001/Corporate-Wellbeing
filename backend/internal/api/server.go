package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/convin/crae/internal/api/handlers"
	"github.com/convin/crae/internal/config"
	"github.com/convin/crae/internal/database"
	"github.com/convin/crae/internal/logger"
	"github.com/convin/crae/internal/middleware"
	"github.com/convin/crae/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Server struct {
	Router     *gin.Engine
	db         *sqlx.DB
	logger     *zap.Logger
	cfg        *config.Config
	httpServer *http.Server
}

func NewServer(db *sqlx.DB, cfg *config.Config) (*Server, error) {
	// Initialize logger
	appLogger, err := logger.NewLogger(cfg.LogLevel, cfg.LogFormat, cfg.LogFile)
	if err != nil {
		return nil, err
	}

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	router := gin.New()

	// Add recovery middleware with structured logging
	router.Use(middleware.RecoveryMiddleware(appLogger))

	// Add request logging middleware
	router.Use(middleware.LoggerMiddleware(appLogger))

	// Configure CORS
	corsConfig := cors.Config{
		AllowOrigins:     cfg.CORSAllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-Tenant-ID", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: cfg.CORSAllowCredentials,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	// Add rate limiting if enabled
	if cfg.RateLimitEnabled {
		limiter := middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)
		router.Use(middleware.RateLimitMiddleware(limiter))
	}

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
		// Webhooks for Live Call Flow Integration
		// ====================================================================
		if cfg.EnableWebhooks {
			webhooks := v1.Group("/webhooks")
			{
				// Convin webhook with signature verification
				convinWebhooks := webhooks.Group("/convin")
				if cfg.ConvinWebhookSecret != "" {
					convinWebhooks.Use(middleware.WebhookSignatureMiddleware(cfg.ConvinWebhookSecret))
				}
				convinWebhooks.POST("", h.HandleConvinWebhook)

				// Generic telephony webhook
				telephonyWebhooks := webhooks.Group("/telephony")
				if cfg.TelephonyWebhookSecret != "" {
					telephonyWebhooks.Use(middleware.WebhookSignatureMiddleware(cfg.TelephonyWebhookSecret))
				}
				telephonyWebhooks.POST("", h.HandleGenericTelephonyWebhook)
			}
		}

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
			realtime.POST("/alerts", h.CreateAlert)                      // Create alert (for testing)
			realtime.GET("/alerts", h.GetAlerts)                         // Get alerts
			realtime.POST("/alerts/:id/acknowledge", h.AcknowledgeAlert) // Acknowledge alert
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
		// Check database connectivity
		dbHealthy := true
		if err := database.HealthCheck(db); err != nil {
			dbHealthy = false
			appLogger.Error("Database health check failed", zap.Error(err))
		}

		status := http.StatusOK
		if !dbHealthy {
			status = http.StatusServiceUnavailable
		}

		c.JSON(status, gin.H{
			"status":    "ok",
			"service":   "Convin Revenue Attribution Engine (CRAE)",
			"version":   "2.0.0-production",
			"timestamp": time.Now().UTC(),
			"database":  map[string]interface{}{"healthy": dbHealthy},
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
				"Live Call Flow Integration",
			},
		})
	})

	// Readiness probe
	router.GET("/ready", func(c *gin.Context) {
		if err := database.HealthCheck(db); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	// Liveness probe
	router.GET("/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})

	return &Server{
		Router: router,
		db:     db,
		logger: appLogger,
		cfg:    cfg,
	}, nil
}

// Start starts the HTTP server with graceful shutdown
func (s *Server) Start() error {
	s.httpServer = &http.Server{
		Addr:         ":" + s.cfg.Port,
		Handler:      s.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		s.logger.Info("Server starting", zap.String("port", s.cfg.Port), zap.String("environment", s.cfg.Environment))
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("Server forced to shutdown", zap.Error(err))
		return err
	}

	s.logger.Info("Server exited gracefully")
	return nil
}

// Close closes database connections
func (s *Server) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
