# Telephony Integration & Attribution Guide

## Overview

This guide explains how clients integrate their telephony systems (Convin, Twilio, or custom) with the Convin Revenue Attribution Engine (CRAE) and how attribution works with call data.

---

## Integration Methods

### Method 1: Webhook Integration (Recommended)

**Best for**: Real-time call tracking, automatic ingestion

Clients configure their telephony provider to send webhooks to CRAE when call events occur.

#### Convin Integration

**Webhook URL**: `https://yourdomain.com/v1/webhooks/convin`

**Authentication**: HMAC-SHA256 signature verification
- Header: `X-Webhook-Signature`
- Secret: Configured in `CONVIN_WEBHOOK_SECRET` environment variable

**Setup Steps**:
1. Configure webhook URL in Convin dashboard
2. Set webhook secret (same as `CONVIN_WEBHOOK_SECRET` in your `.env`)
3. Enable events: `call.started`, `call.ended`, `call.transcript.updated`, `call.intent.detected`

#### Generic Telephony Integration

**Webhook URL**: `https://yourdomain.com/v1/webhooks/telephony`

**Authentication**: HMAC-SHA256 signature verification (optional)
- Header: `X-Webhook-Signature`
- Secret: Configured in `TELEPHONY_WEBHOOK_SECRET` environment variable

### Method 2: Direct API Integration

**Best for**: Custom integrations, batch processing, manual ingestion

Clients can directly call the CRAE API to ingest call data.

**Endpoint**: `POST /v1/interactions`

**Authentication**: API Key or JWT Token
- Header: `X-API-Key` or `Authorization: Bearer <token>`
- Header: `X-Tenant-ID: <tenant_id>`

---

## Required Data Fields

### Minimum Required Fields

For a call to be ingested and attributed, you need:

```json
{
  "external_interaction_id": "call-12345",  // Unique call ID from your system
  "channel": "call",                        // Channel type
  "started_at": "2024-01-15T10:30:00Z",     // ISO 8601 timestamp
  "customer_identifiers": [                 // At least one identifier
    {
      "type": "phone",                      // "phone" or "email"
      "value": "+1234567890"                // Phone number or email
    }
  ]
}
```

### Recommended Fields (For Better Attribution)

```json
{
  "external_interaction_id": "call-12345",
  "channel": "call",
  "vendor_code": "convin",                 // Vendor/BPO identifier
  "started_at": "2024-01-15T10:30:00Z",
  "ended_at": "2024-01-15T10:45:00Z",      // Call end time
  "direction": "inbound",                   // "inbound" or "outbound"
  "language": "en",                         // Language code
  "customer_identifiers": [
    {
      "type": "phone",
      "value": "+1234567890"
    },
    {
      "type": "email",
      "value": "customer@example.com"
    }
  ],
  "participants": [                         // Agent information
    {
      "participant_type": "agent",
      "external_agent_id": "agent-001",
      "role": "sales_rep"
    }
  ],
  "transcript_url": "https://...",         // Transcript location
  "primary_intent": "purchase_inquiry",     // Detected intent
  "secondary_intents": ["product_info", "pricing"],
  "outcome_prediction": "likely_purchase",  // Predicted outcome
  "purchase_probability": 0.75,             // 0.0 to 1.0
  "raw_metadata": {                         // Additional data
    "call_duration": 900,
    "ivr_path": "sales->pricing",
    "sentiment": "positive"
  }
}
```

---

## Webhook Event Types

### 1. Call Started (`call.started`)

**When**: Call is initiated

**Payload Example**:
```json
{
  "event_type": "call.started",
  "call_id": "call-12345",
  "tenant_id": 3,
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {
    "phone_number": "+1234567890",
    "email": "customer@example.com",
    "direction": "inbound",
    "language": "en",
    "vendor_code": "convin",
    "started_at": "2024-01-15T10:30:00Z",
    "agent_id": "agent-001",
    "team_id": "team-sales"
  }
}
```

