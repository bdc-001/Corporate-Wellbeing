-- ============================================
-- USER MANAGEMENT SYSTEM
-- ============================================

-- Roles table
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    code_names TEXT[] DEFAULT '{}', -- Array of permission code names
    can_be_edited BOOLEAN DEFAULT TRUE,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    UNIQUE (tenant_id, name)
);

CREATE INDEX idx_roles_tenant ON roles(tenant_id);

-- Role-Team mapping (team restrictions)
CREATE TABLE IF NOT EXISTS role_teams (
    role_id INT NOT NULL,
    team_id INT NOT NULL,
    PRIMARY KEY (role_id, team_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE
);

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    role_id INT,
    manager_id BIGINT, -- References users.id
    auditor_id BIGINT, -- References users.id
    team_id INT,
    user_type VARCHAR(50) DEFAULT 'product_user', -- product_user (can access based on role), observer (no functional access)
    location VARCHAR(255),
    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMP,
    password_reset_token VARCHAR(255),
    password_reset_expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE SET NULL,
    FOREIGN KEY (manager_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (auditor_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE SET NULL,
    UNIQUE (tenant_id, email)
);

CREATE INDEX idx_users_tenant ON users(tenant_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role_id);
CREATE INDEX idx_users_manager ON users(manager_id);
CREATE INDEX idx_users_team ON users(team_id);

-- Permission groups (for UI organization)
CREATE TABLE IF NOT EXISTS permission_groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    display_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Permissions (code names with metadata)
CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    code_name VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    group_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES permission_groups(id) ON DELETE SET NULL
);

CREATE INDEX idx_permissions_group ON permissions(group_id);

-- Insert default permission groups
INSERT INTO permission_groups (name, description, display_order) VALUES
    ('Core Analytics', 'Access to core analytics dashboards', 1),
    ('Advanced Analytics', 'Access to advanced analytics features', 2),
    ('Platform Management', 'Access to platform configuration', 3),
    ('User Management', 'Access to user and role management', 4),
    ('Data Management', 'Access to data ingestion and management', 5),
    ('Reporting', 'Access to custom reports and exports', 6)
ON CONFLICT DO NOTHING;

-- Insert default permissions
INSERT INTO permissions (code_name, name, description, group_id) VALUES
    -- Core Analytics
    ('analytics.agents.view', 'View Agents Dashboard', 'View agent performance analytics', 1),
    ('analytics.agents.edit', 'Edit Agents', 'Edit agent information', 1),
    ('analytics.vendors.view', 'View Vendors Dashboard', 'View vendor comparison analytics', 1),
    ('analytics.vendors.edit', 'Edit Vendors', 'Edit vendor information', 1),
    ('analytics.intents.view', 'View Intents Dashboard', 'View intent-based analytics', 1),
    ('analytics.journey.view', 'View Customer Journey', 'View customer journey timeline', 1),
    
    -- Advanced Analytics
    ('analytics.mmm.view', 'View MMM Dashboard', 'View Marketing Mix Modeling', 2),
    ('analytics.mmm.run', 'Run MMM Analysis', 'Execute MMM analysis', 2),
    ('analytics.abm.view', 'View ABM Dashboard', 'View Account-Based Marketing', 2),
    ('analytics.abm.edit', 'Edit ABM Accounts', 'Edit account information', 2),
    ('analytics.lead_scoring.view', 'View Lead Scoring', 'View lead scoring dashboard', 2),
    ('analytics.lead_scoring.edit', 'Edit Lead Scores', 'Edit lead scores and models', 2),
    ('analytics.cohorts.view', 'View Cohort Analysis', 'View cohort analysis dashboard', 2),
    ('analytics.realtime.view', 'View Real-time Analytics', 'View real-time metrics', 2),
    
    -- Platform Management
    ('platform.integrations.view', 'View Integrations', 'View integration list', 3),
    ('platform.integrations.edit', 'Edit Integrations', 'Add/edit integrations', 3),
    ('platform.experiments.view', 'View Experiments', 'View A/B tests and experiments', 3),
    ('platform.experiments.edit', 'Edit Experiments', 'Create/edit experiments', 3),
    ('platform.reports.view', 'View Reports', 'View custom reports', 3),
    ('platform.reports.edit', 'Edit Reports', 'Create/edit custom reports', 3),
    
    -- User Management
    ('users.view', 'View Users', 'View user list', 4),
    ('users.create', 'Create Users', 'Create new users', 4),
    ('users.edit', 'Edit Users', 'Edit user information', 4),
    ('users.delete', 'Delete Users', 'Delete users', 4),
    ('users.bulk_upload', 'Bulk Upload Users', 'Upload users via CSV/Excel', 4),
    ('roles.view', 'View Roles', 'View role list', 4),
    ('roles.create', 'Create Roles', 'Create new roles', 4),
    ('roles.edit', 'Edit Roles', 'Edit role permissions', 4),
    ('roles.delete', 'Delete Roles', 'Delete roles', 4),
    
    -- Data Management
    ('data.ingest', 'Ingest Data', 'Ingest interactions and conversions', 5),
    ('data.attribution.run', 'Run Attribution', 'Execute attribution runs', 5),
    ('data.attribution.view', 'View Attribution', 'View attribution results', 5),
    ('data.export', 'Export Data', 'Export data to CSV/Excel', 5),
    
    -- Reporting
    ('reports.view', 'View Reports', 'View all reports', 6),
    ('reports.create', 'Create Reports', 'Create custom reports', 6),
    ('reports.edit', 'Edit Reports', 'Edit custom reports', 6),
    ('reports.delete', 'Delete Reports', 'Delete reports', 6)
ON CONFLICT (code_name) DO NOTHING;

-- Insert default roles for a tenant (example - will be created per tenant)
-- These will be created via API when tenant is set up

