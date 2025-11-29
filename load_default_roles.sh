#!/bin/bash

# Load Default Roles for Convin Economics
# This script creates default roles for all tenants

echo "=========================================="
echo "Loading Default Roles for Convin Economics"
echo "=========================================="

# Check if database exists
if ! psql -lqt | cut -d \| -f 1 | grep -qw convin_crae; then
    echo "❌ Database 'convin_crae' does not exist!"
    echo "Please create it first: createdb convin_crae"
    exit 1
fi

echo ""
echo "📋 Available tenants:"
psql convin_crae -c "SELECT id, name, code FROM tenants ORDER BY id;"

echo ""
read -p "Enter tenant ID to create roles for (or 'all' for all tenants): " tenant_input

if [ "$tenant_input" = "all" ]; then
    echo ""
    echo "🔄 Creating default roles for all tenants..."
    
    # Get all tenant IDs
    tenant_ids=$(psql convin_crae -t -c "SELECT id FROM tenants;")
    
    for tenant_id in $tenant_ids; do
        if [ ! -z "$tenant_id" ]; then
            echo ""
            echo "Creating roles for tenant_id: $tenant_id"
            psql convin_crae -c "
            DO \$\$
            DECLARE
                v_tenant_id BIGINT := $tenant_id;
                v_admin_role_id INT;
            BEGIN
                -- Check if roles already exist
                IF EXISTS (SELECT 1 FROM roles WHERE tenant_id = v_tenant_id AND is_default = true) THEN
                    RAISE NOTICE 'Default roles already exist for tenant_id: %', v_tenant_id;
                ELSE
                    -- Create Admin role
                    INSERT INTO roles (tenant_id, name, description, code_names, can_be_edited, is_default)
                    VALUES (
                        v_tenant_id,
                        'Admin',
                        'Full system access with all permissions',
                        ARRAY[
                            'analytics.agents.view', 'analytics.agents.edit',
                            'analytics.vendors.view', 'analytics.vendors.edit',
                            'analytics.intents.view', 'analytics.journey.view',
                            'analytics.mmm.view', 'analytics.mmm.run',
                            'analytics.abm.view', 'analytics.abm.edit',
                            'analytics.lead_scoring.view', 'analytics.lead_scoring.edit',
                            'analytics.cohorts.view', 'analytics.realtime.view',
                            'platform.integrations.view', 'platform.integrations.edit',
                            'platform.experiments.view', 'platform.experiments.edit',
                            'platform.reports.view', 'platform.reports.edit',
                            'users.view', 'users.create', 'users.edit', 'users.delete', 'users.bulk_upload',
                            'roles.view', 'roles.create', 'roles.edit', 'roles.delete',
                            'data.ingest', 'data.attribution.run', 'data.attribution.view', 'data.export',
                            'reports.view', 'reports.create', 'reports.edit', 'reports.delete'
                        ],
                        false, true
                    ) RETURNING id INTO v_admin_role_id;
                    
                    RAISE NOTICE 'Created Admin role (ID: %) for tenant_id: %', v_admin_role_id, v_tenant_id;
                END IF;
            END \$\$;
            "
        fi
    done
    
    echo ""
    echo "✅ Running full default roles script..."
    psql convin_crae -f database/migrations/add_default_roles.sql
    
else
    tenant_id=$tenant_input
    echo ""
    echo "🔄 Creating default roles for tenant_id: $tenant_id"
    
    # Update the SQL file with the tenant ID
    sed "s/v_tenant_id BIGINT := 3;/v_tenant_id BIGINT := $tenant_id;/" database/migrations/add_default_roles.sql > /tmp/add_default_roles_temp.sql
    
    psql convin_crae -f /tmp/add_default_roles_temp.sql
    
    rm /tmp/add_default_roles_temp.sql
fi

echo ""
echo "✅ Default roles created successfully!"
echo ""
echo "📊 Verification - Roles created:"
psql convin_crae -c "
SELECT 
    tenant_id,
    name,
    description,
    array_length(code_names, 1) as permission_count,
    is_default,
    can_be_edited
FROM roles
WHERE is_default = true
ORDER BY tenant_id, name;
"

