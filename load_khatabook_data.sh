#!/bin/bash

# ============================================================================
# Load Khatabook B2B Company Dummy Data
# ============================================================================

set -e

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║         Loading Khatabook B2B Company Dummy Data              ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

DB_NAME="convin_crae"
DB_USER="${PGUSER:-$(whoami)}"
DB_HOST="${PGHOST:-localhost}"
DB_PORT="${PGPORT:-5432}"

echo "📊 Database: $DB_NAME"
echo "👤 User: $DB_USER"
echo "🌐 Host: $DB_HOST:$DB_PORT"
echo ""

# Check if database exists
if ! psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -lqt | cut -d \| -f 1 | grep -qw "$DB_NAME"; then
    echo "❌ Error: Database '$DB_NAME' does not exist!"
    echo "   Please create it first: createdb $DB_NAME"
    exit 1
fi

echo "✅ Database found"
echo ""

# Load the seed data
echo "📥 Loading Khatabook seed data..."
if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f database/seed_khatabook_data.sql; then
    echo ""
    echo "✅ Khatabook data loaded successfully!"
    echo ""
    
    # Show summary
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║                    DATA SUMMARY                                 ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo ""
    
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF
SELECT 
    'Tenants' as module, COUNT(*) as count FROM tenants WHERE code = 'khatabook'
UNION ALL
SELECT 'Accounts', COUNT(*) FROM accounts WHERE tenant_id = (SELECT id FROM tenants WHERE code = 'khatabook')
UNION ALL
SELECT 'Customers', COUNT(*) FROM customers WHERE tenant_id = (SELECT id FROM tenants WHERE code = 'khatabook')
UNION ALL
SELECT 'Interactions', COUNT(*) FROM interactions WHERE tenant_id = (SELECT id FROM tenants WHERE code = 'khatabook')
UNION ALL
SELECT 'Conversions', COUNT(*) FROM conversion_events WHERE tenant_id = (SELECT id FROM tenants WHERE code = 'khatabook')
UNION ALL
SELECT 'Agents', COUNT(*) FROM agents WHERE tenant_id = (SELECT id FROM tenants WHERE code = 'khatabook')
UNION ALL
SELECT 'Campaigns', COUNT(*) FROM campaigns WHERE tenant_id = (SELECT id FROM tenants WHERE code = 'khatabook')
UNION ALL
SELECT 'Lead Scores', COUNT(*) FROM lead_scores WHERE tenant_id = (SELECT id FROM tenants WHERE code = 'khatabook')
UNION ALL
SELECT 'ABM Accounts', COUNT(*) FROM accounts WHERE tenant_id = (SELECT id FROM tenants WHERE code = 'khatabook') AND target_account = true
UNION ALL
SELECT 'Integrations', COUNT(*) FROM integrations WHERE tenant_id = (SELECT id FROM tenants WHERE code = 'khatabook')
UNION ALL
SELECT 'MMM Models', COUNT(*) FROM mmm_models WHERE tenant_id = (SELECT id FROM tenants WHERE code = 'khatabook')
UNION ALL
SELECT 'Experiments', COUNT(*) FROM experiments WHERE tenant_id = (SELECT id FROM tenants WHERE code = 'khatabook')
ORDER BY module;
EOF

    echo ""
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║              REVENUE SUMMARY                                   ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo ""
    
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF
SELECT 
    COUNT(*) as total_conversions,
    SUM(revenue_amount) as total_revenue_inr,
    AVG(revenue_amount) as avg_revenue_inr,
    MIN(event_timestamp) as first_conversion,
    MAX(event_timestamp) as latest_conversion
FROM conversion_events
WHERE tenant_id = (SELECT id FROM tenants WHERE code = 'khatabook');
EOF

    echo ""
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║              TOP MERCHANTS BY REVENUE                         ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo ""
    
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF
SELECT 
    a.name as merchant_name,
    a.industry,
    COUNT(DISTINCT ce.id) as conversions,
    SUM(ce.revenue_amount) as total_revenue_inr,
    a.lifecycle_stage
FROM accounts a
JOIN customers c ON c.account_id = a.id
JOIN conversion_events ce ON ce.customer_id = c.id
WHERE a.tenant_id = (SELECT id FROM tenants WHERE code = 'khatabook')
GROUP BY a.id, a.name, a.industry, a.lifecycle_stage
ORDER BY total_revenue_inr DESC
LIMIT 5;
EOF

    echo ""
    echo "✅ All Khatabook data loaded and verified!"
    echo ""
    echo "🚀 Next steps:"
    echo "   1. Start backend: cd backend && go run cmd/server/main.go"
    echo "   2. Start frontend: cd frontend && npm start"
    echo "   3. Access dashboard: http://localhost:3000"
    echo "   4. Use X-Tenant-ID header: 1 (or the tenant ID from database)"
    echo ""
    
else
    echo ""
    echo "❌ Error loading Khatabook data!"
    exit 1
fi