**What Happens**:
- Creates a new interaction record
- Identifies or creates customer based on phone/email
- Links to vendor/team if provided
- Stores call metadata

### 2. Call Ended (`call.ended`)

**When**: Call is completed

**Payload Example**:
```json
{
  "event_type": "call.ended",
  "call_id": "call-12345",
  "tenant_id": 3,
  "timestamp": "2024-01-15T10:45:00Z",
  "data": {
    "ended_at": "2024-01-15T10:45:00Z",
    "transcript_url": "https://storage.example.com/transcripts/call-12345.json",
    "outcome": "purchase",
    "duration_seconds": 900
  }
}
```

**What Happens**:
- Updates interaction with end time and duration
- Stores transcript location
- Updates outcome if available

### 3. Transcript Updated (`call.transcript.updated`)

**When**: Transcript is available or updated

**Payload Example**:
```json
{
  "event_type": "call.transcript.updated",
  "call_id": "call-12345",
  "tenant_id": 3,
  "timestamp": "2024-01-15T10:46:00Z",
  "data": {
    "transcript_url": "https://storage.example.com/transcripts/call-12345.json"
  }
}
```

**What Happens**:
- Updates interaction with transcript location
- Transcript can be used for intent detection and analysis

### 4. Intent Detected (`call.intent.detected`)

**When**: AI detects intent from call

**Payload Example**:
```json
{
  "event_type": "call.intent.detected",
  "call_id": "call-12345",
  "tenant_id": 3,
  "timestamp": "2024-01-15T10:47:00Z",
  "data": {
    "primary_intent": "purchase_inquiry",
    "secondary_intents": ["product_info", "pricing"],
    "purchase_probability": 0.75
  }
}
```

**What Happens**:
- Updates interaction with detected intents
- Stores purchase probability
- Used for attribution weighting

---

## How Attribution Works

### Step 1: Call Ingestion

When a call webhook is received:

1. **Customer Identification**
   - System looks up customer by phone number or email
   - If not found, creates new customer record
   - Links call to customer

2. **Interaction Creation**
   - Creates interaction record with:
     - Call metadata (duration, direction, language)
     - Customer ID
     - Vendor/Team/Agent information
     - Intent and outcome data
     - Timestamp

3. **Data Enrichment**
   - Links to vendor/BPO if vendor_code provided
   - Links to team/agent if participant data provided
   - Stores transcript URL for later analysis

### Step 2: Conversion Tracking

When a conversion occurs (purchase, signup, etc.):

**API Call**:
```bash
POST /v1/conversions
{
  "external_event_id": "order-789",
  "event_type": "purchase",
  "customer_identifiers": [
    {"type": "phone", "value": "+1234567890"}
  ],
  "amount_decimal": 5000.00,
  "currency": "USD",
  "occurred_at": "2024-01-20T14:00:00Z",
  "event_source": "billing"
}
```

**What Happens**:
- Creates conversion event
- Links to customer
- Stores conversion amount and metadata

### Step 3: Attribution Run

Attribution connects calls to conversions:

**API Call**:
```bash
POST /v1/attribution/runs
{
  "model_code": "TIME_DECAY",
  "name": "Q1 2024 Attribution",
  "config": {
    "time_window_hours": 72,
    "include_channels": ["call"],
    "event_types": ["purchase"],
    "min_purchase_amount": 100.00
  }
}

# Then execute
POST /v1/attribution/runs/{run_id}/execute
```

**Attribution Process**:

1. **Time Window Lookup**
   - For each conversion, finds all interactions within time window (default: 72 hours before conversion)
   - Only includes interactions from specified channels

2. **Attribution Model Application**

   **Available Models**:
   
   - **FIRST_TOUCH**: 100% credit to first interaction
   - **LAST_TOUCH**: 100% credit to last interaction
   - **LINEAR**: Equal credit to all interactions
   - **TIME_DECAY**: More credit to recent interactions (exponential decay)
   - **AI_WEIGHTED**: Credit based on purchase probability and intent

