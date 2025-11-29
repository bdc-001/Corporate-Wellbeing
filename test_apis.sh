#!/bin/bash

# ============================================================================
# API Testing Script - Comprehensive Feature Set
# Tests all major endpoints with the loaded sample data
# ============================================================================

BASE_URL="http://localhost:8080"
TENANT_ID="1"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Function to test an endpoint
test_endpoint() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    
    echo -e "${BLUE}Testing: ${name}${NC}"
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" -H "X-Tenant-ID: $TENANT_ID" "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" -H "Content-Type: application/json" -H "X-Tenant-ID: $TENANT_ID" -d "$data" "$BASE_URL$endpoint")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        echo -e "${GREEN}✓ PASS${NC} (HTTP $http_code)"
        echo "$body" | jq -C '.' 2>/dev/null | head -20
        ((TESTS_PASSED++))
    else
        echo -e "${RED}✗ FAIL${NC} (HTTP $http_code)"
        echo "$body"
        ((TESTS_FAILED++))
    fi
    echo ""
    sleep 0.5
}

echo "============================================================================"
echo "🧪 CONVIN CRAE - API TEST SUITE"
echo "============================================================================"
echo "Base URL: $BASE_URL"
echo "Tenant ID: $TENANT_ID"
echo ""

# Check if server is running
echo "Checking if server is running..."
if ! curl -s "$BASE_URL/health" > /dev/null; then
    echo -e "${RED}❌ Server is not running!${NC}"
    echo "Please start the server first:"
    echo "  cd backend && go run cmd/server/main.go"
    exit 1
fi
echo -e "${GREEN}✓ Server is running${NC}"
echo ""

# ============================================================================
# 1. HEALTH CHECK
# ============================================================================
echo "============================================================================"
echo "1. HEALTH & SYSTEM INFO"
echo "============================================================================"
test_endpoint "Health Check" "GET" "/health"

# ============================================================================
# 2. ABM (Account-Based Marketing)
# ============================================================================
echo "============================================================================"
echo "2. ACCOUNT-BASED MARKETING (ABM)"
echo "============================================================================"
test_endpoint "List All Accounts" "GET" "/v1/abm/accounts"
test_endpoint "List Target Accounts Only" "GET" "/v1/abm/accounts?target_only=true"
test_endpoint "Get Specific Account" "GET" "/v1/abm/accounts/1"
test_endpoint "Get Account Summary" "GET" "/v1/abm/accounts/1/summary"
test_endpoint "Target Account Insights" "GET" "/v1/abm/insights/target-accounts"

# ============================================================================
# 3. LEAD SCORING & PREDICTIONS
# ============================================================================
echo "============================================================================"
echo "3. LEAD SCORING & PREDICTIVE ANALYTICS"
echo "============================================================================"
test_endpoint "Get High-Value Leads (score > 70)" "GET" "/v1/leads/high-value?min_score=70"
test_endpoint "Get High-Value Leads (score > 80)" "GET" "/v1/leads/high-value?min_score=80"

# ============================================================================
# 4. COHORT ANALYSIS
# ============================================================================
echo "============================================================================"
echo "4. COHORT ANALYSIS"
echo "============================================================================"
test_endpoint "Get Retention Curve for Segment 2" "GET" "/v1/cohorts/segments/2/retention"

# ============================================================================
# 5. REAL-TIME & ALERTS
# ============================================================================
echo "============================================================================"
echo "5. REAL-TIME DATA & ALERTS"
echo "============================================================================"
test_endpoint "Get All Alerts" "GET" "/v1/realtime/alerts"
test_endpoint "Get Unresolved Alerts" "GET" "/v1/realtime/alerts?unresolved=true"
test_endpoint "Get Critical Alerts" "GET" "/v1/realtime/alerts?severity=critical"
test_endpoint "Real-Time Metrics (15 min window)" "GET" "/v1/analytics/realtime/metrics?window=15"

# ============================================================================
# 6. BEHAVIOR ANALYTICS
# ============================================================================
echo "============================================================================"
echo "6. USER BEHAVIOR ANALYTICS"
echo "============================================================================"
test_endpoint "Top Pages (Last 7 Days)" "GET" "/v1/behavior/pages/top?limit=10"
test_endpoint "Session Details" "GET" "/v1/behavior/sessions/sess_abc123"

# ============================================================================
# 7. INTEGRATIONS
# ============================================================================
echo "============================================================================"
echo "7. INTEGRATIONS"
echo "============================================================================"
test_endpoint "List All Integrations" "GET" "/v1/integrations"
test_endpoint "List Active Integrations" "GET" "/v1/integrations?active_only=true"
test_endpoint "List CRM Integrations" "GET" "/v1/integrations?type=crm"

# ============================================================================
# 8. CUSTOM REPORTS
# ============================================================================
echo "============================================================================"
echo "8. CUSTOM REPORTS"
echo "============================================================================"
test_endpoint "List All Reports" "GET" "/v1/reports"
test_endpoint "List Public Reports" "GET" "/v1/reports?public_only=true"

# ============================================================================
# 9. EXPERIMENTS & FEATURE FLAGS
# ============================================================================
echo "============================================================================"
echo "9. EXPERIMENTS & FEATURE FLAGS"
echo "============================================================================"
test_endpoint "List Feature Flags" "GET" "/v1/features/flags"

# ============================================================================
# 10. CORE ANALYTICS (Existing)
# ============================================================================
echo "============================================================================"
echo "10. CORE ANALYTICS"
echo "============================================================================"
test_endpoint "Agent Revenue Summary" "GET" "/v1/analytics/agents/revenue"
test_endpoint "Vendor Comparison" "GET" "/v1/analytics/vendors/comparison"
test_endpoint "Intent Profitability" "GET" "/v1/analytics/intents/revenue"

# ============================================================================
# 11. ADVANCED ANALYTICS (Factors.ai-style)
# ============================================================================
echo "============================================================================"
echo "11. ADVANCED ANALYTICS (FACTORS.AI-STYLE)"
echo "============================================================================"
test_endpoint "Funnel Stage Metrics" "GET" "/v1/analytics/funnel/stages"
test_endpoint "Content Engagement" "GET" "/v1/analytics/content/engagement"
test_endpoint "Multi-Channel ROI" "GET" "/v1/analytics/channels/roi"
test_endpoint "Journey Velocity" "GET" "/v1/analytics/journey/velocity"

# ============================================================================
# 12. CUSTOMER JOURNEY
# ============================================================================
echo "============================================================================"
echo "12. CUSTOMER JOURNEY"
echo "============================================================================"
test_endpoint "Customer Journey (Customer 1)" "GET" "/v1/customers/1/journey"
test_endpoint "Customer Journey (Customer 3)" "GET" "/v1/customers/3/journey"

# ============================================================================
# SUMMARY
# ============================================================================
echo "============================================================================"
echo "📊 TEST SUMMARY"
echo "============================================================================"
TOTAL_TESTS=$((TESTS_PASSED + TESTS_FAILED))
echo "Total Tests: $TOTAL_TESTS"
echo -e "${GREEN}Passed: $TESTS_PASSED${NC}"
echo -e "${RED}Failed: $TESTS_FAILED${NC}"

if [ $TESTS_FAILED -eq 0 ]; then
    echo ""
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    echo "============================================================================"
    exit 0
else
    echo ""
    echo -e "${YELLOW}⚠️  Some tests failed. Check the output above for details.${NC}"
    echo "============================================================================"
    exit 1
fi

