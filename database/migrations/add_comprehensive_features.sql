-- Comprehensive Feature Set Migration
-- Inspired by: Factors.ai, Bizible, HockeyStack, DreamData, 6sense, Segment, Mixpanel, Amplitude

-- ============================================================================
-- 1. ACCOUNT-BASED MARKETING (ABM) FEATURES
-- ============================================================================

-- Accounts table for B2B ABM tracking
CREATE TABLE IF NOT EXISTS accounts (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    domain VARCHAR(255),
    industry VARCHAR(100),
    company_size VARCHAR(50),
    annual_revenue DECIMAL(15, 2),
    location VARCHAR(255),
    target_account BOOLEAN DEFAULT FALSE,
    account_tier VARCHAR(50), -- Enterprise, Mid-Market, SMB
    health_score DECIMAL(5, 2), -- 0-100 score
    engagement_score DECIMAL(5, 2), -- 0-100 score
    intent_score DECIMAL(5, 2), -- 0-100 score from intent signals
    lifecycle_stage VARCHAR(50), -- Target, Engaged, MQL, SQL, Opportunity, Customer
    crm_account_id VARCHAR(100),
    owner_id VARCHAR(100), -- Account owner/rep
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    metadata JSONB
);

CREATE INDEX idx_accounts_tenant_id ON accounts(tenant_id);
CREATE INDEX idx_accounts_domain ON accounts(domain);
CREATE INDEX idx_accounts_target ON accounts(target_account);
CREATE INDEX idx_accounts_lifecycle ON accounts(lifecycle_stage);

-- Link customers to accounts
ALTER TABLE customers ADD COLUMN IF NOT EXISTS account_id BIGINT REFERENCES accounts(id);
CREATE INDEX IF NOT EXISTS idx_customers_account_id ON customers(account_id);

-- Account engagement events
CREATE TABLE IF NOT EXISTS account_engagements (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    account_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    engagement_type VARCHAR(100) NOT NULL, -- website_visit, content_download, demo_request, etc.
    engagement_date TIMESTAMPTZ NOT NULL,
    touchpoint_id BIGINT, -- References interaction or other event
    score_impact DECIMAL(5, 2), -- How much this impacted engagement score
    created_at TIMESTAMPTZ DEFAULT NOW(),
    metadata JSONB
);

CREATE INDEX idx_account_engagements_tenant ON account_engagements(tenant_id);
CREATE INDEX idx_account_engagements_account ON account_engagements(account_id);
CREATE INDEX idx_account_engagements_date ON account_engagements(engagement_date);

-- ============================================================================
-- 2. LEAD SCORING & PREDICTIVE ANALYTICS
-- ============================================================================