3. **Weight Calculation**

   Example with TIME_DECAY model:
   ```
   Conversion: $1000 purchase
   Interactions:
     - Call 1 (3 days before): 10% weight = $100
     - Call 2 (1 day before): 30% weight = $300
     - Call 3 (same day): 60% weight = $600
   ```

4. **Attribution Results**

   For each interaction, creates attribution result:
   ```json
   {
     "conversion_event_id": 123,
     "interaction_id": 456,
     "agent_id": 789,
     "team_id": 10,
     "vendor_id": 5,
     "attribution_weight": 0.30,
     "attributed_amount": 300.00,
     "is_primary_touch": false
   }
   ```

### Step 4: Analytics & Reporting

Attribution results are used for:

1. **Agent Performance**
   - Revenue attributed to each agent
   - Conversion rates by agent
   - Average deal size per agent

2. **Vendor/BPO Comparison**
   - Revenue by vendor
   - Cost per acquisition by vendor
   - ROI comparison

3. **Intent Analysis**
   - Revenue by intent type
   - Intent-to-conversion rates
   - Most valuable intents

4. **Channel Performance**
   - Revenue by channel (call, email, chat, etc.)
   - Multi-channel attribution
   - Channel ROI

---

## Complete Integration Example

### Scenario: E-commerce Company with Convin

**Setup**:
1. Company uses Convin for call tracking
2. Has billing system for conversions
3. Wants to attribute revenue to calls

### Step 1: Configure Webhook

In Convin dashboard:
- Webhook URL: `https://api.company.com/v1/webhooks/convin`
- Webhook Secret: `your-secret-key`
- Events: Enable all call events

### Step 2: Calls Flow Automatically

When customer calls:
```
1. Customer calls → Convin receives call
2. Convin sends webhook → CRAE receives "call.started"
3. CRAE creates interaction → Links to customer
4. Call ends → Convin sends "call.ended"
5. CRAE updates interaction → Stores duration, transcript
6. AI analyzes → Convin sends "call.intent.detected"
7. CRAE updates interaction → Stores intent, probability
```

### Step 3: Track Conversions

When customer purchases:
```bash
POST /v1/conversions
{
  "external_event_id": "order-12345",
  "event_type": "purchase",
  "customer_identifiers": [
    {"type": "phone", "value": "+1234567890"}
  ],
  "amount_decimal": 2500.00,
  "currency": "USD",
  "occurred_at": "2024-01-20T15:00:00Z",
  "event_source": "billing"
}
```

### Step 4: Run Attribution

```bash
# Create attribution run
POST /v1/attribution/runs
{
  "model_code": "TIME_DECAY",
  "name": "January 2024",
  "config": {
    "time_window_hours": 72,
    "include_channels": ["call"],
    "event_types": ["purchase"]
  }
}

# Execute
POST /v1/attribution/runs/{run_id}/execute
```

### Step 5: View Results

```bash
# Get agent revenue
GET /v1/analytics/agents/revenue?start_date=2024-01-01&end_date=2024-01-31

# Get vendor comparison
GET /v1/analytics/vendors/comparison?start_date=2024-01-01&end_date=2024-01-31

# Get intent profitability
GET /v1/analytics/intents/revenue?start_date=2024-01-01&end_date=2024-01-31
```

---

## Attribution Models Explained

### FIRST_TOUCH
**Best for**: Understanding initial customer acquisition

- 100% credit to first touchpoint
- Ignores all subsequent interactions
- Simple but may undervalue nurturing

**Example**:
```
Call 1 (Day 1): $1000 → 100% = $1000
Call 2 (Day 2): $1000 → 0% = $0
Call 3 (Day 3): $1000 → 0% = $0
Purchase (Day 3): $1000
```

### LAST_TOUCH
**Best for**: Understanding final conversion drivers

- 100% credit to last touchpoint
- Ignores all previous interactions
- Simple but may undervalue awareness

