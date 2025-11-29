-- Add funnel stage tracking to interactions and conversions
-- Based on Factors.ai analytics features

-- Add funnel_stage to interactions
ALTER TABLE interactions ADD COLUMN IF NOT EXISTS funnel_stage VARCHAR(50);
-- Values: MQL, SQL, Opportunity, Closed-Won, etc.

-- Add funnel_stage to conversion_events  
ALTER TABLE conversion_events ADD COLUMN IF NOT EXISTS funnel_stage VARCHAR(50);

-- Create campaigns table for multi-channel tracking
CREATE TABLE IF NOT EXISTS campaigns (
    id SERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    channel_id INT,
    vendor_id INT,
    campaign_type VARCHAR(50), -- 'paid', 'organic', 'email', etc.
    external_campaign_id VARCHAR(255), -- LinkedIn, Google, Meta campaign ID
    platform VARCHAR(50), -- 'linkedin', 'google', 'meta', 'bing'
    start_date DATE,
    end_date DATE,
    budget DECIMAL(18, 4),
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (channel_id) REFERENCES channels(id),
    FOREIGN KEY (vendor_id) REFERENCES vendors(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_campaigns_tenant ON campaigns(tenant_id);
CREATE INDEX IF NOT EXISTS idx_campaigns_platform ON campaigns(platform);

-- Link interactions to campaigns
ALTER TABLE interactions ADD COLUMN IF NOT EXISTS campaign_id INT;
ALTER TABLE interactions ADD CONSTRAINT fk_interactions_campaign 
    FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE SET NULL;

-- Create content/assets table
CREATE TABLE IF NOT EXISTS content_assets (
    id SERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    content_type VARCHAR(50), -- 'whitepaper', 'ebook', 'webinar', 'case_study', etc.
    url TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_content_assets_tenant ON content_assets(tenant_id);

-- Track content engagement
CREATE TABLE IF NOT EXISTS content_engagements (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    content_id INT NOT NULL,
    interaction_id BIGINT,
    engagement_type VARCHAR(50), -- 'view', 'download', 'click', 'share'
    engaged_at TIMESTAMP NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY (content_id) REFERENCES content_assets(id) ON DELETE CASCADE,
    FOREIGN KEY (interaction_id) REFERENCES interactions(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_content_engagements_tenant ON content_engagements(tenant_id);
CREATE INDEX IF NOT EXISTS idx_content_engagements_customer ON content_engagements(customer_id);
CREATE INDEX IF NOT EXISTS idx_content_engagements_content ON content_engagements(content_id);

-- Buyer intent signals table
CREATE TABLE IF NOT EXISTS buyer_intent_signals (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    signal_type VARCHAR(50) NOT NULL, -- 'g2_profile_view', 'competitor_visit', 'category_interest', etc.
    signal_source VARCHAR(50), -- 'g2', 'linkedin', 'website', etc.
    signal_value JSONB,
    detected_at TIMESTAMP NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_buyer_intent_tenant ON buyer_intent_signals(tenant_id);
CREATE INDEX IF NOT EXISTS idx_buyer_intent_customer ON buyer_intent_signals(customer_id);
CREATE INDEX IF NOT EXISTS idx_buyer_intent_type ON buyer_intent_signals(signal_type);

-- Ad spend tracking for multi-channel ROI
CREATE TABLE IF NOT EXISTS ad_spend (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    campaign_id INT,
    platform VARCHAR(50) NOT NULL, -- 'linkedin', 'google', 'meta', 'bing'
    date DATE NOT NULL,
    spend_amount DECIMAL(18, 4) NOT NULL,
    impressions BIGINT,
    clicks BIGINT,
    conversions BIGINT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (campaign_id) REFERENCES campaigns(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_ad_spend_tenant ON ad_spend(tenant_id);
CREATE INDEX IF NOT EXISTS idx_ad_spend_campaign ON ad_spend(campaign_id);
CREATE INDEX IF NOT EXISTS idx_ad_spend_date ON ad_spend(date);

-- View-through attribution tracking
ALTER TABLE interactions ADD COLUMN IF NOT EXISTS is_view_through BOOLEAN DEFAULT FALSE;
ALTER TABLE interactions ADD COLUMN IF NOT EXISTS ad_viewed_at TIMESTAMP;
ALTER TABLE interactions ADD COLUMN IF NOT EXISTS ad_platform VARCHAR(50);

-- Add segments table for advanced segmentation
CREATE TABLE IF NOT EXISTS segments (
    id SERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    criteria JSONB, -- JSON criteria for segment definition
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_segments_tenant ON segments(tenant_id);

-- Link customers to segments
CREATE TABLE IF NOT EXISTS customer_segments (
    customer_id BIGINT NOT NULL,
    segment_id INT NOT NULL,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (customer_id, segment_id),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY (segment_id) REFERENCES segments(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_customer_segments_customer ON customer_segments(customer_id);
CREATE INDEX IF NOT EXISTS idx_customer_segments_segment ON customer_segments(segment_id);

