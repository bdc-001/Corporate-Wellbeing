-- ============================================
-- Convin Revenue Attribution Engine - Schema
-- ============================================

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS attribution_results CASCADE;
DROP TABLE IF EXISTS attribution_runs CASCADE;
DROP TABLE IF EXISTS conversion_events CASCADE;
DROP TABLE IF EXISTS interaction_participants CASCADE;
DROP TABLE IF EXISTS interactions CASCADE;
DROP TABLE IF EXISTS customer_identifiers CASCADE;
DROP TABLE IF EXISTS customers CASCADE;
DROP TABLE IF EXISTS agents CASCADE;
DROP TABLE IF EXISTS teams CASCADE;
DROP TABLE IF EXISTS vendors CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS attribution_models CASCADE;
DROP TABLE IF EXISTS channels CASCADE;
DROP TABLE IF EXISTS event_sources CASCADE;
DROP TABLE IF EXISTS currencies CASCADE;
DROP TABLE IF EXISTS tenants CASCADE;

-- ============================================
-- TENANTS
-- ============================================
CREATE TABLE tenants (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100) UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tenants_code ON tenants(code);

-- ============================================
-- LOOKUP TABLES
-- ============================================

-- Currencies
CREATE TABLE currencies (
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Channels: Voice, WhatsApp, Email, In-App, Webchat, etc.
CREATE TABLE channels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_channels_name ON channels(name);

-- Event Sources: CRM, Billing, OMS, etc.
CREATE TABLE event_sources (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    type VARCHAR(50),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_event_sources_type ON event_sources(type);

-- Attribution Models: FIRST, LAST, LINEAR, TIME_DECAY, AI_WEIGHTED
CREATE TABLE attribution_models (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    params JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Products / Plans (optional, can be linked from external system)
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    external_product_id VARCHAR(255) UNIQUE,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100),
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- ORGANIZATIONAL STRUCTURE
-- ============================================

-- Vendors / BPOs
CREATE TABLE vendors (
    id SERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    UNIQUE (tenant_id, code)
);

CREATE INDEX idx_vendors_tenant ON vendors(tenant_id);
CREATE INDEX idx_vendors_name ON vendors(name);
CREATE INDEX idx_vendors_active ON vendors(is_active);

-- Teams (within vendors or internal)
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    vendor_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    manager_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vendor_id) REFERENCES vendors(id) ON DELETE CASCADE
);

CREATE INDEX idx_teams_vendor ON teams(vendor_id);
CREATE INDEX idx_teams_name ON teams(name);

-- Agents
CREATE TABLE agents (
    id SERIAL PRIMARY KEY,
    vendor_id INT NOT NULL,
    team_id INT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    external_agent_id VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vendor_id) REFERENCES vendors(id) ON DELETE CASCADE,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE SET NULL
);

CREATE INDEX idx_agents_vendor ON agents(vendor_id);
CREATE INDEX idx_agents_team ON agents(team_id);
CREATE INDEX idx_agents_external ON agents(external_agent_id);
CREATE INDEX idx_agents_active ON agents(is_active);

-- Add foreign key constraint for team manager
ALTER TABLE teams
ADD CONSTRAINT fk_teams_manager
FOREIGN KEY (manager_id) REFERENCES agents(id) ON DELETE SET NULL;

-- ============================================
-- CUSTOMER & IDENTITY GRAPH
-- ============================================

-- Master Customers
CREATE TABLE customers (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);

CREATE INDEX idx_customers_tenant ON customers(tenant_id);

-- Customer Identifiers (phone, email, CRM ID, etc.)
CREATE TABLE customer_identifiers (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    type VARCHAR(50) NOT NULL,
    value VARCHAR(255) NOT NULL,
    source_system VARCHAR(100),
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    UNIQUE (customer_id, type, value)
);

CREATE INDEX idx_customer_identifiers_type_value ON customer_identifiers(type, value);
CREATE INDEX idx_customer_identifiers_source ON customer_identifiers(source_system);

-- ============================================
-- INTERACTIONS (CALLS, CHATS, ETC.)
-- ============================================

CREATE TABLE interactions (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    customer_id BIGINT,
    external_interaction_id VARCHAR(255) NOT NULL,
    channel_id INT NOT NULL,
    vendor_id INT,
    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP,
    duration_seconds INT,
    direction VARCHAR(20),
    language VARCHAR(20),
    transcript_location TEXT,
    primary_intent VARCHAR(100),
    secondary_intents JSONB,
    outcome_prediction VARCHAR(100),
    purchase_probability DECIMAL(5, 4),
    raw_metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL,
    FOREIGN KEY (channel_id) REFERENCES channels(id),
    FOREIGN KEY (vendor_id) REFERENCES vendors(id) ON DELETE SET NULL,
    UNIQUE (tenant_id, external_interaction_id)
);

CREATE INDEX idx_interactions_tenant ON interactions(tenant_id);
CREATE INDEX idx_interactions_customer ON interactions(customer_id);
CREATE INDEX idx_interactions_channel_time ON interactions(channel_id, started_at);
CREATE INDEX idx_interactions_vendor_time ON interactions(vendor_id, started_at);
CREATE INDEX idx_interactions_started_at ON interactions(started_at);

