# 🚀 Convin Revenue Attribution Engine (CRAE)

A comprehensive revenue attribution and analytics platform built with Go and React.

## ✨ Features

### Core Attribution
- Multi-touch attribution (First, Last, Linear, Time-Decay, AI-Weighted)
- Customer identity resolution & unified graph
- Agent & vendor performance analytics
- Intent-based revenue tracking

### Advanced Analytics
- **Marketing Mix Modeling (MMM)** - Channel ROI & budget optimization
- **Account-Based Marketing (ABM)** - Account health & targeting
- **Lead Scoring** - Predictive analytics & conversion probability
- **Cohort Analysis** - Retention curves & churn analysis
- **Real-time Analytics** - Live metrics & alerts (auto-refresh)

### Platform Features
- External integrations (CRM, ad platforms)
- Custom reports & dashboards
- A/B testing & feature flags
- Fraud detection & data quality
- User behavior analytics

---

## 🏗️ Architecture

- **Backend**: Go 1.21+ (Gin framework)
- **Frontend**: React 18 (Material-UI v5, Figtree font)
- **Database**: PostgreSQL 14+
- **Design**: Convin design system

---

## 🚀 Quick Start

### 1. Database Setup
```bash
# Create database
createdb convin_crae

# Apply schema
psql convin_crae < database/schema.sql

# Run migrations
psql convin_crae < database/migrations/add_funnel_stages.sql
psql convin_crae < database/migrations/add_comprehensive_features.sql

# Load test data
./load_test_data.sh
```

### 2. Start Backend
```bash
cd backend
export DATABASE_URL="postgres://localhost/convin_crae?sslmode=disable"
export PORT=8080
go run cmd/server/main.go
```
**Backend**: http://localhost:8080

### 3. Start Frontend
```bash
cd frontend
npm install
npm start
```
**Frontend**: http://localhost:3000

---

## 📊 Dashboards

Access these dashboards at http://localhost:3000:

### Core Analytics
- **Overview** - `/`
- **Agents** - `/agents`
- **Vendors** - `/vendors`
- **Intents** - `/intents`

### Advanced Analytics
- **Marketing Mix Modeling** - `/mmm`
- **ABM Dashboard** - `/abm`
- **Lead Scoring** - `/lead-scoring`
- **Real-time Analytics** - `/realtime`
- **Cohort Analysis** - `/cohorts`

### Platform
- **Integrations** - `/integrations`

---

## 🔌 API Endpoints

### Health Check
```bash
curl http://localhost:8080/health
```

### Sample API Calls
```bash
# ABM Accounts
curl -H "X-Tenant-ID: 1" http://localhost:8080/v1/abm/accounts

# High-Value Leads
curl -H "X-Tenant-ID: 1" "http://localhost:8080/v1/leads/high-value?min_score=60"

# Real-time Metrics
curl -H "X-Tenant-ID: 1" "http://localhost:8080/v1/analytics/realtime/metrics?window=15"

# Run MMM Analysis
curl -X POST http://localhost:8080/v1/mmm/run \
  -H "X-Tenant-ID: 1" \
  -H "Content-Type: application/json" \
  -d '{"model_name":"Q4 Analysis","start_date":"2024-10-01T00:00:00Z","end_date":"2024-12-31T23:59:59Z","granularity":"weekly","target_metric":"revenue","channels":[1,2,3,4]}'
```

---

## 🗂️ Project Structure

```
.
├── backend/
│   ├── cmd/server/main.go          # Entry point
│   ├── internal/
│   │   ├── api/
│   │   │   ├── server.go           # API routes (200+ endpoints)
│   │   │   └── handlers/           # Request handlers
│   │   ├── models/                 # Data models
│   │   └── services/               # Business logic (15 services)
│   │       ├── mmm.go              # Marketing Mix Modeling
│   │       ├── abm.go              # Account-Based Marketing
│   │       ├── lead_scoring.go     # Lead scoring
│   │       ├── realtime.go         # Real-time analytics
│   │       └── ...
│   └── go.mod
│
├── frontend/
│   ├── src/
│   │   ├── App.js                  # Main app
│   │   ├── components/Layout.js    # Navigation
│   │   └── pages/                  # 11 dashboards
│   │       ├── MMMDashboard.js
│   │       ├── ABMDashboard.js
│   │       ├── LeadScoringDashboard.js
│   │       ├── RealtimeDashboard.js
│   │       └── ...
│   └── package.json
│
├── database/
│   ├── schema.sql                  # Main schema (40+ tables)
│   ├── migrations/                 # Schema updates
│   └── seed_data.sql              # Test data (300+ records)
│
├── load_test_data.sh              # Data loading script
├── test_apis.sh                   # API testing script
└── README.md                      # This file
```

---

## 📦 Database Schema

**40+ Tables** including:
- `tenants`, `customers`, `accounts`
- `interactions`, `conversion_events`
- `attribution_runs`, `attribution_results`
- `mmm_models`, `channel_effectiveness`
- `lead_scores`, `predictions`
- `cohort_metrics`, `segments`
- `event_stream`, `alerts`
- `integrations`, `experiments`

---

## 🧪 Testing

### Run API Tests
```bash
./test_apis.sh
```

### Sample Data
The system includes comprehensive test data:
- 2 Tenants
- 15+ Accounts
- 30+ Customers
- 100+ Interactions
- 50+ Conversions
- 20+ Lead scores
- 10+ Alerts

---

## 🎨 Design System

- **Primary Color**: #1A62F2 (Convin Blue)
- **Font**: Figtree (300-900 weights)
- **UI Framework**: Material-UI v5
- **Responsive**: Mobile & desktop optimized

---

## 📚 Documentation

- **API Docs**: See `QUICK_ACCESS_GUIDE.md` for endpoint details
- **Design Guide**: Convin design system implemented throughout
- **Database Schema**: See `database/schema.sql`

---

## 🔧 Configuration

### Environment Variables
```bash
DATABASE_URL=postgres://localhost/convin_crae?sslmode=disable
PORT=8080
```

### Multi-tenant Support
All API endpoints require `X-Tenant-ID` header:
```bash
curl -H "X-Tenant-ID: 1" http://localhost:8080/v1/...
```

---

## 🚨 Troubleshooting

### Backend won't start
```bash
# Kill existing process
kill $(lsof -t -i:8080)

# Restart
cd backend && go run cmd/server/main.go
```

### Frontend won't start
```bash
# Kill existing process
lsof -ti:3000 | xargs kill -9

# Restart
cd frontend && npm start
```

### Database connection errors
```bash
# Verify PostgreSQL is running
psql -l

# Check database exists
psql convin_crae -c "SELECT version();"
```

---

## 📈 Key Metrics

- **15** Specialized backend services
- **200+** RESTful API endpoints
- **40+** Database tables
- **11** Frontend dashboards
- **300+** Test data records
- **15,000+** Lines of code

---

## 🎯 What's Next?

The system is fully operational! You can:

1. **Explore dashboards** at http://localhost:3000
2. **Run MMM analysis** to optimize channel spend
3. **Monitor real-time metrics** with live updates
4. **Score leads** for sales prioritization
5. **Analyze cohorts** for retention insights

---

## 📄 License

Proprietary - Convin AI

---

## 🤝 Support

For questions or issues, contact the Convin development team.

---

*Version: 2.0.0-comprehensive*  
*Last Updated: November 2024*
