# Features Documentation

## Table of Contents

1. [Overview](#overview)
2. [Core Features](#core-features)
3. [Analytics & Reporting](#analytics--reporting)
4. [User & Access Management](#user--access-management)
5. [Team Management](#team-management)
6. [Attribution Engine](#attribution-engine)
7. [Advanced Analytics](#advanced-analytics)
8. [Integrations](#integrations)
9. [Real-Time Features](#real-time-features)
10. [API Reference](#api-reference)

---

## Overview

**Economics** is a comprehensive Revenue Attribution Engine (CRAE) designed to track, analyze, and attribute revenue across multiple touchpoints. The platform provides end-to-end visibility into customer journeys, from initial contact through conversion, enabling data-driven decision making.

### Key Capabilities

- **Multi-Touch Attribution**: Track and attribute revenue across all customer touchpoints
- **Live Call Flow Integration**: Real-time ingestion of call data from telephony providers
- **Advanced Analytics**: Deep insights into agent, vendor, and intent performance
- **Account-Based Marketing**: Target and track enterprise accounts
- **Predictive Analytics**: Lead scoring and purchase probability predictions
- **Cohort Analysis**: Understand customer behavior over time
- **Marketing Mix Modeling**: Optimize marketing spend allocation
- **Real-Time Monitoring**: Live dashboards and alerts

---

## Core Features

### 1. Dashboard Overview

**Location**: `/` (Home Dashboard)

**Features**:
- Welcome message with personalized user greeting
- Key performance metrics at a glance:
  - Total Revenue
  - Conversions
  - Attribution Runs
  - Active Campaigns
- Quick Actions:
  - Run Attribution Analysis
  - Upload Data
  - View Analytics
  - Data Library
- Recent Activity Feed
- Analytics & Insights Cards

**Use Cases**:
- Daily performance monitoring
- Quick access to common tasks
- Activity tracking

---

### 2. Data Ingestion

#### Interaction Ingestion

**Endpoint**: `POST /v1/interactions`

**Purpose**: Ingest customer interactions (calls, emails, chats, etc.)

**Required Fields**:
- `external_interaction_id`: Unique identifier from source system
- `channel`: Interaction channel (call, email, chat, etc.)
- `started_at`: Timestamp when interaction began
- `customer_identifiers`: At least one identifier (phone, email)

**Optional Fields**:
- `vendor_code`: BPO/Vendor identifier
- `ended_at`: Interaction end time
- `direction`: inbound/outbound
- `language`: Language code
- `primary_intent`: Detected intent
- `purchase_probability`: ML prediction (0.0-1.0)
- `transcript_url`: Transcript location
- `participants`: Agent/team information

**Example**:
```json
{
  "external_interaction_id": "call-12345",
  "channel": "call",
  "vendor_code": "convin",
  "started_at": "2024-01-15T10:30:00Z",
  "ended_at": "2024-01-15T10:45:00Z",
  "direction": "inbound",
  "customer_identifiers": [
    {"type": "phone", "value": "+1234567890"}
  ],
  "primary_intent": "purchase_inquiry",
  "purchase_probability": 0.75
}
```

#### Conversion Tracking

**Endpoint**: `POST /v1/conversions`

**Purpose**: Track customer conversions (purchases, signups, etc.)

**Required Fields**:
- `external_event_id`: Unique conversion identifier
- `event_type`: Type of conversion (purchase, signup, etc.)
- `customer_identifiers`: Customer identification
- `amount_decimal`: Conversion amount
- `currency`: Currency code
- `occurred_at`: Conversion timestamp
- `event_source`: Source system

**Example**:
```json
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

#### Event Tracking

**Endpoint**: `POST /v1/events`

**Purpose**: Track real-time events (page views, form submissions, etc.)

**Endpoint**: `POST /v1/page-views`

**Purpose**: Track page view events

---

### 3. Customer Identity & Journey

#### Customer Journey View

**Location**: `/journey`

**Features**:
- Complete customer interaction timeline
- Multi-channel touchpoint visualization
- Conversion events highlighted
- Intent progression tracking
- Revenue attribution breakdown

**Endpoint**: `GET /v1/customers/:customer_id/journey`

**Returns**:
- Chronological list of all interactions
- Conversion events
- Attribution results
- Customer metadata

---

## Analytics & Reporting

### 1. Agents Dashboard

**Location**: `/agents`

**Features**:
- Agent performance metrics
- Revenue attributed per agent
- Conversion rates
- Average deal size
- Call volume and duration
- Top performing agents

**Endpoint**: `GET /v1/analytics/agents/revenue`

**Parameters**:
- `start_date`: Start date (YYYY-MM-DD)
- `end_date`: End date (YYYY-MM-DD)
- `agent_id`: Filter by specific agent (optional)

**Metrics Provided**:
- Total attributed revenue
- Number of conversions
- Conversion rate
- Average deal size
- Total interactions
- Average call duration

---

### 2. Vendors Dashboard

**Location**: `/vendors`

**Features**:
- Vendor/BPO performance comparison
- Revenue by vendor
- Cost per acquisition (CPA)
- ROI comparison
- Volume metrics
- Quality metrics

**Endpoint**: `GET /v1/analytics/vendors/comparison`

**Parameters**:
- `start_date`: Start date
- `end_date`: End date
- `vendor_id`: Filter by vendor (optional)

**Metrics Provided**:
- Total revenue per vendor
- Number of conversions
- Average deal size
- Total interactions
- Conversion rate
- Cost efficiency

---

### 3. Intents Dashboard

**Location**: `/intents`

**Features**:
- Revenue by intent type
- Intent-to-conversion rates
- Most valuable intents
- Intent distribution
- Conversion funnel by intent

**Endpoint**: `GET /v1/analytics/intents/revenue`

**Parameters**:
- `start_date`: Start date
- `end_date`: End date
- `intent`: Filter by intent type (optional)

**Metrics Provided**:
- Revenue per intent
- Conversion rate by intent
- Average deal size by intent
- Intent frequency
- Purchase probability by intent

---

### 4. Advanced Analytics

#### Funnel Stage Metrics

**Endpoint**: `GET /v1/analytics/funnel/stages`

**Purpose**: Analyze customer progression through funnel stages

**Metrics**:
- Stage conversion rates
- Drop-off points
- Time between stages
- Volume at each stage

#### Content Engagement

**Endpoint**: `GET /v1/analytics/content/engagement`

**Purpose**: Track content performance and engagement

**Metrics**:
- Content views
- Engagement rates
- Conversion attribution
- Top performing content

#### Multi-Channel ROI

**Endpoint**: `GET /v1/analytics/channels/roi`

**Purpose**: Calculate ROI across all marketing channels

**Metrics**:
- Revenue by channel
- Cost by channel
- ROI per channel
- Channel efficiency

#### Journey Velocity

**Endpoint**: `GET /v1/analytics/journey/velocity`

**Purpose**: Measure time-to-conversion across customer journeys

**Metrics**:
- Average journey duration
- Fastest conversion paths
- Bottleneck identification
- Velocity trends

#### Custom Reports

**Endpoint**: `POST /v1/analytics/reports/custom`

**Purpose**: Create custom analytics reports with flexible queries

**Features**:
- Custom date ranges
- Multiple metric selection
- Filtering options
- Export capabilities

---

## User & Access Management

### 1. User Manager

**Location**: `/settings` → User Manager Tab

**Features**:
- Create, edit, and delete users
- Bulk user upload (CSV)
- User profile management
- Team assignment
- Role assignment
- User type management (Product User / Observer)

**User Fields**:
- First Name (required)
- Last Name (optional)
- Email (required, unique)
- Phone (optional)
- Role
- Team
- Sub-team
- User Type (product_user / observer)
- Location
- Password (editable)

**User Types**:
- **Product User**: Can access features based on role permissions
- **Observer**: No functional access, view-only

**API Endpoints**:
- `POST /v1/users` - Create user
- `POST /v1/users/bulk` - Bulk create users
- `GET /v1/users` - List all users
- `GET /v1/users/:id` - Get user details
- `PUT /v1/users/:id` - Update user
- `DELETE /v1/users/:id` - Delete user
- `POST /v1/users/login` - User authentication

**Bulk Upload Format** (CSV):
```
email,name,phone,role_id,team_id,user_type,location
user1@example.com,John Doe,+1234567890,1,1,product_user,USA
user2@example.com,Jane Smith,,2,2,product_user,UK
```

---

### 2. Role Manager

**Location**: `/settings` → Role Manager Tab

**Features**:
- Create custom roles
- Assign permissions to roles
- Role-based access control (RBAC)
- Permission groups organization
- Default roles (Admin, Agent)

**Role Fields**:
- Name
- Description
- Permissions (code_names array)
- Allowed Team IDs (optional)
- Can be edited flag
- Is default flag

**Permission Groups**:
1. **Core Analytics**: View/edit agents, vendors, intents, journey
2. **Advanced Analytics**: MMM, ABM, lead scoring, cohorts, real-time
3. **Platform Management**: Integrations, experiments, reports
4. **User Management**: Users, roles CRUD operations
5. **Data Management**: Data ingestion, attribution, export
6. **Reporting**: View, create, edit, delete reports

**API Endpoints**:
- `POST /v1/roles` - Create role
- `GET /v1/roles` - List all roles
- `GET /v1/roles/:id` - Get role details
- `PUT /v1/roles/:id` - Update role
- `DELETE /v1/roles/:id` - Delete role
- `GET /v1/permissions` - List all permissions

**Default Permissions**:
- `analytics.agents.view` - View Agents Dashboard
- `analytics.vendors.view` - View Vendors Dashboard
- `analytics.intents.view` - View Intents Dashboard
- `users.view` - View users
- `users.create` - Create users
- `users.edit` - Edit users
- `users.delete` - Delete users
- And 50+ more permissions

---

### 3. Profile Management

**Location**: `/profile`

**Features**:
- View user profile
- Edit name (first name, last name)
- Change password
- Update timezone
- View organizational details (role, team, etc.)

**Profile Sections**:
1. **Primary Details**:
   - Name (editable)
   - Email (read-only)
   - Phone
   - Location

2. **Organizational Details**:
   - Role
   - Team
   - Sub-team
   - Manager
   - User Type

3. **Account Settings**:
   - Password (changeable)
   - Timezone (editable)
   - Account status

**API Endpoints**:
- `GET /v1/users/:id` - Get user profile
- `PUT /v1/users/:id` - Update profile (name, password, timezone)

---

## Team Management

### 1. Team Manager

**Location**: `/settings` → Team Manager Tab

**Features**:
- Create parent teams
- Create sub-teams (hierarchical structure)
- Assign users to teams and sub-teams
- Team-vendor association
- Collapsible sub-team view
- Safe team deletion with member transfer

**Team Structure**:
```
Parent Team
  ├── Sub-team 1
  ├── Sub-team 2
  └── Sub-team 3
```

**Team Fields**:
- Name (required)
- Description (optional)
- Vendor (optional)
- Parent Team (for sub-teams)

**Features**:
- **Collapsible Sub-teams**: Expand/collapse to view sub-teams
- **Add Users**: Assign users to teams or sub-teams
- **Create Sub-team**: Create sub-teams directly from parent team
- **Safe Deletion**: Transfer members to another team before deletion

**User Assignment Rules**:
- One user can only be in one team at a time
- Users can be assigned to parent teams or sub-teams
- Assignment automatically removes user from previous team

**API Endpoints**:
- `POST /v1/teams` - Create team
- `GET /v1/teams` - List all teams
- `GET /v1/teams/:id` - Get team details
- `PUT /v1/teams/:id` - Update team
- `DELETE /v1/teams/:id` - Delete team (with member transfer option)
- `POST /v1/teams/:id/members` - Add team members

---

## Attribution Engine

### 1. Attribution Models

The platform supports five attribution models:

#### FIRST_TOUCH
- **Description**: 100% credit to first interaction
- **Best For**: Understanding initial customer acquisition
- **Use Case**: Measuring top-of-funnel effectiveness

#### LAST_TOUCH
- **Description**: 100% credit to last interaction
- **Best For**: Understanding final conversion drivers
- **Use Case**: Measuring bottom-of-funnel effectiveness

#### LINEAR
- **Description**: Equal credit to all interactions
- **Best For**: Fair distribution across touchpoints
- **Use Case**: Recognizing all touchpoints equally

#### TIME_DECAY
- **Description**: More credit to recent interactions (exponential decay)
- **Best For**: Balancing awareness and conversion
- **Use Case**: Most common model for multi-touch attribution

#### AI_WEIGHTED
- **Description**: Credit based on purchase probability and intent
- **Best For**: Data-driven attribution using ML
- **Use Case**: Most accurate attribution when ML data available

### 2. Creating Attribution Runs

**Endpoint**: `POST /v1/attribution/runs`

**Request Body**:
```json
{
  "model_code": "TIME_DECAY",
  "name": "Q1 2024 Attribution",
  "config": {
    "time_window_hours": 72,
    "include_channels": ["call", "email", "chat"],
    "event_types": ["purchase", "signup"],
    "min_purchase_amount": 100.00
  }
}
```

**Configuration Options**:
- `time_window_hours`: Lookback window (default: 72 hours)
- `include_channels`: Channels to include (optional)
- `event_types`: Conversion types to include (optional)
- `min_purchase_amount`: Minimum conversion amount (optional)

### 3. Executing Attribution

**Endpoint**: `POST /v1/attribution/runs/:run_id/execute`

**Process**:
1. Finds all conversions matching criteria
2. For each conversion, finds interactions within time window
3. Applies selected attribution model
4. Calculates weights and attributed amounts
5. Stores attribution results

**Attribution Results Include**:
- Conversion event ID
- Interaction ID
- Agent ID
- Team ID
- Vendor ID
- Attribution weight (0.0-1.0)
- Attributed amount
- Is primary touch flag

### 4. Viewing Attribution Results

**Endpoint**: `GET /v1/attribution/runs/:run_id`

**Returns**:
- Run status (pending, running, completed, failed)
- Start and completion times
- Number of conversions processed
- Configuration used

---

## Advanced Analytics

### 1. Account-Based Marketing (ABM)

**Location**: `/abm`

**Features**:
- Target account management
- Account engagement tracking
- Account scoring
- Revenue by account
- Engagement metrics

**API Endpoints**:
- `POST /v1/abm/accounts` - Create account
- `GET /v1/abm/accounts` - List accounts
- `GET /v1/abm/accounts/:id` - Get account details
- `GET /v1/abm/accounts/:id/summary` - Account summary
- `POST /v1/abm/accounts/engagements` - Track engagement
- `GET /v1/abm/insights/target-accounts` - Target account insights

**Account Fields**:
- Name
- Domain
- Industry
- Employee Count
- Annual Revenue
- Location
- Lifecycle Stage (MQL, SQL, Opportunity, Customer)
- Engagement Score
- Lead Score

---

### 2. Lead Scoring & Predictive Analytics

**Location**: `/lead-scoring`

**Features**:
- Automated lead scoring
- Purchase probability prediction
- High-value lead identification
- Scoring model management
- Predictive insights

**API Endpoints**:
- `POST /v1/leads/customers/:customer_id/score` - Calculate lead score
- `GET /v1/leads/high-value` - Get high-value leads
- `POST /v1/leads/predictions` - Create prediction

**Scoring Factors**:
- Interaction history
- Intent signals
- Engagement level
- Purchase probability
- Account attributes

---

### 3. Cohort Analysis

**Location**: `/cohorts`

**Features**:
- Customer cohort segmentation
- Retention analysis
- Revenue by cohort
- Cohort comparison
- Retention curves

**API Endpoints**:
- `POST /v1/cohorts/compute` - Compute cohort metrics
- `GET /v1/cohorts/segments/:segment_id/retention` - Get retention curve

**Cohort Types**:
- Acquisition cohorts (by signup date)
- Behavioral cohorts (by action)
- Revenue cohorts (by first purchase)

---

### 4. Marketing Mix Modeling (MMM)

**Location**: `/mmm`

**Features**:
- Marketing spend optimization
- Channel effectiveness analysis
- Budget allocation recommendations
- ROI by channel
- Model creation and execution

**API Endpoints**:
- `POST /v1/mmm/run` - Run MMM analysis
- `GET /v1/mmm/models` - List MMM models
- `GET /v1/mmm/models/:model_id/results` - Get model results

**MMM Capabilities**:
- Multi-channel attribution
- Budget optimization
- Channel mix recommendations
- Effectiveness measurement

---

### 5. Real-Time Analytics

**Location**: `/realtime`

**Features**:
- Live metrics dashboard
- Real-time event streaming
- Alert management
- System health monitoring
- Performance metrics

**API Endpoints**:
- `GET /v1/analytics/realtime/metrics` - Get real-time metrics
- `GET /v1/realtime/alerts` - Get alerts
- `POST /v1/realtime/alerts/:id/acknowledge` - Acknowledge alert

**Real-Time Metrics**:
- Active interactions
- Conversions today
- Revenue today
- System performance
- Data quality scores

---

## Integrations

### 1. Integrations Dashboard

**Location**: `/integrations`

**Features**:
- Connect CRM platforms (Salesforce, HubSpot, Marketo)
- Connect ad platforms (Google Ads, LinkedIn, Facebook)
- Sync configuration
- Sync status monitoring
- Integration health checks

**API Endpoints**:
- `POST /v1/integrations` - Create integration
- `GET /v1/integrations` - List integrations
- `POST /v1/integrations/:id/sync` - Sync integration

**Supported Platforms**:
- **CRM**: Salesforce, HubSpot, Marketo
- **Ad Platforms**: Google Ads, LinkedIn Ads, Facebook Ads
- **Analytics**: Google Analytics, Adobe Analytics
- **Email**: Mailchimp, SendGrid

---

### 2. Live Call Flow Integration

**Location**: Webhook endpoints

**Features**:
- Real-time call ingestion
- Automatic customer identification
- Intent detection integration
- Transcript storage
- Call outcome tracking

**Webhook Endpoints**:
- `POST /v1/webhooks/convin` - Convin webhook
- `POST /v1/webhooks/telephony` - Generic telephony webhook

**Supported Events**:
- `call.started` - Call initiation
- `call.ended` - Call completion
- `call.transcript.updated` - Transcript available
- `call.intent.detected` - Intent detected

**See**: [TELEPHONY_INTEGRATION_GUIDE.md](./TELEPHONY_INTEGRATION_GUIDE.md) for detailed integration instructions.

---

## Real-Time Features

### 1. Event Streaming

**Endpoint**: `POST /v1/events`

**Purpose**: Ingest real-time events for immediate processing

**Event Types**:
- Page views
- Form submissions
- Button clicks
- Video views
- Downloads
- Custom events

**Processing**:
- Events processed in real-time
- Immediate customer journey updates
- Real-time analytics updates

---

### 2. Alerts & Monitoring

**Features**:
- System alerts
- Data quality alerts
- Performance alerts
- Custom alert rules
- Alert acknowledgment

**Alert Types**:
- **Data Quality**: Missing data, invalid formats
- **Performance**: Slow queries, high latency
- **Business**: Revenue drops, conversion anomalies
- **System**: Database issues, service failures

**API Endpoints**:
- `GET /v1/realtime/alerts` - Get alerts
- `POST /v1/realtime/alerts/:id/acknowledge` - Acknowledge alert

---

### 3. Fraud Detection

**Endpoint**: `GET /v1/fraud/detect`

**Purpose**: Detect fraudulent interactions and conversions

**Detection Methods**:
- Pattern analysis
- Anomaly detection
- Velocity checks
- Duplicate detection
- Behavioral analysis

**Endpoint**: `GET /v1/fraud/incidents`

**Returns**: List of detected fraud incidents

---

### 4. Data Quality

**Endpoint**: `GET /v1/quality/scores`

**Purpose**: Calculate data quality scores

**Metrics**:
- Completeness score
- Accuracy score
- Consistency score
- Timeliness score
- Overall quality score

---

## Additional Features

### 1. Custom Reports

**Location**: Reports section

**Features**:
- Create custom reports
- Save report templates
- Schedule reports
- Export reports (CSV, PDF)
- Share reports

**API Endpoints**:
- `POST /v1/reports` - Create report
- `GET /v1/reports` - List reports
- `POST /v1/reports/:id/execute` - Execute report

---

### 2. A/B Testing & Experiments

**Features**:
- Create experiments
- Variant assignment
- Results tracking
- Statistical significance
- Winner determination

**API Endpoints**:
- `POST /v1/experiments` - Create experiment
- `GET /v1/experiments/:id/results` - Get results

---

### 3. Feature Flags

**Features**:
- Gradual feature rollout
- Segment targeting
- Percentage rollout
- Feature enable/disable

**API Endpoints**:
- `POST /v1/features/flags` - Create feature flag
- `GET /v1/features/flags` - List feature flags

---

## API Reference

### Authentication

**Methods**:
1. **API Key**: `X-API-Key: <your-api-key>`
2. **JWT Token**: `Authorization: Bearer <token>`
3. **Tenant ID**: `X-Tenant-ID: <tenant-id>` (required for all requests)

### Base URL

- Development: `http://localhost:8080/v1`
- Production: `https://api.yourdomain.com/v1`

### Common Response Format

**Success Response**:
```json
{
  "data": { ... },
  "message": "Success message"
}
```

**Error Response**:
```json
{
  "error": "Error message",
  "code": "ERROR_CODE"
}
```

### Status Codes

- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Invalid request
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Duplicate resource
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error

---

## Complete API Endpoint List

### Data Ingestion
- `POST /v1/interactions` - Ingest interaction
- `POST /v1/conversions` - Track conversion
- `POST /v1/events` - Ingest event
- `POST /v1/page-views` - Track page view

### Customer & Journey
- `GET /v1/customers/:customer_id/journey` - Get customer journey

### Attribution
- `POST /v1/attribution/runs` - Create attribution run
- `GET /v1/attribution/runs/:run_id` - Get attribution run
- `POST /v1/attribution/runs/:run_id/execute` - Execute attribution

### Analytics
- `GET /v1/analytics/agents/revenue` - Agent revenue
- `GET /v1/analytics/vendors/comparison` - Vendor comparison
- `GET /v1/analytics/intents/revenue` - Intent revenue
- `GET /v1/analytics/funnel/stages` - Funnel metrics
- `GET /v1/analytics/content/engagement` - Content engagement
- `GET /v1/analytics/channels/roi` - Multi-channel ROI
- `GET /v1/analytics/journey/velocity` - Journey velocity
- `POST /v1/analytics/reports/custom` - Custom report
- `GET /v1/analytics/realtime/metrics` - Real-time metrics

### ABM
- `POST /v1/abm/accounts` - Create account
- `GET /v1/abm/accounts` - List accounts
- `GET /v1/abm/accounts/:id` - Get account
- `GET /v1/abm/accounts/:id/summary` - Account summary
- `POST /v1/abm/accounts/engagements` - Track engagement
- `GET /v1/abm/insights/target-accounts` - Target insights

### Lead Scoring
- `POST /v1/leads/customers/:customer_id/score` - Calculate score
- `GET /v1/leads/high-value` - High-value leads
- `POST /v1/leads/predictions` - Create prediction

### Cohorts
- `POST /v1/cohorts/compute` - Compute metrics
- `GET /v1/cohorts/segments/:segment_id/retention` - Retention curve

### MMM
- `POST /v1/mmm/run` - Run MMM analysis
- `GET /v1/mmm/models` - List models
- `GET /v1/mmm/models/:model_id/results` - Get results

### Real-Time
- `GET /v1/realtime/alerts` - Get alerts
- `POST /v1/realtime/alerts/:id/acknowledge` - Acknowledge alert

### Fraud & Quality
- `GET /v1/fraud/detect` - Detect fraud
- `GET /v1/fraud/incidents` - Fraud incidents
- `GET /v1/quality/scores` - Data quality scores

### Behavior
- `GET /v1/behavior/sessions/:session_id` - Session details
- `GET /v1/behavior/pages/top` - Top pages

### Integrations
- `POST /v1/integrations` - Create integration
- `GET /v1/integrations` - List integrations
- `POST /v1/integrations/:id/sync` - Sync integration

### Reports
- `POST /v1/reports` - Create report
- `GET /v1/reports` - List reports
- `POST /v1/reports/:id/execute` - Execute report

### Experiments
- `POST /v1/experiments` - Create experiment
- `GET /v1/experiments/:id/results` - Get results

### Feature Flags
- `POST /v1/features/flags` - Create flag
- `GET /v1/features/flags` - List flags

### Users
- `POST /v1/users` - Create user
- `POST /v1/users/bulk` - Bulk create
- `POST /v1/users/login` - Login
- `GET /v1/users` - List users
- `GET /v1/users/:id` - Get user
- `PUT /v1/users/:id` - Update user
- `DELETE /v1/users/:id` - Delete user

### Roles
- `POST /v1/roles` - Create role
- `GET /v1/roles` - List roles
- `GET /v1/roles/:id` - Get role
- `PUT /v1/roles/:id` - Update role
- `DELETE /v1/roles/:id` - Delete role

### Permissions
- `GET /v1/permissions` - List permissions

### Teams
- `POST /v1/teams` - Create team
- `GET /v1/teams` - List teams
- `GET /v1/teams/:id` - Get team
- `PUT /v1/teams/:id` - Update team
- `DELETE /v1/teams/:id` - Delete team
- `POST /v1/teams/:id/members` - Add members

### Use Cases
- `POST /v1/use-cases` - Create use case
- `GET /v1/use-cases` - List use cases
- `DELETE /v1/use-cases/:id` - Delete use case

### Vendors
- `GET /v1/vendors` - List vendors

### Webhooks
- `POST /v1/webhooks/convin` - Convin webhook
- `POST /v1/webhooks/telephony` - Generic telephony webhook

### Health
- `GET /health` - Health check
- `GET /ready` - Readiness probe
- `GET /live` - Liveness probe

---

## Use Cases

### 1. Call Center Revenue Attribution

**Scenario**: BPO company wants to attribute revenue to specific agents and teams.

**Steps**:
1. Integrate Convin webhook for call tracking
2. Calls automatically ingested as interactions
3. Track conversions from billing system
4. Run attribution with TIME_DECAY model
5. View agent and team performance in dashboards

**Outcome**: Clear visibility into which agents/teams drive revenue.

---

### 2. Multi-Channel Marketing Attribution

**Scenario**: Company wants to understand ROI across all marketing channels.

**Steps**:
1. Ingest interactions from all channels (calls, emails, ads)
2. Track conversions
3. Run multi-channel attribution
4. Analyze channel ROI
5. Optimize budget allocation

**Outcome**: Data-driven marketing budget decisions.

---

### 3. Enterprise Account Tracking

**Scenario**: B2B company wants to track and score enterprise accounts.

**Steps**:
1. Create target accounts in ABM
2. Track account engagements
3. Calculate account scores
4. Monitor account progression
5. Attribute revenue to accounts

**Outcome**: Better enterprise sales pipeline management.

---

### 4. Lead Qualification

**Scenario**: Sales team wants to prioritize high-value leads.

**Steps**:
1. Ingest all customer interactions
2. Calculate lead scores automatically
3. Identify high-value leads
4. Predict purchase probability
5. Route leads to appropriate sales reps

**Outcome**: Improved conversion rates and sales efficiency.

---

## Best Practices

### 1. Data Quality
- Always provide accurate timestamps
- Use consistent customer identifiers
- Include all relevant metadata
- Validate data before ingestion

### 2. Attribution
- Choose model based on business goals
- Set appropriate time windows
- Regularly review and adjust models
- Compare multiple models

### 3. User Management
- Use roles for access control
- Assign users to appropriate teams
- Regularly audit user access
- Use bulk upload for efficiency

### 4. Performance
- Use appropriate date ranges in queries
- Cache frequently accessed data
- Monitor system performance
- Optimize database queries

### 5. Security
- Use strong passwords
- Rotate API keys regularly
- Enable webhook signature verification
- Monitor access logs

---

## Support & Resources

### Documentation
- **Telephony Integration**: See [TELEPHONY_INTEGRATION_GUIDE.md](./TELEPHONY_INTEGRATION_GUIDE.md)

### Health Checks
- `GET /health` - Full system health
- `GET /ready` - Readiness for traffic
- `GET /live` - Service liveness

### Logs
- Backend logs: `docker-compose logs backend`
- Frontend logs: `docker-compose logs frontend`
- Application logs: `logs/app.log`

---

## Version Information

**Current Version**: 2.0.0-production

**Features Included**:
- Multi-Touch Attribution
- Account-Based Marketing (ABM)
- Lead Scoring & Predictive Analytics
- Cohort Analysis
- Real-Time Data Streaming
- Fraud Detection
- User Behavior Analytics
- CRM & Ad Platform Integrations
- Custom Reports & Dashboards
- A/B Testing & Experiments
- Feature Flags
- Marketing Mix Modeling (MMM)
- Live Call Flow Integration

---

## Getting Started

1. **Access the Application**: Navigate to `http://localhost:3000`
2. **Login**: Use your credentials or create an account
3. **Configure Teams**: Set up your team structure
4. **Integrate Telephony**: Configure webhooks for call tracking
5. **Track Conversions**: Start tracking conversion events
6. **Run Attribution**: Create and execute attribution runs
7. **View Analytics**: Explore dashboards and reports

For detailed integration instructions, see [TELEPHONY_INTEGRATION_GUIDE.md](./TELEPHONY_INTEGRATION_GUIDE.md).