-- Lead scoring models
CREATE TABLE IF NOT EXISTS lead_scoring_models (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    model_type VARCHAR(50) NOT NULL, -- rule_based, ml_based, hybrid
    is_active BOOLEAN DEFAULT TRUE,
    version INTEGER DEFAULT 1,
    scoring_rules JSONB, -- Rule definitions
    model_config JSONB, -- ML model configuration
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_lead_scoring_models_tenant ON lead_scoring_models(tenant_id);
CREATE INDEX idx_lead_scoring_models_active ON lead_scoring_models(is_active);

-- Lead scores (historical tracking)
CREATE TABLE IF NOT EXISTS lead_scores (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    customer_id BIGINT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    model_id BIGINT NOT NULL REFERENCES lead_scoring_models(id) ON DELETE CASCADE,
    score DECIMAL(5, 2) NOT NULL,
    score_breakdown JSONB, -- Component scores
    calculated_at TIMESTAMPTZ DEFAULT NOW(),
    factors JSONB -- Contributing factors
);

CREATE INDEX idx_lead_scores_tenant ON lead_scores(tenant_id);
CREATE INDEX idx_lead_scores_customer ON lead_scores(customer_id);
CREATE INDEX idx_lead_scores_calculated ON lead_scores(calculated_at);

-- Predictive analytics results
CREATE TABLE IF NOT EXISTS predictions (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    customer_id BIGINT REFERENCES customers(id) ON DELETE CASCADE,
    account_id BIGINT REFERENCES accounts(id) ON DELETE CASCADE,
    prediction_type VARCHAR(100) NOT NULL, -- churn, conversion, ltv, next_action
    predicted_value DECIMAL(15, 2),
    predicted_probability DECIMAL(5, 4), -- 0-1
    confidence_level DECIMAL(5, 4), -- 0-1
    prediction_date TIMESTAMPTZ NOT NULL,
    expiry_date TIMESTAMPTZ,
    model_version VARCHAR(50),
    features_used JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_predictions_tenant ON predictions(tenant_id);
CREATE INDEX idx_predictions_customer ON predictions(customer_id);
CREATE INDEX idx_predictions_account ON predictions(account_id);
CREATE INDEX idx_predictions_type ON predictions(prediction_type);

-- ============================================================================
-- 3. COHORT ANALYSIS & SEGMENTATION
-- ============================================================================

-- Enhanced segments table with cohort capabilities
ALTER TABLE segments ADD COLUMN IF NOT EXISTS segment_type VARCHAR(50) DEFAULT 'static'; -- static, dynamic, cohort
ALTER TABLE segments ADD COLUMN IF NOT EXISTS cohort_date_field VARCHAR(100); -- Field to use for cohort grouping
ALTER TABLE segments ADD COLUMN IF NOT EXISTS cohort_interval VARCHAR(50); -- daily, weekly, monthly
ALTER TABLE segments ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE;
ALTER TABLE segments ADD COLUMN IF NOT EXISTS last_computed_at TIMESTAMPTZ;

-- Cohort metrics (computed periodically)
CREATE TABLE IF NOT EXISTS cohort_metrics (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    segment_id BIGINT NOT NULL REFERENCES segments(id) ON DELETE CASCADE,
    cohort_period VARCHAR(50) NOT NULL, -- 2024-W01, 2024-01, 2024-Q1
    period_offset INTEGER NOT NULL, -- Days/weeks/months since cohort start
    metric_name VARCHAR(100) NOT NULL, -- retention_rate, conversion_rate, revenue, etc.
    metric_value DECIMAL(15, 4),
    customer_count INTEGER,
    computed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_cohort_metrics_tenant ON cohort_metrics(tenant_id);
CREATE INDEX idx_cohort_metrics_segment ON cohort_metrics(segment_id);
CREATE INDEX idx_cohort_metrics_period ON cohort_metrics(cohort_period);

-- ============================================================================
-- 4. REAL-TIME DATA STREAMING & EVENTS
-- ============================================================================

-- Event stream for real-time processing
CREATE TABLE IF NOT EXISTS event_stream (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    event_type VARCHAR(100) NOT NULL,
    event_timestamp TIMESTAMPTZ NOT NULL,
    customer_id BIGINT REFERENCES customers(id) ON DELETE CASCADE,
    account_id BIGINT REFERENCES accounts(id) ON DELETE CASCADE,
    session_id VARCHAR(255),
    event_data JSONB NOT NULL,
    processed BOOLEAN DEFAULT FALSE,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_event_stream_tenant ON event_stream(tenant_id);
CREATE INDEX idx_event_stream_timestamp ON event_stream(event_timestamp);
CREATE INDEX idx_event_stream_processed ON event_stream(processed);
CREATE INDEX idx_event_stream_customer ON event_stream(customer_id);
CREATE INDEX idx_event_stream_account ON event_stream(account_id);

-- Real-time alerts and triggers
CREATE TABLE IF NOT EXISTS alerts (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    alert_type VARCHAR(100) NOT NULL,
    severity VARCHAR(20) NOT NULL, -- info, warning, critical
    title VARCHAR(255) NOT NULL,
    description TEXT,
    entity_type VARCHAR(50), -- customer, account, campaign, etc.
    entity_id BIGINT,
    triggered_at TIMESTAMPTZ DEFAULT NOW(),
    acknowledged BOOLEAN DEFAULT FALSE,
    acknowledged_by VARCHAR(255),
    acknowledged_at TIMESTAMPTZ,
    resolved BOOLEAN DEFAULT FALSE,
    resolved_at TIMESTAMPTZ,
    metadata JSONB
);

CREATE INDEX idx_alerts_tenant ON alerts(tenant_id);
CREATE INDEX idx_alerts_triggered ON alerts(triggered_at);
CREATE INDEX idx_alerts_acknowledged ON alerts(acknowledged);

-- ============================================================================
-- 5. FRAUD DETECTION & DATA QUALITY
-- ============================================================================

-- Fraud detection rules
CREATE TABLE IF NOT EXISTS fraud_detection_rules (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    rule_name VARCHAR(255) NOT NULL,
    rule_type VARCHAR(50) NOT NULL, -- click_fraud, lead_fraud, attribution_gaming
    is_active BOOLEAN DEFAULT TRUE,
    severity VARCHAR(20) NOT NULL, -- low, medium, high
    detection_logic JSONB NOT NULL,
    action VARCHAR(50) NOT NULL, -- flag, block, review
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_fraud_rules_tenant ON fraud_detection_rules(tenant_id);
CREATE INDEX idx_fraud_rules_active ON fraud_detection_rules(is_active);

-- Fraud incidents
CREATE TABLE IF NOT EXISTS fraud_incidents (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    rule_id BIGINT REFERENCES fraud_detection_rules(id) ON DELETE SET NULL,
    incident_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    entity_type VARCHAR(50), -- interaction, conversion, customer, etc.
    entity_id BIGINT,
    detected_at TIMESTAMPTZ DEFAULT NOW(),
    confidence_score DECIMAL(5, 4), -- 0-1
    status VARCHAR(50) DEFAULT 'pending', -- pending, confirmed, false_positive, resolved
    reviewed_by VARCHAR(255),
    reviewed_at TIMESTAMPTZ,
    evidence JSONB,
    actions_taken JSONB
);

CREATE INDEX idx_fraud_incidents_tenant ON fraud_incidents(tenant_id);
CREATE INDEX idx_fraud_incidents_detected ON fraud_incidents(detected_at);
CREATE INDEX idx_fraud_incidents_status ON fraud_incidents(status);

-- Data quality scores
CREATE TABLE IF NOT EXISTS data_quality_scores (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL, -- customer, interaction, conversion, account
    entity_id BIGINT NOT NULL,
    completeness_score DECIMAL(5, 2), -- 0-100
    accuracy_score DECIMAL(5, 2), -- 0-100
    consistency_score DECIMAL(5, 2), -- 0-100
    timeliness_score DECIMAL(5, 2), -- 0-100
    overall_score DECIMAL(5, 2), -- 0-100
    issues JSONB, -- List of data quality issues
    calculated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_data_quality_tenant ON data_quality_scores(tenant_id);
CREATE INDEX idx_data_quality_entity ON data_quality_scores(entity_type, entity_id);

-- ============================================================================
-- 6. INTEGRATIONS & EXTERNAL PLATFORMS
-- ============================================================================

-- Integration connections
CREATE TABLE IF NOT EXISTS integrations (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    platform VARCHAR(100) NOT NULL, -- salesforce, hubspot, marketo, google_ads, linkedin, etc.
    integration_type VARCHAR(50) NOT NULL, -- crm, ad_platform, analytics, email, etc.
    is_active BOOLEAN DEFAULT TRUE,
    credentials JSONB, -- Encrypted credentials
    sync_config JSONB, -- Sync settings
    last_sync_at TIMESTAMPTZ,
    last_sync_status VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_integrations_tenant ON integrations(tenant_id);
CREATE INDEX idx_integrations_platform ON integrations(platform);
CREATE INDEX idx_integrations_active ON integrations(is_active);

-- Sync logs
CREATE TABLE IF NOT EXISTS sync_logs (
    id BIGSERIAL PRIMARY KEY,
    integration_id BIGINT NOT NULL REFERENCES integrations(id) ON DELETE CASCADE,
    sync_started_at TIMESTAMPTZ NOT NULL,
    sync_completed_at TIMESTAMPTZ,
    status VARCHAR(50) NOT NULL, -- running, completed, failed
    records_processed INTEGER DEFAULT 0,
    records_created INTEGER DEFAULT 0,
    records_updated INTEGER DEFAULT 0,
    records_failed INTEGER DEFAULT 0,
    error_details JSONB
);

CREATE INDEX idx_sync_logs_integration ON sync_logs(integration_id);
CREATE INDEX idx_sync_logs_started ON sync_logs(sync_started_at);

-- ============================================================================
-- 7. USER BEHAVIOR ANALYTICS
-- ============================================================================

-- Page views and sessions
CREATE TABLE IF NOT EXISTS page_views (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    session_id VARCHAR(255) NOT NULL,
    customer_id BIGINT REFERENCES customers(id) ON DELETE CASCADE,
    page_url TEXT NOT NULL,
    page_title VARCHAR(500),
    referrer TEXT,
    view_timestamp TIMESTAMPTZ NOT NULL,
    time_on_page INTEGER, -- seconds
    scroll_depth INTEGER, -- percentage
    exit_page BOOLEAN DEFAULT FALSE,
    device_type VARCHAR(50),
    browser VARCHAR(100),
    os VARCHAR(100),
    location JSONB, -- city, country, region
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_page_views_tenant ON page_views(tenant_id);
CREATE INDEX idx_page_views_session ON page_views(session_id);
CREATE INDEX idx_page_views_customer ON page_views(customer_id);
CREATE INDEX idx_page_views_timestamp ON page_views(view_timestamp);

-- User actions/events
CREATE TABLE IF NOT EXISTS user_actions (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    session_id VARCHAR(255) NOT NULL,
    customer_id BIGINT REFERENCES customers(id) ON DELETE CASCADE,
    action_type VARCHAR(100) NOT NULL, -- click, form_submit, video_play, etc.
    action_target TEXT, -- Element clicked, form name, etc.
    action_timestamp TIMESTAMPTZ NOT NULL,
    page_url TEXT,
    action_data JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_user_actions_tenant ON user_actions(tenant_id);
CREATE INDEX idx_user_actions_session ON user_actions(session_id);
CREATE INDEX idx_user_actions_customer ON user_actions(customer_id);
CREATE INDEX idx_user_actions_timestamp ON user_actions(action_timestamp);

-- Session summaries
CREATE TABLE IF NOT EXISTS sessions (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    session_id VARCHAR(255) NOT NULL UNIQUE,
    customer_id BIGINT REFERENCES customers(id) ON DELETE CASCADE,
    account_id BIGINT REFERENCES accounts(id) ON DELETE CASCADE,
    session_start TIMESTAMPTZ NOT NULL,
    session_end TIMESTAMPTZ,
    duration INTEGER, -- seconds
    page_views_count INTEGER DEFAULT 0,
    actions_count INTEGER DEFAULT 0,
    entry_page TEXT,
    exit_page TEXT,
    utm_source VARCHAR(255),
    utm_medium VARCHAR(255),
    utm_campaign VARCHAR(255),
    utm_term VARCHAR(255),
    utm_content VARCHAR(255),
    device_type VARCHAR(50),
    browser VARCHAR(100),
    os VARCHAR(100),
    location JSONB,
    converted BOOLEAN DEFAULT FALSE,
    conversion_value DECIMAL(15, 2),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_sessions_tenant ON sessions(tenant_id);
CREATE INDEX idx_sessions_customer ON sessions(customer_id);
CREATE INDEX idx_sessions_start ON sessions(session_start);
CREATE INDEX idx_sessions_session_id ON sessions(session_id);

-- ============================================================================
-- 8. MARKETING MIX MODELING (MMM)
-- ============================================================================

-- MMM models
CREATE TABLE IF NOT EXISTS mmm_models (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    model_name VARCHAR(255) NOT NULL,
    time_period_start DATE NOT NULL,
    time_period_end DATE NOT NULL,
    granularity VARCHAR(20) NOT NULL, -- daily, weekly, monthly
    target_metric VARCHAR(100) NOT NULL, -- revenue, conversions, leads
    model_config JSONB, -- Model parameters
    model_results JSONB, -- Coefficients, R-squared, etc.
    status VARCHAR(50) DEFAULT 'pending', -- pending, training, completed, failed
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_mmm_models_tenant ON mmm_models(tenant_id);
CREATE INDEX idx_mmm_models_status ON mmm_models(status);

-- Channel effectiveness from MMM
CREATE TABLE IF NOT EXISTS channel_effectiveness (
    id BIGSERIAL PRIMARY KEY,
    mmm_model_id BIGINT NOT NULL REFERENCES mmm_models(id) ON DELETE CASCADE,
    channel_name VARCHAR(255) NOT NULL,
    contribution_percentage DECIMAL(5, 2), -- % of total outcome
    roi DECIMAL(10, 2),
    coefficient DECIMAL(10, 4),
    confidence_interval_lower DECIMAL(10, 4),
    confidence_interval_upper DECIMAL(10, 4),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_channel_effectiveness_model ON channel_effectiveness(mmm_model_id);

-- ============================================================================
-- 9. CUSTOM REPORTS & SAVED QUERIES
-- ============================================================================

-- Saved reports
CREATE TABLE IF NOT EXISTS saved_reports (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    report_name VARCHAR(255) NOT NULL,
    report_type VARCHAR(100) NOT NULL, -- attribution, funnel, cohort, custom, etc.
    description TEXT,
    query_config JSONB NOT NULL, -- Filters, dimensions, metrics
    visualization_config JSONB, -- Chart types, settings
    schedule VARCHAR(50), -- For automated reports: daily, weekly, monthly
    is_public BOOLEAN DEFAULT FALSE,
    created_by VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_saved_reports_tenant ON saved_reports(tenant_id);
CREATE INDEX idx_saved_reports_type ON saved_reports(report_type);

-- Report snapshots (for historical tracking)
CREATE TABLE IF NOT EXISTS report_snapshots (
    id BIGSERIAL PRIMARY KEY,
    report_id BIGINT NOT NULL REFERENCES saved_reports(id) ON DELETE CASCADE,
    snapshot_date TIMESTAMPTZ NOT NULL,
    result_data JSONB NOT NULL,
    generated_by VARCHAR(50) DEFAULT 'system', -- system, user
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_report_snapshots_report ON report_snapshots(report_id);
CREATE INDEX idx_report_snapshots_date ON report_snapshots(snapshot_date);

-- ============================================================================
-- 10. ADVANCED ANALYTICS FEATURES
-- ============================================================================

-- A/B Tests and Experiments
CREATE TABLE IF NOT EXISTS experiments (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    experiment_name VARCHAR(255) NOT NULL,
    experiment_type VARCHAR(50) NOT NULL, -- ab_test, multivariate, etc.
    hypothesis TEXT,
    start_date DATE NOT NULL,
    end_date DATE,
    status VARCHAR(50) DEFAULT 'draft', -- draft, running, completed, paused
    variants JSONB NOT NULL, -- Variant definitions
    success_metrics JSONB NOT NULL, -- Metrics to track
    results JSONB, -- Statistical results
    winner VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_experiments_tenant ON experiments(tenant_id);
CREATE INDEX idx_experiments_status ON experiments(status);

-- Experiment assignments
CREATE TABLE IF NOT EXISTS experiment_assignments (
    id BIGSERIAL PRIMARY KEY,
    experiment_id BIGINT NOT NULL REFERENCES experiments(id) ON DELETE CASCADE,
    customer_id BIGINT REFERENCES customers(id) ON DELETE CASCADE,
    session_id VARCHAR(255),
    variant VARCHAR(100) NOT NULL,
    assigned_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_experiment_assignments_experiment ON experiment_assignments(experiment_id);
CREATE INDEX idx_experiment_assignments_customer ON experiment_assignments(customer_id);

-- Feature flags for gradual rollout
CREATE TABLE IF NOT EXISTS feature_flags (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    flag_name VARCHAR(255) NOT NULL,
    flag_key VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    is_enabled BOOLEAN DEFAULT FALSE,
    rollout_percentage INTEGER DEFAULT 0, -- 0-100
    target_segments JSONB, -- Segment IDs to target
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_feature_flags_tenant ON feature_flags(tenant_id);
CREATE INDEX idx_feature_flags_key ON feature_flags(flag_key);

-- ============================================================================
-- 11. ENHANCEMENTS TO EXISTING TABLES
-- ============================================================================

-- Add more fields to interactions
ALTER TABLE interactions ADD COLUMN IF NOT EXISTS session_id VARCHAR(255);
ALTER TABLE interactions ADD COLUMN IF NOT EXISTS page_url TEXT;
ALTER TABLE interactions ADD COLUMN IF NOT EXISTS referrer TEXT;
ALTER TABLE interactions ADD COLUMN IF NOT EXISTS device_type VARCHAR(50);
ALTER TABLE interactions ADD COLUMN IF NOT EXISTS browser VARCHAR(100);
ALTER TABLE interactions ADD COLUMN IF NOT EXISTS os VARCHAR(100);
ALTER TABLE interactions ADD COLUMN IF NOT EXISTS location JSONB;

-- Add more fields to conversion_events
ALTER TABLE conversion_events ADD COLUMN IF NOT EXISTS session_id VARCHAR(255);
ALTER TABLE conversion_events ADD COLUMN IF NOT EXISTS attribution_window INTEGER; -- days
ALTER TABLE conversion_events ADD COLUMN IF NOT EXISTS conversion_path JSONB; -- Full touchpoint path

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_interactions_session ON interactions(session_id);
CREATE INDEX IF NOT EXISTS idx_conversions_session ON conversion_events(session_id);

-- ============================================================================
-- 12. VIEWS FOR COMMON QUERIES
-- ============================================================================

-- Customer 360 view
CREATE OR REPLACE VIEW customer_360 AS
SELECT 
    c.id,
    c.tenant_id,
    c.unified_customer_id,
    c.account_id,
    a.name as account_name,
    c.created_at as customer_since,
    COUNT(DISTINCT i.id) as total_interactions,
    COUNT(DISTINCT ce.id) as total_conversions,
    SUM(ce.revenue) as total_revenue,
    MAX(ls.score) as latest_lead_score,
    MAX(p.predicted_probability) as conversion_probability,
    COUNT(DISTINCT s.session_id) as total_sessions,
    AVG(s.duration) as avg_session_duration
FROM customers c
LEFT JOIN accounts a ON c.account_id = a.id
LEFT JOIN interactions i ON c.id = i.customer_id
LEFT JOIN conversion_events ce ON c.id = ce.customer_id
LEFT JOIN lead_scores ls ON c.id = ls.customer_id
LEFT JOIN predictions p ON c.id = p.customer_id AND p.prediction_type = 'conversion'
LEFT JOIN sessions s ON c.id = s.customer_id
GROUP BY c.id, c.tenant_id, c.unified_customer_id, c.account_id, a.name, c.created_at;

-- Account engagement summary
CREATE OR REPLACE VIEW account_engagement_summary AS
SELECT 
    a.id,
    a.tenant_id,
    a.name,
    a.lifecycle_stage,
    a.health_score,
    a.engagement_score,
    COUNT(DISTINCT c.id) as contacts_count,
    COUNT(DISTINCT ae.id) as engagement_events_count,
    COUNT(DISTINCT s.session_id) as sessions_count,
    SUM(ce.revenue) as total_revenue,
    MAX(s.session_start) as last_activity_date
FROM accounts a
LEFT JOIN customers c ON a.id = c.account_id
LEFT JOIN account_engagements ae ON a.id = ae.account_id
LEFT JOIN sessions s ON a.id = s.account_id
LEFT JOIN conversion_events ce ON c.id = ce.customer_id
GROUP BY a.id, a.tenant_id, a.name, a.lifecycle_stage, a.health_score, a.engagement_score;

-- ============================================================================
-- 13. PRIVACY & COMPLIANCE
-- ============================================================================

-- Data retention policies
CREATE TABLE IF NOT EXISTS data_retention_policies (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL, -- customer, interaction, event, etc.
    retention_period_days INTEGER NOT NULL,
    anonymize_after_days INTEGER, -- Anonymize instead of delete
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_retention_policies_tenant ON data_retention_policies(tenant_id);

-- Consent management
CREATE TABLE IF NOT EXISTS consent_records (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    customer_id BIGINT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    consent_type VARCHAR(100) NOT NULL, -- tracking, marketing, analytics, etc.
    consent_given BOOLEAN NOT NULL,
    consent_date TIMESTAMPTZ NOT NULL,
    consent_source VARCHAR(100), -- website, email, form, etc.
    ip_address INET,
    user_agent TEXT,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_consent_records_tenant ON consent_records(tenant_id);
CREATE INDEX idx_consent_records_customer ON consent_records(customer_id);
CREATE INDEX idx_consent_records_type ON consent_records(consent_type);

-- Data access logs (for GDPR compliance)
CREATE TABLE IF NOT EXISTS data_access_logs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    access_type VARCHAR(50) NOT NULL, -- read, write, delete, export
    entity_type VARCHAR(50) NOT NULL,
    entity_id BIGINT,
    accessed_by VARCHAR(255) NOT NULL,
    access_reason VARCHAR(255),
    ip_address INET,
    accessed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_data_access_logs_tenant ON data_access_logs(tenant_id);
CREATE INDEX idx_data_access_logs_accessed ON data_access_logs(accessed_at);
CREATE INDEX idx_data_access_logs_by ON data_access_logs(accessed_by);

-- ============================================================================
-- MIGRATION COMPLETE
-- ============================================================================