-- Who participated in the interaction? (Agent, Bot, System)
CREATE TABLE interaction_participants (
    id BIGSERIAL PRIMARY KEY,
    interaction_id BIGINT NOT NULL,
    participant_type VARCHAR(50) NOT NULL,
    agent_id INT,
    role VARCHAR(50),
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (interaction_id) REFERENCES interactions(id) ON DELETE CASCADE,
    FOREIGN KEY (agent_id) REFERENCES agents(id) ON DELETE SET NULL
);

CREATE INDEX idx_interaction_participants_interaction ON interaction_participants(interaction_id);
CREATE INDEX idx_interaction_participants_agent ON interaction_participants(agent_id);

-- ============================================
-- CONVERSION EVENTS (CRM / BILLING / OMS)
-- ============================================

CREATE TABLE conversion_events (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    event_source_id INT NOT NULL,
    external_event_id VARCHAR(255) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    product_id INT,
    currency_id INT NOT NULL,
    amount_decimal DECIMAL(18, 4) NOT NULL,
    occurred_at TIMESTAMP NOT NULL,
    raw_payload JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY (event_source_id) REFERENCES event_sources(id),
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL,
    FOREIGN KEY (currency_id) REFERENCES currencies(id),
    UNIQUE (tenant_id, event_source_id, external_event_id)
);

CREATE INDEX idx_conversion_events_tenant_time ON conversion_events(tenant_id, occurred_at);
CREATE INDEX idx_conversion_events_customer_time ON conversion_events(customer_id, occurred_at);
CREATE INDEX idx_conversion_events_type ON conversion_events(event_type);

-- ============================================
-- ATTRIBUTION RUNS & RESULTS
-- ============================================

-- Attribution Runs: represents a batch or streaming config execution
CREATE TABLE attribution_runs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    model_id INT NOT NULL,
    name VARCHAR(255),
    description TEXT,
    config JSONB,
    status VARCHAR(50) DEFAULT 'pending',
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (model_id) REFERENCES attribution_models(id)
);

CREATE INDEX idx_attribution_runs_tenant ON attribution_runs(tenant_id);
CREATE INDEX idx_attribution_runs_status ON attribution_runs(status);

-- Attribution Results: 1..N rows per conversion event
CREATE TABLE attribution_results (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    attribution_run_id BIGINT NOT NULL,
    conversion_event_id BIGINT NOT NULL,
    interaction_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    agent_id INT,
    team_id INT,
    vendor_id INT,
    model_id INT NOT NULL,
    attribution_weight DECIMAL(10, 6) NOT NULL,
    attributed_amount DECIMAL(18, 4) NOT NULL,
    is_primary_touch BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (attribution_run_id) REFERENCES attribution_runs(id) ON DELETE CASCADE,
    FOREIGN KEY (conversion_event_id) REFERENCES conversion_events(id) ON DELETE CASCADE,
    FOREIGN KEY (interaction_id) REFERENCES interactions(id) ON DELETE CASCADE,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY (agent_id) REFERENCES agents(id) ON DELETE SET NULL,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE SET NULL,
    FOREIGN KEY (vendor_id) REFERENCES vendors(id) ON DELETE SET NULL,
    FOREIGN KEY (model_id) REFERENCES attribution_models(id)
);

CREATE INDEX idx_attribution_results_tenant ON attribution_results(tenant_id);
CREATE INDEX idx_attribution_results_conversion ON attribution_results(conversion_event_id);
CREATE INDEX idx_attribution_results_interaction ON attribution_results(interaction_id);
CREATE INDEX idx_attribution_results_agent ON attribution_results(agent_id);
CREATE INDEX idx_attribution_results_vendor ON attribution_results(vendor_id);
CREATE INDEX idx_attribution_results_model ON attribution_results(model_id);

-- ============================================
-- SEED DATA
-- ============================================

-- Insert default currencies
INSERT INTO currencies (code, name) VALUES
    ('USD', 'US Dollar'),
    ('INR', 'Indian Rupee'),
    ('EUR', 'Euro'),
    ('GBP', 'British Pound');

-- Insert default channels
INSERT INTO channels (name, description) VALUES
    ('Voice', 'Phone calls'),
    ('WhatsApp', 'WhatsApp messages'),
    ('Email', 'Email interactions'),
    ('In-App', 'In-app chat'),
    ('Webchat', 'Web chat widget'),
    ('SMS', 'SMS messages');

-- Insert default attribution models
INSERT INTO attribution_models (code, name, description) VALUES
    ('FIRST_TOUCH', 'First Touch', 'Attributes 100% to the first interaction'),
    ('LAST_TOUCH', 'Last Touch', 'Attributes 100% to the last interaction'),
    ('LINEAR', 'Linear', 'Distributes attribution equally across all touches'),
    ('TIME_DECAY', 'Time Decay', 'Gives more weight to recent interactions'),
    ('AI_WEIGHTED', 'AI Weighted', 'Uses ML to weight interactions based on intent and outcome');

-- Insert default event sources
INSERT INTO event_sources (name, type, description) VALUES
    ('Salesforce', 'crm', 'Salesforce CRM'),
    ('Stripe', 'billing', 'Stripe payment processor'),
    ('Shopify', 'oms', 'Shopify order management'),
    ('HubSpot', 'crm', 'HubSpot CRM'),
    ('Zendesk', 'ticketing', 'Zendesk ticketing system');

