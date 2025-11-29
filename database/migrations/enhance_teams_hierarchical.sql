-- ============================================
-- ENHANCE TEAMS TABLE FOR HIERARCHICAL STRUCTURE
-- ============================================

-- Add columns for hierarchical team structure
ALTER TABLE teams 
ADD COLUMN IF NOT EXISTS group_id INT,
ADD COLUMN IF NOT EXISTS about VARCHAR(512) DEFAULT '',
ADD COLUMN IF NOT EXISTS use_case_id INT,
ADD COLUMN IF NOT EXISTS tenant_id BIGINT;

-- Update tenant_id from vendors (for teams that already exist)
UPDATE teams t
SET tenant_id = v.tenant_id
FROM vendors v
WHERE t.vendor_id = v.id AND t.tenant_id IS NULL;

-- Make tenant_id NOT NULL after populating
ALTER TABLE teams 
ALTER COLUMN tenant_id SET NOT NULL;

-- Add foreign key constraints
ALTER TABLE teams
ADD CONSTRAINT fk_teams_group 
FOREIGN KEY (group_id) REFERENCES teams(id) ON DELETE RESTRICT;

ALTER TABLE teams
ADD CONSTRAINT fk_teams_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;

-- Add unique constraint for (name, group_id) to allow same name in different parent teams
ALTER TABLE teams
ADD CONSTRAINT unique_team_name_group UNIQUE (name, group_id);

-- Create use_cases table
CREATE TABLE IF NOT EXISTS use_cases (
    id SERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    UNIQUE (tenant_id, name)
);

CREATE INDEX idx_use_cases_tenant ON use_cases(tenant_id);

-- Add foreign key for use_case_id
ALTER TABLE teams
ADD CONSTRAINT fk_teams_use_case
FOREIGN KEY (use_case_id) REFERENCES use_cases(id) ON DELETE RESTRICT;

-- Update manager_id to reference users instead of agents
-- First, drop the old constraint if it exists
ALTER TABLE teams
DROP CONSTRAINT IF EXISTS fk_teams_manager;

-- Add new constraint referencing users
ALTER TABLE teams
ADD CONSTRAINT fk_teams_manager
FOREIGN KEY (manager_id) REFERENCES users(id) ON DELETE RESTRICT;

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_teams_group ON teams(group_id);
CREATE INDEX IF NOT EXISTS idx_teams_tenant ON teams(tenant_id);
CREATE INDEX IF NOT EXISTS idx_teams_use_case ON teams(use_case_id);
CREATE INDEX IF NOT EXISTS idx_teams_manager ON teams(manager_id);

-- Add comment
COMMENT ON TABLE teams IS 'Hierarchical team structure with parent teams (group_id) and subteams';
COMMENT ON COLUMN teams.group_id IS 'Parent team ID (NULL for top-level teams)';
COMMENT ON COLUMN teams.use_case_id IS 'Use case association (Sales, Collection, Support, etc.)';

