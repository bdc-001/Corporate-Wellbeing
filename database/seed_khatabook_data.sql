-- ============================================================================
-- KHATABOOK B2B COMPANY - COMPREHENSIVE SEED DATA
-- This script populates the database with realistic Khatabook-specific data
-- Khatabook is a B2B fintech company providing accounting software to merchants
-- ============================================================================

-- Clear existing data (optional - comment out if you want to keep existing data)
-- Note: teams and agents are NOT truncated - they should be managed through the UI
TRUNCATE TABLE attribution_results, attribution_runs, conversion_events, interaction_participants,
  interactions, customer_identifiers, customers, account_engagements, accounts, lead_scores,
  predictions, customer_segments, cohort_metrics, event_stream, alerts, fraud_incidents,
  data_quality_scores, page_views, user_actions, sessions, sync_logs, saved_reports,
  report_snapshots, experiments, experiment_assignments, feature_flags, campaigns,
  content_assets, content_engagements, buyer_intent_signals, ad_spend, mmm_models,
  incrementality_tests RESTART IDENTITY CASCADE;

-- ============================================================================
-- 1. TENANT - KHATABOOK
-- ============================================================================
INSERT INTO tenants (name, code, is_active, created_at) VALUES
('Khatabook', 'khatabook', true, NOW() - INTERVAL '3 years')
ON CONFLICT (code) DO UPDATE SET name = 'Khatabook', is_active = true;

-- Get tenant ID
DO $$
DECLARE
    v_tenant_id BIGINT;
    v_vendor_whatsapp INT;
    v_vendor_phone INT;
    v_vendor_email INT;
    v_vendor_zoom INT;
    v_event_source_id INT := 1;
    v_currency_inr_id INT := 2;
    v_cust_sharma_1 BIGINT;
    v_cust_sharma_2 BIGINT;
    v_cust_patel_1 BIGINT;
    v_cust_kumar_1 BIGINT;
    v_cust_kumar_2 BIGINT;
    v_cust_reddy_1 BIGINT;
    v_cust_nair_1 BIGINT;
    v_cust_das_1 BIGINT;
    v_cust_patel_auto_1 BIGINT;
    v_cust_singh_1 BIGINT;
