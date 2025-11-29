#!/bin/bash

# Test Notification Script
# This script creates test notifications to verify the notification center

API_URL="${API_URL:-http://localhost:8080/v1}"
TENANT_ID="${TENANT_ID:-1}"

# Get auth token (you may need to login first)
TOKEN="${TOKEN:-}"
if [ -z "$TOKEN" ]; then
    echo "⚠️  No token provided. Please login first and set TOKEN environment variable."
    echo "   Example: TOKEN=your_token_here ./test_notifications.sh"
    echo ""
    echo "   Or login via:"
    echo "   curl -X POST $API_URL/users/login -H 'Content-Type: application/json' -d '{\"email\":\"your@email.com\",\"password\":\"yourpassword\"}'"
    exit 1
fi

echo "🔔 Creating test notifications..."
echo ""

# Test Notification 1: Error Alert
echo "Creating Error Alert..."
curl -X POST "$API_URL/realtime/alerts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d '{
    "alert_type": "system_error",
    "severity": "error",
    "title": "Database Connection Failed",
    "description": "Unable to connect to primary database. Fallback to replica.",
    "entity_type": "system",
    "entity_id": null,
    "metadata": {
      "error_code": "DB_CONN_001",
      "retry_count": 3
    }
  }' | jq '.'

echo ""
sleep 1

# Test Notification 2: Warning Alert
echo "Creating Warning Alert..."
curl -X POST "$API_URL/realtime/alerts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d '{
    "alert_type": "data_quality",
    "severity": "warning",
    "title": "Low Data Quality Score",
    "description": "Data quality score dropped below 80% for the last hour.",
    "entity_type": "data_quality",
    "entity_id": null,
    "metadata": {
      "score": 75,
      "threshold": 80
    }
  }' | jq '.'

echo ""
sleep 1

# Test Notification 3: Info Alert
echo "Creating Info Alert..."
curl -X POST "$API_URL/realtime/alerts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d '{
    "alert_type": "attribution_complete",
    "severity": "info",
    "title": "Attribution Run Completed",
    "description": "Q1 2024 attribution analysis has been completed successfully.",
    "entity_type": "attribution_run",
    "entity_id": 123,
    "metadata": {
      "run_id": 123,
      "conversions_processed": 1250,
      "duration_seconds": 45
    }
  }' | jq '.'

echo ""
sleep 1

# Test Notification 4: Success Alert
echo "Creating Success Alert..."
curl -X POST "$API_URL/realtime/alerts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d '{
    "alert_type": "integration_success",
    "severity": "success",
    "title": "Salesforce Sync Successful",
    "description": "Successfully synced 250 contacts from Salesforce.",
    "entity_type": "integration",
    "entity_id": 5,
    "metadata": {
      "integration_id": 5,
      "records_synced": 250,
      "sync_duration": 120
    }
  }' | jq '.'

echo ""
sleep 1

# Test Notification 5: Critical Alert
echo "Creating Critical Alert..."
curl -X POST "$API_URL/realtime/alerts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Tenant-ID: $TENANT_ID" \
  -d '{
    "alert_type": "fraud_detected",
    "severity": "error",
    "title": "Potential Fraud Detected",
    "description": "Unusual pattern detected in conversion events. Requires immediate review.",
    "entity_type": "fraud_incident",
    "entity_id": 789,
    "metadata": {
      "incident_id": 789,
      "risk_score": 0.95,
      "pattern_type": "velocity_anomaly"
    }
  }' | jq '.'

echo ""
echo "✅ Test notifications created!"
echo ""
echo "📱 Check your notification center in the UI to see the new alerts."
echo "   The notification icon should show a badge with the unread count."