**Example**:
```
Call 1 (Day 1): $1000 → 0% = $0
Call 2 (Day 2): $1000 → 0% = $0
Call 3 (Day 3): $1000 → 100% = $1000
Purchase (Day 3): $1000
```

### LINEAR
**Best for**: Fair distribution across all touchpoints

- Equal credit to all touchpoints
- Recognizes all interactions
- May overvalue low-impact touches

**Example**:
```
Call 1 (Day 1): $1000 → 33.3% = $333.33
Call 2 (Day 2): $1000 → 33.3% = $333.33
Call 3 (Day 3): $1000 → 33.3% = $333.33
Purchase (Day 3): $1000
```

### TIME_DECAY
**Best for**: Balancing all touchpoints with recency bias

- More credit to recent interactions
- Exponential decay formula
- Recognizes both awareness and conversion

**Example**:
```
Call 1 (3 days before): $1000 → 10% = $100
Call 2 (1 day before): $1000 → 30% = $300
Call 3 (same day): $1000 → 60% = $600
Purchase: $1000
```

### AI_WEIGHTED
**Best for**: Data-driven attribution using ML

- Credit based on purchase probability
- Considers intent, sentiment, duration
- Most accurate but requires ML model

**Example**:
```
Call 1 (low probability): $1000 → 15% = $150
Call 2 (medium probability): $1000 → 25% = $250
Call 3 (high probability): $1000 → 60% = $600
Purchase: $1000
```

---

## Best Practices

### 1. Customer Identification
- **Always provide phone number** (most reliable for calls)
- Include email if available (for cross-channel matching)
- Use consistent format (E.164 for phone numbers)

### 2. Timestamp Accuracy
- Use ISO 8601 format with timezone
- Ensure server clocks are synchronized
- Include both started_at and ended_at for accurate duration

### 3. Vendor/Team Tracking
- Provide vendor_code for BPO comparison
- Include agent/team information for granular analysis
- Use consistent identifiers across calls

### 4. Intent Data
- Send intent detection events when available
- Include purchase_probability for AI_WEIGHTED model
- Update intents as call progresses

### 5. Conversion Tracking
- Track conversions immediately after they occur
- Use same customer identifiers as calls
- Include accurate timestamps

### 6. Attribution Configuration
- Choose model based on business goals
- Set appropriate time window (typically 7-30 days)
- Filter by channel/event type as needed

---

## Troubleshooting

### Calls Not Appearing
1. Check webhook URL is correct
2. Verify webhook signature secret matches
3. Check logs: `docker-compose logs backend`
4. Verify customer identifiers are provided

### Attribution Not Working
1. Ensure conversions are tracked
2. Check time window includes calls before conversion
3. Verify customer identifiers match between calls and conversions
4. Check attribution run status: `GET /v1/attribution/runs/{run_id}`

### Missing Revenue
1. Verify all calls are being ingested
2. Check attribution model configuration
3. Ensure time window is sufficient
4. Verify customer matching is working

---

## API Reference

### Webhook Endpoints
- `POST /v1/webhooks/convin` - Convin webhook
- `POST /v1/webhooks/telephony` - Generic telephony webhook

### Ingestion Endpoints
- `POST /v1/interactions` - Manual interaction ingestion
- `POST /v1/conversions` - Conversion tracking

### Attribution Endpoints
- `POST /v1/attribution/runs` - Create attribution run
- `POST /v1/attribution/runs/{run_id}/execute` - Execute attribution
- `GET /v1/attribution/runs/{run_id}` - Get run status

### Analytics Endpoints
- `GET /v1/analytics/agents/revenue` - Agent revenue
- `GET /v1/analytics/vendors/comparison` - Vendor comparison
- `GET /v1/analytics/intents/revenue` - Intent revenue

---

## Support

For integration help:
1. Check logs: `docker-compose logs backend`
2. Test webhook: Use curl to send test payload
3. Verify health: `GET /health`
4. Contact support with webhook payload examples