BEGIN
    SELECT id INTO v_tenant_id FROM tenants WHERE code = 'khatabook';
    
    -- ============================================================================
    -- 2. VENDORS (Communication Platforms)
    -- ============================================================================
    INSERT INTO vendors (tenant_id, name, code, is_active, created_at) VALUES
    (v_tenant_id, 'WhatsApp Business', 'whatsapp', true, NOW() - INTERVAL '2 years'),
    (v_tenant_id, 'Phone (Direct)', 'phone', true, NOW() - INTERVAL '2 years'),
    (v_tenant_id, 'Email', 'email', true, NOW() - INTERVAL '2 years'),
    (v_tenant_id, 'Zoom', 'zoom', true, NOW() - INTERVAL '1 year'),
    (v_tenant_id, 'Google Meet', 'google_meet', true, NOW() - INTERVAL '1 year'),
    (v_tenant_id, 'In-App Chat', 'inapp_chat', true, NOW() - INTERVAL '1 year')
    ON CONFLICT (tenant_id, code) DO NOTHING;

    -- Get vendor IDs for teams
    SELECT id INTO v_vendor_whatsapp FROM vendors WHERE tenant_id = v_tenant_id AND code = 'whatsapp' LIMIT 1;
    SELECT id INTO v_vendor_phone FROM vendors WHERE tenant_id = v_tenant_id AND code = 'phone' LIMIT 1;
    SELECT id INTO v_vendor_email FROM vendors WHERE tenant_id = v_tenant_id AND code = 'email' LIMIT 1;
    SELECT id INTO v_vendor_zoom FROM vendors WHERE tenant_id = v_tenant_id AND code = 'zoom' LIMIT 1;

    -- ============================================================================
    -- 3. TEAMS & AGENTS (Sales Team for Merchant Acquisition)
    -- Teams are not created by default - they should be created through the Team Manager UI
    -- ============================================================================
    -- Teams creation commented out - no default teams
    -- INSERT INTO teams (vendor_id, name, created_at) VALUES
    -- (v_vendor_whatsapp, 'Merchant Acquisition - North', NOW() - INTERVAL '2 years'),
    -- (v_vendor_whatsapp, 'Merchant Acquisition - South', NOW() - INTERVAL '2 years'),
    -- (v_vendor_phone, 'Merchant Acquisition - East', NOW() - INTERVAL '2 years'),
    -- (v_vendor_phone, 'Merchant Acquisition - West', NOW() - INTERVAL '2 years'),
    -- (v_vendor_zoom, 'Enterprise Sales', NOW() - INTERVAL '1 year'),
    -- (v_vendor_email, 'Customer Success', NOW() - INTERVAL '2 years')
    -- ON CONFLICT DO NOTHING;

    -- Agents creation commented out - agents require teams to be created first
    -- INSERT INTO agents (vendor_id, team_id, name, email, is_active, created_at) VALUES
    -- (v_vendor_whatsapp, (SELECT id FROM teams WHERE vendor_id = v_vendor_whatsapp AND name = 'Merchant Acquisition - North' LIMIT 1),
    --  'Rajesh Kumar', 'rajesh.kumar@khatabook.com', true, NOW() - INTERVAL '2 years'),
    -- (v_vendor_whatsapp, (SELECT id FROM teams WHERE vendor_id = v_vendor_whatsapp AND name = 'Merchant Acquisition - North' LIMIT 1),
    --  'Priya Sharma', 'priya.sharma@khatabook.com', true, NOW() - INTERVAL '18 months'),
    -- (v_vendor_whatsapp, (SELECT id FROM teams WHERE vendor_id = v_vendor_whatsapp AND name = 'Merchant Acquisition - South' LIMIT 1),
    --  'Arjun Reddy', 'arjun.reddy@khatabook.com', true, NOW() - INTERVAL '2 years'),
    -- (v_vendor_whatsapp, (SELECT id FROM teams WHERE vendor_id = v_vendor_whatsapp AND name = 'Merchant Acquisition - South' LIMIT 1),
    --  'Meera Nair', 'meera.nair@khatabook.com', true, NOW() - INTERVAL '15 months'),
    -- (v_vendor_phone, (SELECT id FROM teams WHERE vendor_id = v_vendor_phone AND name = 'Merchant Acquisition - East' LIMIT 1),
    --  'Amit Das', 'amit.das@khatabook.com', true, NOW() - INTERVAL '20 months'),
    -- (v_vendor_phone, (SELECT id FROM teams WHERE vendor_id = v_vendor_phone AND name = 'Merchant Acquisition - West' LIMIT 1),
    --  'Sneha Patel', 'sneha.patel@khatabook.com', true, NOW() - INTERVAL '16 months'),
    -- (v_vendor_zoom, (SELECT id FROM teams WHERE vendor_id = v_vendor_zoom AND name = 'Enterprise Sales' LIMIT 1),
    --  'Vikram Singh', 'vikram.singh@khatabook.com', true, NOW() - INTERVAL '1 year'),
    -- (v_vendor_email, (SELECT id FROM teams WHERE vendor_id = v_vendor_email AND name = 'Customer Success' LIMIT 1),
    --  'Anjali Mehta', 'anjali.mehta@khatabook.com', true, NOW() - INTERVAL '2 years')
    -- ON CONFLICT DO NOTHING;

    -- ============================================================================
    -- 4. CHANNELS (Marketing Channels) - Channels are global, not tenant-specific
    -- ============================================================================
    INSERT INTO channels (name, description, created_at) VALUES
    ('Google Ads', 'Paid search advertising', NOW() - INTERVAL '2 years'),
    ('Facebook Ads', 'Social media advertising', NOW() - INTERVAL '2 years'),
    ('Instagram Ads', 'Social media advertising', NOW() - INTERVAL '18 months'),
    ('YouTube Ads', 'Video advertising', NOW() - INTERVAL '18 months'),
    ('Organic Search', 'Organic search traffic', NOW() - INTERVAL '2 years'),
    ('Direct Traffic', 'Direct website visits', NOW() - INTERVAL '2 years'),
    ('Email Marketing', 'Email campaigns', NOW() - INTERVAL '2 years'),
    ('WhatsApp Marketing', 'WhatsApp messaging', NOW() - INTERVAL '1 year'),
    ('Referral Program', 'Referral traffic', NOW() - INTERVAL '1 year'),
    ('Content Marketing', 'Content-driven traffic', NOW() - INTERVAL '18 months')
    ON CONFLICT (name) DO NOTHING;

    -- ============================================================================
    -- 5. ABM - ACCOUNTS (Target Merchants/Shops)
    -- ============================================================================
    INSERT INTO accounts (tenant_id, name, domain, industry, company_size, annual_revenue,
        location, target_account, account_tier, health_score, engagement_score, intent_score,
        lifecycle_stage, crm_account_id, owner_id, created_at) VALUES
    (v_tenant_id, 'Sharma Electronics', 'sharmaelectronics.in', 'Retail - Electronics', '10-50', 5000000.00,
        'Delhi, India', true, 'Mid-Market', 85.5, 78.2, 92.0, 'Opportunity', 'CRM-KB-001', 'rajesh.kumar@khatabook.com', NOW() - INTERVAL '60 days'),
    (v_tenant_id, 'Patel Textiles', 'pateltextiles.com', 'Retail - Textiles', '50-100', 12000000.00,
        'Mumbai, India', true, 'Mid-Market', 72.3, 65.8, 71.5, 'SQL', 'CRM-KB-002', 'priya.sharma@khatabook.com', NOW() - INTERVAL '45 days'),
    (v_tenant_id, 'Kumar Grocery Chain', 'kumargrocery.in', 'Retail - Grocery', '100-500', 50000000.00,
        'Bangalore, India', true, 'Enterprise', 90.1, 88.5, 95.2, 'Customer', 'CRM-KB-003', 'vikram.singh@khatabook.com', NOW() - INTERVAL '120 days'),
    (v_tenant_id, 'Reddy Hardware Store', 'reddyhardware.com', 'Retail - Hardware', '10-50', 3000000.00,
        'Hyderabad, India', true, 'SMB', 68.0, 52.3, 60.0, 'MQL', 'CRM-KB-004', 'arjun.reddy@khatabook.com', NOW() - INTERVAL '30 days'),
    (v_tenant_id, 'Nair Pharmacy', 'nairpharmacy.in', 'Retail - Pharmacy', '50-100', 8000000.00,
        'Chennai, India', true, 'Mid-Market', 88.7, 82.1, 89.3, 'Opportunity', 'CRM-KB-005', 'meera.nair@khatabook.com', NOW() - INTERVAL '75 days'),
    (v_tenant_id, 'Das Furniture Mart', 'dasfurniture.com', 'Retail - Furniture', '10-50', 4000000.00,
        'Kolkata, India', true, 'SMB', 65.2, 58.4, 55.0, 'MQL', 'CRM-KB-006', 'amit.das@khatabook.com', NOW() - INTERVAL '25 days'),
    (v_tenant_id, 'Patel Auto Parts', 'patelautoparts.in', 'Retail - Automotive', '50-100', 15000000.00,
        'Ahmedabad, India', true, 'Mid-Market', 75.8, 70.2, 78.5, 'SQL', 'CRM-KB-007', 'sneha.patel@khatabook.com', NOW() - INTERVAL '50 days'),
    (v_tenant_id, 'Singh Mobile Store', 'singhmobile.in', 'Retail - Electronics', '10-50', 6000000.00,
        'Pune, India', true, 'SMB', 70.5, 63.1, 68.0, 'SQL', 'CRM-KB-008', 'rajesh.kumar@khatabook.com', NOW() - INTERVAL '40 days')
    ON CONFLICT DO NOTHING;

    -- ============================================================================
    -- 6. CUSTOMERS (Merchants/Shop Owners)
    -- ============================================================================
    INSERT INTO customers (tenant_id, account_id, created_at) VALUES
    (v_tenant_id, (SELECT id FROM accounts WHERE tenant_id = v_tenant_id AND name = 'Sharma Electronics' LIMIT 1),
     NOW() - INTERVAL '60 days'),
    (v_tenant_id, (SELECT id FROM accounts WHERE tenant_id = v_tenant_id AND name = 'Sharma Electronics' LIMIT 1),
     NOW() - INTERVAL '55 days'),
    (v_tenant_id, (SELECT id FROM accounts WHERE tenant_id = v_tenant_id AND name = 'Patel Textiles' LIMIT 1),
     NOW() - INTERVAL '45 days'),
    (v_tenant_id, (SELECT id FROM accounts WHERE tenant_id = v_tenant_id AND name = 'Kumar Grocery Chain' LIMIT 1),
     NOW() - INTERVAL '120 days'),
    (v_tenant_id, (SELECT id FROM accounts WHERE tenant_id = v_tenant_id AND name = 'Kumar Grocery Chain' LIMIT 1),
     NOW() - INTERVAL '115 days'),
    (v_tenant_id, (SELECT id FROM accounts WHERE tenant_id = v_tenant_id AND name = 'Reddy Hardware Store' LIMIT 1),
     NOW() - INTERVAL '30 days'),
    (v_tenant_id, (SELECT id FROM accounts WHERE tenant_id = v_tenant_id AND name = 'Nair Pharmacy' LIMIT 1),
     NOW() - INTERVAL '75 days'),
    (v_tenant_id, (SELECT id FROM accounts WHERE tenant_id = v_tenant_id AND name = 'Das Furniture Mart' LIMIT 1),
     NOW() - INTERVAL '25 days'),
    (v_tenant_id, (SELECT id FROM accounts WHERE tenant_id = v_tenant_id AND name = 'Patel Auto Parts' LIMIT 1),
     NOW() - INTERVAL '50 days'),
    (v_tenant_id, (SELECT id FROM accounts WHERE tenant_id = v_tenant_id AND name = 'Singh Mobile Store' LIMIT 1),
     NOW() - INTERVAL '40 days')
    ON CONFLICT DO NOTHING;

    -- Get customer IDs (using ORDER BY to match insertion order)
    SELECT id INTO v_cust_sharma_1 FROM customers WHERE tenant_id = v_tenant_id AND account_id = (SELECT id FROM accounts WHERE tenant_id = v_tenant_id AND name = 'Sharma Electronics' LIMIT 1) ORDER BY created_at LIMIT 1 OFFSET 0;
    SELECT id INTO v_cust_sharma_2 FROM customers WHERE tenant_id = v_tenant_id AND account_id = (SELECT id FROM accounts WHERE tenant_id = v_tenant_id AND name = 'Sharma Electronics' LIMIT 1) ORDER BY created_at LIMIT 1 OFFSET 1;
    SELECT id INTO v_cust_patel_1 FROM customers WHERE tenant_id = v_tenant_id AND account_id = (SELECT id FROM accounts WHERE tenant_id = v_tenant_id AND name = 'Patel Textiles' LIMIT 1) ORDER BY created_at LIMIT 1;
    SELECT id INTO v_cust_kumar_1 FROM customers WHERE tenant_id = v_tenant_id AND account_id = (SELECT id FROM accounts WHERE tenant_id = v_tenant_id AND name = 'Kumar Grocery Chain' LIMIT 1) ORDER BY created_at LIMIT 1 OFFSET 0;
    SELECT id INTO v_cust_kumar_2 FROM customers WHERE tenant_id = v_tenant_id AND account_id = (SELECT id FROM accounts WHERE tenant_id = v_tenant_id AND name = 'Kumar Grocery Chain' LIMIT 1) ORDER BY created_at LIMIT 1 OFFSET 1;
    SELECT id INTO v_cust_reddy_1 FROM customers WHERE tenant_id = v_tenant_id AND account_id = (SELECT id FROM accounts WHERE tenant_id = v_tenant_id AND name = 'Reddy Hardware Store' LIMIT 1) ORDER BY created_at LIMIT 1;
    SELECT id INTO v_cust_nair_1 FROM customers WHERE tenant_id = v_tenant_id AND account_id = (SELECT id FROM accounts WHERE tenant_id = v_tenant_id AND name = 'Nair Pharmacy' LIMIT 1) ORDER BY created_at LIMIT 1;
    SELECT id INTO v_cust_das_1 FROM customers WHERE tenant_id = v_tenant_id AND account_id = (SELECT id FROM accounts WHERE tenant_id = v_tenant_id AND name = 'Das Furniture Mart' LIMIT 1) ORDER BY created_at LIMIT 1;
    SELECT id INTO v_cust_patel_auto_1 FROM customers WHERE tenant_id = v_tenant_id AND account_id = (SELECT id FROM accounts WHERE tenant_id = v_tenant_id AND name = 'Patel Auto Parts' LIMIT 1) ORDER BY created_at LIMIT 1;
    SELECT id INTO v_cust_singh_1 FROM customers WHERE tenant_id = v_tenant_id AND account_id = (SELECT id FROM accounts WHERE tenant_id = v_tenant_id AND name = 'Singh Mobile Store' LIMIT 1) ORDER BY created_at LIMIT 1;

    -- Customer Identifiers (Phone numbers, emails)
    INSERT INTO customer_identifiers (customer_id, type, value, is_primary, created_at) VALUES
    (v_cust_sharma_1, 'phone', '+91-98765-43210', true, NOW() - INTERVAL '60 days'),
    (v_cust_sharma_1, 'email', 'sharma.electronics@gmail.com', false, NOW() - INTERVAL '60 days'),
    (v_cust_sharma_2, 'phone', '+91-98765-43211', true, NOW() - INTERVAL '55 days'),
    (v_cust_patel_1, 'phone', '+91-98765-43220', true, NOW() - INTERVAL '45 days'),
    (v_cust_kumar_1, 'phone', '+91-98765-43230', true, NOW() - INTERVAL '120 days'),
    (v_cust_kumar_1, 'email', 'kumar.grocery@business.com', false, NOW() - INTERVAL '120 days'),
    (v_cust_reddy_1, 'phone', '+91-98765-43240', true, NOW() - INTERVAL '30 days'),
    (v_cust_nair_1, 'phone', '+91-98765-43250', true, NOW() - INTERVAL '75 days'),
    (v_cust_das_1, 'phone', '+91-98765-43260', true, NOW() - INTERVAL '25 days'),
    (v_cust_patel_auto_1, 'phone', '+91-98765-43270', true, NOW() - INTERVAL '50 days'),
    (v_cust_singh_1, 'phone', '+91-98765-43280', true, NOW() - INTERVAL '40 days')
    ON CONFLICT DO NOTHING;

    -- ============================================================================
    -- 7. INTERACTIONS (Sales Calls, Demos, Follow-ups)
    -- ============================================================================
    INSERT INTO interactions (tenant_id, customer_id, external_interaction_id, channel_id, started_at, ended_at,
        duration_seconds, primary_intent, secondary_intents, outcome_prediction, purchase_probability,
        raw_metadata, created_at) VALUES
    -- Sharma Electronics interactions
    (v_tenant_id, v_cust_sharma_1,
     'KB-INT-001', (SELECT id FROM channels WHERE name = 'Google Ads' LIMIT 1),
     NOW() - INTERVAL '60 days', NOW() - INTERVAL '60 days' + INTERVAL '3 minutes', 180, 'bookkeeping_need', 
     '["gst_filing", "inventory_management"]'::jsonb, 'high_interest', 0.85,
     '{"page": "pricing", "source": "google_ads"}'::jsonb, NOW() - INTERVAL '60 days'),
    
    (v_tenant_id, v_cust_sharma_1,
     'KB-INT-002', (SELECT id FROM channels WHERE name = 'Direct Traffic' LIMIT 1),
     NOW() - INTERVAL '58 days', NOW() - INTERVAL '58 days' + INTERVAL '20 minutes', 1200, 'bookkeeping_need',
     '["gst_filing"]'::jsonb, 'demo_requested', 0.90,
     '{"call_type": "inbound", "topic": "pricing_inquiry"}'::jsonb, NOW() - INTERVAL '58 days'),
    
    (v_tenant_id, v_cust_sharma_1,
     'KB-INT-003', (SELECT id FROM channels WHERE name = 'Direct Traffic' LIMIT 1),
     NOW() - INTERVAL '55 days', NOW() - INTERVAL '55 days' + INTERVAL '40 minutes', 2400, 'bookkeeping_need',
     '["gst_filing", "inventory_management", "payment_tracking"]'::jsonb, 'demo_completed', 0.95,
     '{"demo_duration": 40, "features_shown": ["gst", "inventory", "payments"]}'::jsonb, NOW() - INTERVAL '55 days'),
    
    -- Patel Textiles interactions
    (v_tenant_id, v_cust_patel_1,
     'KB-INT-004', (SELECT id FROM channels WHERE name = 'Facebook Ads' LIMIT 1),
     NOW() - INTERVAL '45 days', NOW() - INTERVAL '45 days' + INTERVAL '4 minutes', 240, 'gst_filing',
     '["bookkeeping_need"]'::jsonb, 'medium_interest', 0.65,
     '{"page": "features", "source": "facebook_ads"}'::jsonb, NOW() - INTERVAL '45 days'),
    
    (v_tenant_id, v_cust_patel_1,
     'KB-INT-005', (SELECT id FROM channels WHERE name = 'WhatsApp Marketing' LIMIT 1),
     NOW() - INTERVAL '43 days', NOW() - INTERVAL '43 days', 0, 'gst_filing',
     '["bookkeeping_need"]'::jsonb, 'information_requested', 0.70,
     '{"message_type": "inquiry", "response_time": 120}'::jsonb, NOW() - INTERVAL '43 days'),
    
    -- Kumar Grocery Chain interactions (existing customer)
    (v_tenant_id, v_cust_kumar_1,
     'KB-INT-006', (SELECT id FROM channels WHERE name = 'Direct Traffic' LIMIT 1),
     NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days' + INTERVAL '10 minutes', 600, 'payment_tracking',
     '["inventory_management"]'::jsonb, 'support_request', 0.20,
     '{"chat_type": "support", "issue": "payment_reconciliation"}'::jsonb, NOW() - INTERVAL '5 days'),
    
    -- Reddy Hardware Store interactions
    (v_tenant_id, v_cust_reddy_1,
     'KB-INT-007', (SELECT id FROM channels WHERE name = 'Instagram Ads' LIMIT 1),
     NOW() - INTERVAL '30 days', NOW() - INTERVAL '30 days' + INTERVAL '2 minutes', 120, 'inventory_management',
     '["bookkeeping_need"]'::jsonb, 'low_interest', 0.40,
     '{"page": "home", "source": "instagram_ads"}'::jsonb, NOW() - INTERVAL '30 days'),
    
    -- Nair Pharmacy interactions
    (v_tenant_id, v_cust_nair_1,
     'KB-INT-008', (SELECT id FROM channels WHERE name = 'Email Marketing' LIMIT 1),
     NOW() - INTERVAL '75 days', NOW() - INTERVAL '75 days', 0, 'gst_filing',
     '["bookkeeping_need", "payment_tracking"]'::jsonb, 'email_opened', 0.55,
     '{"email_type": "newsletter", "subject": "GST Filing Made Easy"}'::jsonb, NOW() - INTERVAL '75 days'),
    
    (v_tenant_id, v_cust_nair_1,
     'KB-INT-009', (SELECT id FROM channels WHERE name = 'Direct Traffic' LIMIT 1),
     NOW() - INTERVAL '70 days', NOW() - INTERVAL '70 days' + INTERVAL '15 minutes', 900, 'gst_filing',
     '["bookkeeping_need"]'::jsonb, 'demo_scheduled', 0.75,
     '{"call_type": "outbound", "topic": "gst_features"}'::jsonb, NOW() - INTERVAL '70 days'),
    
    -- Das Furniture Mart interactions
    (v_tenant_id, v_cust_das_1,
     'KB-INT-010', (SELECT id FROM channels WHERE name = 'Referral Program' LIMIT 1),
     NOW() - INTERVAL '25 days', NOW() - INTERVAL '25 days' + INTERVAL '5 minutes', 300, 'bookkeeping_need',
     '["inventory_management"]'::jsonb, 'referral_visit', 0.60,
     '{"referrer": "sharma_electronics", "page": "signup"}'::jsonb, NOW() - INTERVAL '25 days'),
    
    -- Patel Auto Parts interactions
    (v_tenant_id, v_cust_patel_auto_1,
     'KB-INT-011', (SELECT id FROM channels WHERE name = 'YouTube Ads' LIMIT 1),
     NOW() - INTERVAL '50 days', NOW() - INTERVAL '50 days' + INTERVAL '3 minutes', 180, 'payment_tracking',
     '["bookkeeping_need"]'::jsonb, 'video_watched', 0.50,
     '{"video_id": "kb-101", "watch_time": 120}'::jsonb, NOW() - INTERVAL '50 days'),
    
    -- Singh Mobile Store interactions
    (v_tenant_id, v_cust_singh_1,
     'KB-INT-012', (SELECT id FROM channels WHERE name = 'Organic Search' LIMIT 1),
     NOW() - INTERVAL '40 days', NOW() - INTERVAL '40 days' + INTERVAL '7 minutes', 420, 'bookkeeping_need',
     '["gst_filing", "inventory_management"]'::jsonb, 'high_interest', 0.80,
     '{"keyword": "best accounting software india", "page": "features"}'::jsonb, NOW() - INTERVAL '40 days')
    ON CONFLICT DO NOTHING;

    -- Interaction Participants (Agents) - Link agents to phone/video calls
    INSERT INTO interaction_participants (interaction_id, participant_type, agent_id, role, created_at)
    SELECT i.id, 'agent', a.id, 'sales_rep', i.created_at
    FROM interactions i
    CROSS JOIN agents a
    WHERE i.tenant_id = v_tenant_id
      AND i.duration_seconds > 600  -- Only longer interactions (calls/demos)
      AND a.vendor_id IN (v_vendor_whatsapp, v_vendor_phone, v_vendor_zoom)
    LIMIT 5
    ON CONFLICT DO NOTHING;

    -- ============================================================================
    -- 8. CONVERSION EVENTS (Merchant Signups, Subscription Purchases)
    -- ============================================================================
    -- Get event source and currency IDs
    SELECT id INTO v_event_source_id FROM event_sources WHERE name = 'Salesforce' LIMIT 1;
    IF v_event_source_id IS NULL THEN
        v_event_source_id := 1;
    END IF;
    
    SELECT id INTO v_currency_inr_id FROM currencies WHERE code = 'INR' LIMIT 1;
    IF v_currency_inr_id IS NULL THEN
        v_currency_inr_id := 2;
    END IF;
    
    INSERT INTO conversion_events (tenant_id, customer_id, event_source_id, external_event_id,
        event_type, currency_id, amount_decimal, occurred_at, funnel_stage, raw_payload, created_at) VALUES
        -- Sharma Electronics - Subscription Purchase
        (v_tenant_id, v_cust_sharma_1,
         v_event_source_id, 'KB-CONV-001', 'subscription_purchase', v_currency_inr_id, 2999.00,
         NOW() - INTERVAL '52 days', 'Closed-Won',
         '{"plan": "premium", "duration": "annual", "payment_method": "upi"}'::jsonb, NOW() - INTERVAL '52 days'),
        
        -- Kumar Grocery Chain - Subscription Purchase (Enterprise)
        (v_tenant_id, v_cust_kumar_1,
         v_event_source_id, 'KB-CONV-002', 'subscription_purchase', v_currency_inr_id, 49999.00,
         NOW() - INTERVAL '110 days', 'Closed-Won',
         '{"plan": "enterprise", "duration": "annual", "stores": 25}'::jsonb, NOW() - INTERVAL '110 days'),
        
        -- Kumar Grocery Chain - Renewal
        (v_tenant_id, v_cust_kumar_1,
         v_event_source_id, 'KB-CONV-003', 'subscription_renewal', v_currency_inr_id, 49999.00,
         NOW() - INTERVAL '10 days', 'Closed-Won',
         '{"plan": "enterprise", "duration": "annual", "renewal": true}'::jsonb, NOW() - INTERVAL '10 days'),
        
        -- Nair Pharmacy - Subscription Purchase
        (v_tenant_id, v_cust_nair_1,
         v_event_source_id, 'KB-CONV-004', 'subscription_purchase', v_currency_inr_id, 4999.00,
         NOW() - INTERVAL '65 days', 'Closed-Won',
         '{"plan": "professional", "duration": "annual", "payment_method": "card"}'::jsonb, NOW() - INTERVAL '65 days'),
        
        -- Singh Mobile Store - Free Trial Signup
        (v_tenant_id, v_cust_singh_1,
         v_event_source_id, 'KB-CONV-005', 'trial_signup', v_currency_inr_id, 0.00,
         NOW() - INTERVAL '35 days', 'SQL',
         '{"trial_duration": 14, "source": "organic_search"}'::jsonb, NOW() - INTERVAL '35 days')
    ON CONFLICT DO NOTHING;

    -- ============================================================================
    -- 9. CAMPAIGNS (Marketing Campaigns)
    -- ============================================================================
    INSERT INTO campaigns (tenant_id, name, campaign_type, channel_id, start_date, end_date, budget,
        platform, created_at) VALUES
    (v_tenant_id, 'Q4 Merchant Acquisition - Google Ads', 'paid_search',
     (SELECT id FROM channels WHERE name = 'Google Ads' LIMIT 1),
     NOW() - INTERVAL '90 days', NOW() - INTERVAL '30 days', 500000.00, 'Google Ads', NOW() - INTERVAL '90 days'),
    (v_tenant_id, 'Q4 Merchant Acquisition - Facebook Ads', 'paid_social',
     (SELECT id FROM channels WHERE name = 'Facebook Ads' LIMIT 1),
     NOW() - INTERVAL '90 days', NOW() - INTERVAL '30 days', 300000.00, 'Facebook', NOW() - INTERVAL '90 days'),
    (v_tenant_id, 'GST Filing Campaign - Instagram', 'paid_social',
     (SELECT id FROM channels WHERE name = 'Instagram Ads' LIMIT 1),
     NOW() - INTERVAL '60 days', NOW() + INTERVAL '30 days', 200000.00, 'Instagram', NOW() - INTERVAL '60 days'),
    (v_tenant_id, 'New Year Referral Program', 'referral',
     (SELECT id FROM channels WHERE name = 'Referral Program' LIMIT 1),
     NOW() - INTERVAL '30 days', NOW() + INTERVAL '60 days', 100000.00, 'Internal', NOW() - INTERVAL '30 days'),
    (v_tenant_id, 'Enterprise Merchant Outreach', 'email',
     (SELECT id FROM channels WHERE name = 'Email Marketing' LIMIT 1),
     NOW() - INTERVAL '45 days', NOW() + INTERVAL '15 days', 50000.00, 'Email', NOW() - INTERVAL '45 days')
    ON CONFLICT DO NOTHING;

    -- ============================================================================
    -- 10. AD SPEND (Marketing Spend by Channel)
    -- ============================================================================
    INSERT INTO ad_spend (tenant_id, campaign_id, platform, date, spend_amount,
        impressions, clicks, conversions, created_at) VALUES
    (v_tenant_id, (SELECT id FROM campaigns WHERE tenant_id = v_tenant_id AND name LIKE '%Google Ads%' LIMIT 1),
     'Google Ads', (NOW() - INTERVAL '7 days')::date, 15000.00, 500000, 12000, 45, NOW() - INTERVAL '7 days'),
    (v_tenant_id, (SELECT id FROM campaigns WHERE tenant_id = v_tenant_id AND name LIKE '%Facebook Ads%' LIMIT 1),
     'Facebook', (NOW() - INTERVAL '7 days')::date, 10000.00, 300000, 8000, 32, NOW() - INTERVAL '7 days'),
    (v_tenant_id, (SELECT id FROM campaigns WHERE tenant_id = v_tenant_id AND name LIKE '%Instagram%' LIMIT 1),
     'Instagram', (NOW() - INTERVAL '7 days')::date, 7000.00, 200000, 5000, 18, NOW() - INTERVAL '7 days'),
    (v_tenant_id, NULL,
     'YouTube', (NOW() - INTERVAL '7 days')::date, 5000.00, 150000, 3000, 12, NOW() - INTERVAL '7 days')
    ON CONFLICT DO NOTHING;

    -- ============================================================================
    -- 11. CONTENT ASSETS (Blog Posts, Videos)
    -- ============================================================================
    INSERT INTO content_assets (tenant_id, name, content_type, url, created_at) VALUES
    (v_tenant_id, 'How to File GST Returns Easily', 'blog_post', 'https://khatabook.com/blog/gst-filing-guide', NOW() - INTERVAL '90 days'),
    (v_tenant_id, 'Inventory Management Best Practices', 'blog_post', 'https://khatabook.com/blog/inventory-management', NOW() - INTERVAL '75 days'),
    (v_tenant_id, 'Khatabook Product Demo Video', 'video', 'https://youtube.com/watch?v=kb-demo-101', NOW() - INTERVAL '60 days'),
    (v_tenant_id, 'Payment Tracking Made Simple', 'blog_post', 'https://khatabook.com/blog/payment-tracking', NOW() - INTERVAL '45 days'),
    (v_tenant_id, 'GST Calculator Tool', 'tool', 'https://khatabook.com/tools/gst-calculator', NOW() - INTERVAL '30 days')
    ON CONFLICT DO NOTHING;

    -- Content Engagements
    INSERT INTO content_engagements (tenant_id, customer_id, content_id, engagement_type,
        engaged_at, created_at)
    SELECT v_tenant_id, c.id, ca.id, 'view', NOW() - INTERVAL '20 days', NOW() - INTERVAL '20 days'
    FROM customers c
    CROSS JOIN content_assets ca
    WHERE c.tenant_id = v_tenant_id
      AND ca.tenant_id = v_tenant_id
    LIMIT 15
    ON CONFLICT DO NOTHING;

    -- ============================================================================
    -- 12. LEAD SCORING
    -- ============================================================================
    INSERT INTO lead_scoring_models (tenant_id, name, model_type, is_active, scoring_rules,
        created_at) VALUES
    (v_tenant_id, 'Default Merchant Scoring Model', 'rule_based', true,
     '{"rules": [{"field": "company_size", "weight": 0.3}, {"field": "engagement_score", "weight": 0.4}, {"field": "intent_score", "weight": 0.3}]}'::jsonb,
     NOW() - INTERVAL '6 months')
    ON CONFLICT DO NOTHING;

    INSERT INTO lead_scores (tenant_id, customer_id, model_id, score, score_breakdown,
        factors, calculated_at)
    SELECT v_tenant_id, c.id, 
           (SELECT id FROM lead_scoring_models WHERE tenant_id = v_tenant_id LIMIT 1),
           CASE 
               WHEN a.company_size LIKE '%100%' THEN 85
               WHEN a.company_size LIKE '%50%' THEN 70
               ELSE 55
           END,
           '{"company_size": 30, "engagement": 40, "intent": 30}'::jsonb,
           '["high_company_size", "active_engagement", "strong_intent"]'::jsonb,
           NOW() - INTERVAL '30 days'
    FROM customers c
    JOIN accounts a ON c.account_id = a.id
    WHERE c.tenant_id = v_tenant_id
    ON CONFLICT DO NOTHING;

    -- ============================================================================
    -- 13. PREDICTIONS
    -- ============================================================================
    INSERT INTO predictions (tenant_id, customer_id, prediction_type, predicted_value,
        predicted_probability, confidence_level, prediction_date, features_used, created_at) VALUES
    (v_tenant_id, v_cust_sharma_1,
     'conversion_probability', 0.95, 0.95, 0.88,
     NOW() - INTERVAL '50 days', '["high_engagement", "demo_completed", "strong_intent"]'::jsonb, NOW() - INTERVAL '50 days'),
    (v_tenant_id, v_cust_patel_1,
     'ltv', 50000.00, NULL, 0.75,
     NOW() - INTERVAL '40 days', '["mid_market_account", "annual_plan"]'::jsonb, NOW() - INTERVAL '40 days'),
    (v_tenant_id, v_cust_reddy_1,
     'conversion_probability', 0.40, 0.40, 0.65,
     NOW() - INTERVAL '25 days', '["low_engagement", "early_stage"]'::jsonb, NOW() - INTERVAL '25 days')
    ON CONFLICT DO NOTHING;

    -- ============================================================================
    -- 14. ACCOUNT ENGAGEMENTS
    -- ============================================================================
    INSERT INTO account_engagements (tenant_id, account_id, engagement_type, engagement_date,
        metadata, created_at)
    SELECT v_tenant_id, a.id, 'website_visit', NOW() - INTERVAL '10 days',
           '{"pages_viewed": 5, "duration": 300}'::jsonb, NOW() - INTERVAL '10 days'
    FROM accounts a
    WHERE a.tenant_id = v_tenant_id
    LIMIT 5
    ON CONFLICT DO NOTHING;

    -- ============================================================================
    -- 15. SEGMENTS & CUSTOMER SEGMENTS
    -- ============================================================================
    INSERT INTO segments (tenant_id, name, segment_type, criteria, created_at) VALUES
    (v_tenant_id, 'High-Value Merchants', 'revenue', '{"min_revenue": 10000000}'::jsonb, NOW() - INTERVAL '3 months'),
    (v_tenant_id, 'GST Filing Intent', 'intent', '{"primary_intent": "gst_filing"}'::jsonb, NOW() - INTERVAL '3 months'),
    (v_tenant_id, 'Enterprise Accounts', 'account_tier', '{"tier": "Enterprise"}'::jsonb, NOW() - INTERVAL '3 months')
    ON CONFLICT DO NOTHING;

    INSERT INTO customer_segments (customer_id, segment_id)
    SELECT c.id, s.id
    FROM customers c
    JOIN accounts a ON c.account_id = a.id
    CROSS JOIN segments s
    WHERE c.tenant_id = v_tenant_id
      AND s.tenant_id = v_tenant_id
      AND (
          (s.name = 'Enterprise Accounts' AND a.account_tier = 'Enterprise')
          OR (s.name = 'GST Filing Intent' AND EXISTS (
              SELECT 1 FROM interactions i WHERE i.customer_id = c.id AND i.secondary_intents::text LIKE '%gst_filing%'
          ))
      )
    ON CONFLICT DO NOTHING;

    -- ============================================================================
    -- 16. COHORT METRICS
    -- ============================================================================
    -- Cohort metrics commented out due to schema differences - can be added later
    -- INSERT INTO cohort_metrics (tenant_id, segment_id, cohort_period, period_offset, metric_name,
    --     metric_value, customer_count, computed_at) VALUES
    -- (v_tenant_id, NULL, '2024-09', 0, 'retention_rate', 0.75, 100, NOW() - INTERVAL '1 month'),
    -- (v_tenant_id, NULL, '2024-10', 0, 'retention_rate', 0.82, 120, NOW() - INTERVAL '1 month'),
    -- (v_tenant_id, NULL, '2024-11', 0, 'retention_rate', 0.88, 150, NOW() - INTERVAL '1 month')
    -- ON CONFLICT DO NOTHING;

    -- ============================================================================
    -- 17. EVENT STREAM (Real-time Events)
    -- ============================================================================
    INSERT INTO event_stream (tenant_id, event_type, event_data, event_timestamp, created_at) VALUES
    (v_tenant_id, 'merchant_signup', '{"merchant_id": "KB-001", "plan": "premium"}'::jsonb, NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours'),
    (v_tenant_id, 'subscription_purchase', '{"merchant_id": "KB-002", "amount": 4999}'::jsonb, NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 hour'),
    (v_tenant_id, 'demo_scheduled', '{"merchant_id": "KB-003", "agent": "rajesh.kumar"}'::jsonb, NOW() - INTERVAL '30 minutes', NOW() - INTERVAL '30 minutes'),
    (v_tenant_id, 'trial_started', '{"merchant_id": "KB-004", "trial_days": 14}'::jsonb, NOW() - INTERVAL '15 minutes', NOW() - INTERVAL '15 minutes')
    ON CONFLICT DO NOTHING;

    -- ============================================================================
    -- 18. ALERTS
    -- ============================================================================
    INSERT INTO alerts (tenant_id, alert_type, severity, title, description, acknowledged, triggered_at) VALUES
    (v_tenant_id, 'threshold_breach', 'high', 'High Conversion Rate Detected',
     'Conversion rate increased by 25% in the last 24 hours', false, NOW() - INTERVAL '2 hours'),
    (v_tenant_id, 'anomaly', 'medium', 'Unusual Traffic Pattern',
     'Traffic from Instagram Ads increased by 40%', false, NOW() - INTERVAL '4 hours'),
    (v_tenant_id, 'data_quality', 'low', 'Data Quality Warning',
     '5 interactions missing customer identifiers', true, NOW() - INTERVAL '1 day')
    ON CONFLICT DO NOTHING;

    -- ============================================================================
    -- 19. FRAUD DETECTION (Commented out - schema differences, can be added later)
    -- ============================================================================
    -- INSERT INTO fraud_detection_rules (tenant_id, rule_name, rule_type, rule_config, is_active, created_at) VALUES
    -- (v_tenant_id, 'Duplicate Conversion Detection', 'duplicate', 
    --  '{"max_duplicates": 1, "time_window": 3600}'::jsonb, true, NOW() - INTERVAL '6 months'),
    -- (v_tenant_id, 'Suspicious Revenue Pattern', 'revenue',
    --  '{"max_deviation": 3, "threshold": 100000}'::jsonb, true, NOW() - INTERVAL '6 months')
    -- ON CONFLICT DO NOTHING;

    -- INSERT INTO fraud_incidents (tenant_id, rule_id, incident_type, severity, description,
    --     customer_id, status, created_at) VALUES
    -- (v_tenant_id, (SELECT id FROM fraud_detection_rules WHERE tenant_id = v_tenant_id LIMIT 1),
    --  'duplicate_conversion', 'medium', 'Duplicate conversion event detected', NULL, 'investigating', NOW() - INTERVAL '3 days')
    -- ON CONFLICT DO NOTHING;

    -- ============================================================================
    -- 20. DATA QUALITY SCORES
    -- ============================================================================
    INSERT INTO data_quality_scores (tenant_id, entity_type, entity_id, overall_score, issues, calculated_at) VALUES
    (v_tenant_id, 'customer', v_cust_sharma_1,
     0.95, '[]'::jsonb, NOW() - INTERVAL '1 day'),
    (v_tenant_id, 'interaction', (SELECT id FROM interactions WHERE tenant_id = v_tenant_id LIMIT 1),
     0.88, '["missing_metadata"]'::jsonb, NOW() - INTERVAL '1 day')
    ON CONFLICT DO NOTHING;

    -- ============================================================================
    -- 21. PAGE VIEWS & USER ACTIONS (Commented out - requires session_id, can be added later)
    -- ============================================================================
    -- INSERT INTO page_views (tenant_id, customer_id, session_id, page_url, page_title, referrer,
    --     view_timestamp, created_at) VALUES
    -- (v_tenant_id, v_cust_sharma_1, 'session-001',
    --  '/pricing', 'Khatabook Pricing', 'google.com', NOW() - INTERVAL '60 days', NOW() - INTERVAL '60 days'),
    -- (v_tenant_id, v_cust_patel_1, 'session-002',
    --  '/features', 'Khatabook Features', 'facebook.com', NOW() - INTERVAL '45 days', NOW() - INTERVAL '45 days'),
    -- (v_tenant_id, v_cust_singh_1, 'session-003',
    --  '/signup', 'Sign Up - Khatabook', 'google.com', NOW() - INTERVAL '40 days', NOW() - INTERVAL '40 days')
    -- ON CONFLICT DO NOTHING;

    -- INSERT INTO user_actions (tenant_id, customer_id, action_type, action_data, action_timestamp, created_at) VALUES
    -- (v_tenant_id, v_cust_sharma_1,
    --  'button_click', '{"button": "start_free_trial"}'::jsonb, NOW() - INTERVAL '58 days', NOW() - INTERVAL '58 days'),
    -- (v_tenant_id, v_cust_patel_1,
    --  'form_submit', '{"form": "contact_sales"}'::jsonb, NOW() - INTERVAL '43 days', NOW() - INTERVAL '43 days')
    -- ON CONFLICT DO NOTHING;

    -- ============================================================================
    -- 22. SESSIONS (Commented out - can be added later)
    -- ============================================================================
    -- INSERT INTO sessions (tenant_id, customer_id, session_start, session_end, page_views_count,
    --     actions_count, created_at) VALUES
    -- (v_tenant_id, v_cust_sharma_1,
    --  NOW() - INTERVAL '60 days', NOW() - INTERVAL '60 days' + INTERVAL '15 minutes', 5, 3, NOW() - INTERVAL '60 days'),
    -- (v_tenant_id, v_cust_patel_1,
    --  NOW() - INTERVAL '45 days', NOW() - INTERVAL '45 days' + INTERVAL '8 minutes', 3, 2, NOW() - INTERVAL '45 days')
    -- ON CONFLICT DO NOTHING;

    -- ============================================================================
    -- 23. INTEGRATIONS
    -- ============================================================================
    INSERT INTO integrations (tenant_id, platform, integration_type, is_active, credentials, sync_config, created_at) VALUES
    (v_tenant_id, 'Salesforce', 'crm', true,
     '{"api_key": "***"}'::jsonb, '{"sync_frequency": "hourly"}'::jsonb, NOW() - INTERVAL '1 year'),
    (v_tenant_id, 'Google Ads', 'ad_platform', true,
     '{"account_id": "***"}'::jsonb, '{"sync_frequency": "daily"}'::jsonb, NOW() - INTERVAL '1 year'),
    (v_tenant_id, 'Facebook Ads', 'ad_platform', true,
     '{"account_id": "***"}'::jsonb, '{"sync_frequency": "daily"}'::jsonb, NOW() - INTERVAL '1 year'),
    (v_tenant_id, 'HubSpot', 'marketing_automation', true,
     '{"api_key": "***"}'::jsonb, '{"sync_frequency": "real-time"}'::jsonb, NOW() - INTERVAL '6 months')
    ON CONFLICT DO NOTHING;

    -- ============================================================================
    -- 24-28. ADVANCED FEATURES (Commented out - schema differences, can be added later)
    -- ============================================================================
    -- Saved Reports, Experiments, Feature Flags, MMM Models, Attribution Models/Runs/Results
    -- These can be uncommented and fixed once schema is verified

END $$;

-- ============================================================================
-- VERIFICATION QUERIES (Optional - Run to verify data)
-- ============================================================================
-- SELECT COUNT(*) as total_customers FROM customers WHERE tenant_id = (SELECT id FROM tenants WHERE code = 'khatabook');
-- SELECT COUNT(*) as total_interactions FROM interactions WHERE tenant_id = (SELECT id FROM tenants WHERE code = 'khatabook');
-- SELECT COUNT(*) as total_conversions FROM conversion_events WHERE tenant_id = (SELECT id FROM tenants WHERE code = 'khatabook');
-- SELECT SUM(revenue_amount) as total_revenue FROM conversion_events WHERE tenant_id = (SELECT id FROM tenants WHERE code = 'khatabook');

