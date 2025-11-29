-- ============================================
-- DEFAULT ROLES FOR CONVIN ECONOMICS
-- ============================================
-- This script creates default roles: Admin and Agent

-- For tenant_id = 3 (Khatabook example)
-- You can modify the tenant_id as needed

DO $$
DECLARE
    v_tenant_id BIGINT := 3; -- Change this to your tenant ID
    v_admin_role_id INT;
    v_agent_role_id INT;
BEGIN
    -- ============================================
    -- 1. ADMIN ROLE - Full Access
    -- ============================================
    INSERT INTO roles (tenant_id, name, description, code_names, can_be_edited, is_default)
    VALUES (
        v_tenant_id,
        'Admin',
        'Full system access with all permissions. Can manage users, roles, and all platform features.',
        ARRAY[
            -- Core Analytics
            'analytics.agents.view', 'analytics.agents.edit',
            'analytics.vendors.view', 'analytics.vendors.edit',
            'analytics.intents.view', 'analytics.journey.view',
            -- Advanced Analytics
            'analytics.mmm.view', 'analytics.mmm.run',
            'analytics.abm.view', 'analytics.abm.edit',
            'analytics.lead_scoring.view', 'analytics.lead_scoring.edit',
            'analytics.cohorts.view', 'analytics.realtime.view',
            -- Platform Management
            'platform.integrations.view', 'platform.integrations.edit',
            'platform.experiments.view', 'platform.experiments.edit',
            'platform.reports.view', 'platform.reports.edit',
            -- User Management
            'users.view', 'users.create', 'users.edit', 'users.delete', 'users.bulk_upload',
            'roles.view', 'roles.create', 'roles.edit', 'roles.delete',
            -- Data Management
            'data.ingest', 'data.attribution.run', 'data.attribution.view', 'data.export',
            -- Reporting
            'reports.view', 'reports.create', 'reports.edit', 'reports.delete'
        ],
        false, -- Cannot be edited
        true   -- Is default role
    )
    RETURNING id INTO v_admin_role_id;

    -- ============================================
    -- 2. AGENT ROLE - Basic Agent Access
    -- ============================================
    INSERT INTO roles (tenant_id, name, description, code_names, can_be_edited, is_default)
    VALUES (
        v_tenant_id,
        'Agent',
        'Basic access for call center agents. Can view their own performance metrics and customer journey data.',
        ARRAY[
            -- Core Analytics (View own data only)
            'analytics.agents.view', 'analytics.journey.view',
            -- Reporting (View only)
            'reports.view'
        ],
        true,  -- Can be edited
        true   -- Is default role
    )
    RETURNING id INTO v_agent_role_id;

    RAISE NOTICE 'Default roles created successfully for tenant_id: %', v_tenant_id;
    RAISE NOTICE 'Admin Role ID: %', v_admin_role_id;
    RAISE NOTICE 'Agent Role ID: %', v_agent_role_id;

END $$;

-- ============================================
-- Verification Query
-- ============================================
-- Run this to verify roles were created:
-- SELECT id, name, description, array_length(code_names, 1) as permission_count, is_default, can_be_edited
-- FROM roles
-- WHERE tenant_id = 3
-- ORDER BY is_default DESC, name;
