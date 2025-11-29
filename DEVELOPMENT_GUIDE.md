# 🚀 Convin Economics - Complete Development Guide

A comprehensive end-to-end guide for building and scaling the revenue attribution platform.

---

## 📋 Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Database Design](#database-design)
3. [Backend Development](#backend-development)
4. [API Development](#api-development)
5. [Frontend Development](#frontend-development)
6. [External Integrations](#external-integrations)
7. [Testing Strategy](#testing-strategy)
8. [Deployment](#deployment)
9. [Scaling & Optimization](#scaling--optimization)
10. [Security Best Practices](#security-best-practices)

---

## 1. Architecture Overview

### System Architecture

```
┌─────────────────┐
│   React Frontend │ ← User Interface (Port 3000)
└────────┬────────┘
         │ HTTP/REST
         ↓
┌─────────────────┐
│   Go Backend    │ ← API Server (Port 8080)
│  (Gin Framework)│
└────────┬────────┘
         │ SQL
         ↓
┌─────────────────┐
│   PostgreSQL    │ ← Database (Port 5432)
└─────────────────┘
         ↑
         │ Sync/ETL
┌─────────────────┐
│  External APIs  │ ← CRM, Ad Platforms, etc.
└─────────────────┘
```

### Technology Stack

**Backend:**
- Language: Go 1.21+
- Framework: Gin (HTTP router)
- Database: sqlx (SQL toolkit)
- Driver: lib/pq (PostgreSQL)

**Frontend:**
- Framework: React 18
- UI Library: Material-UI v5
- State: React Hooks
- HTTP: Axios
- Routing: React Router v6

**Database:**
- PostgreSQL 14+
- Multi-tenant architecture
- JSONB for flexible schemas

---

## 2. Database Design

### 2.1 Core Principles

#### Multi-Tenancy
Every table includes `tenant_id` for data isolation:
```sql
tenant_id BIGINT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE
```

#### Indexing Strategy
```sql
-- Primary indexes
CREATE INDEX idx_interactions_tenant ON interactions(tenant_id);
CREATE INDEX idx_interactions_customer ON interactions(customer_id);
CREATE INDEX idx_interactions_date ON interactions(interaction_date);

-- Composite indexes for common queries
CREATE INDEX idx_interactions_tenant_date ON interactions(tenant_id, interaction_date);
```

### 2.2 Schema Design Process

**Step 1: Identify Core Entities**
```
Tenants → Customers → Interactions → Conversions → Attribution
```

**Step 2: Create Schema File**
```bash
# Location: database/schema.sql
```

**Step 3: Apply Schema**
```bash
psql convin_crae < database/schema.sql
```

### 2.3 Key Tables Structure

#### Tenants (Multi-tenancy)
```sql
CREATE TABLE tenants (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    config JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### Customers (Identity)
```sql
CREATE TABLE customers (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE customer_identifiers (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT REFERENCES customers(id),
    identifier_type VARCHAR(50), -- 'phone', 'email', 'crm_id'
    identifier_value VARCHAR(255),
    UNIQUE(identifier_type, identifier_value)
);
```

#### Interactions (Touchpoints)
```sql
CREATE TABLE interactions (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    customer_id BIGINT REFERENCES customers(id),
    channel_id BIGINT REFERENCES channels(id),
    interaction_date TIMESTAMPTZ NOT NULL,
    duration_seconds INT,
    primary_intent VARCHAR(100),
    secondary_intents TEXT[],
    outcome_prediction VARCHAR(50),
    purchase_probability DECIMAL(5,2),
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### Conversion Events (Revenue)
```sql
CREATE TABLE conversion_events (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    customer_id BIGINT REFERENCES customers(id),
    conversion_date TIMESTAMPTZ NOT NULL,
    event_type VARCHAR(100), -- 'purchase', 'signup', 'demo'
    revenue_amount DECIMAL(15,2),
    currency_code VARCHAR(3),
    funnel_stage VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### 2.4 Migration Strategy

Create incremental migrations:
```bash
database/migrations/
├── 001_initial_schema.sql
├── 002_add_funnel_stages.sql
├── 003_add_abm_features.sql
└── 004_add_mmm_tables.sql
```

**Migration Template:**
```sql
-- Migration: 001_initial_schema.sql
-- Description: Create core tables
-- Date: 2024-11-15

BEGIN;

-- Create tables
CREATE TABLE IF NOT EXISTS tenants (...);
CREATE TABLE IF NOT EXISTS customers (...);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_customers_tenant ON customers(tenant_id);

COMMIT;
```

**Apply Migrations:**
```bash
for file in database/migrations/*.sql; do
  psql convin_crae < "$file"
done
```

---

## 3. Backend Development

### 3.1 Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── api/
│   │   ├── server.go            # Router setup
│   │   └── handlers/
│   │       └── handlers.go      # HTTP handlers
│   ├── models/
│   │   └── models.go            # Data structures
│   └── services/
│       ├── identity.go          # Business logic
│       ├── attribution.go
│       ├── analytics.go
│       └── mmm.go
├── go.mod
└── go.sum
```

### 3.2 Setting Up the Backend

**Step 1: Initialize Go Module**
```bash
cd backend
go mod init github.com/convin/crae
```

**Step 2: Install Dependencies**
```bash
go get github.com/gin-gonic/gin
go get github.com/jmoiron/sqlx
go get github.com/lib/pq
```

**Step 3: Create Main Entry Point**
```go
// cmd/server/main.go
package main

import (
    "log"
    "os"
    
    "github.com/convin/crae/internal/api"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

func main() {
    // Database connection
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        dbURL = "postgres://localhost/convin_crae?sslmode=disable"
    }
    
    db, err := sqlx.Connect("postgres", dbURL)
    if err != nil {
        log.Fatal("Database connection failed:", err)
    }
    defer db.Close()
    
    // Create server
    server := api.NewServer(db, nil)
    
    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("Server starting on port %s", port)
    server.Router.Run(":" + port)
}
```

### 3.3 Service Layer Pattern

**Service Structure:**
```go
// internal/services/analytics.go
package services

import "github.com/jmoiron/sqlx"

type AnalyticsService struct {
    db *sqlx.DB
}

func NewAnalyticsService(db *sqlx.DB) *AnalyticsService {
    return &AnalyticsService{db: db}
}

// Business logic methods
func (s *AnalyticsService) GetRevenueSummary(tenantID int64) (*RevenueSummary, error) {
    // Implementation
}
```

**Why Services?**
- ✅ Separation of concerns
- ✅ Testability
- ✅ Reusability
- ✅ Business logic isolation

### 3.4 Model Definitions

```go
// internal/models/models.go
package models

import (
    "database/sql"
    "time"
)

type Customer struct {
    ID        int64          `db:"id" json:"id"`
    TenantID  int64          `db:"tenant_id" json:"tenant_id"`
    FirstName sql.NullString `db:"first_name" json:"first_name"`
    LastName  sql.NullString `db:"last_name" json:"last_name"`
    Email     sql.NullString `db:"email" json:"email"`
    Phone     sql.NullString `db:"phone" json:"phone"`
    CreatedAt time.Time      `db:"created_at" json:"created_at"`
}

type Interaction struct {
    ID                  int64          `db:"id" json:"id"`
    TenantID            int64          `db:"tenant_id" json:"tenant_id"`
    CustomerID          sql.NullInt64  `db:"customer_id" json:"customer_id"`
    ChannelID           sql.NullInt64  `db:"channel_id" json:"channel_id"`
    InteractionDate     time.Time      `db:"interaction_date" json:"interaction_date"`
    DurationSeconds     sql.NullInt32  `db:"duration_seconds" json:"duration_seconds"`
    PrimaryIntent       sql.NullString `db:"primary_intent" json:"primary_intent"`
    PurchaseProbability sql.NullFloat64 `db:"purchase_probability" json:"purchase_probability"`
    CreatedAt           time.Time      `db:"created_at" json:"created_at"`
}
```

**Key Points:**
- Use `sql.Null*` types for nullable columns
- Use `db` tags for database mapping
- Use `json` tags for API responses

---

## 4. API Development

### 4.1 API Design Principles

**RESTful Conventions:**
```
GET    /v1/customers           # List
GET    /v1/customers/:id       # Retrieve
POST   /v1/customers           # Create
PUT    /v1/customers/:id       # Update
DELETE /v1/customers/:id       # Delete
```

**Multi-tenant Headers:**
```
X-Tenant-ID: 1
```

### 4.2 Handler Implementation

```go
// internal/api/handlers/handlers.go
package handlers

import (
    "net/http"
    "strconv"
    
    "github.com/convin/crae/internal/services"
    "github.com/gin-gonic/gin"
)

type Handlers struct {
    analyticsSvc *services.AnalyticsService
    // ... other services
}

func NewHandlers(analyticsSvc *services.AnalyticsService) *Handlers {
    return &Handlers{analyticsSvc: analyticsSvc}
}

// Get tenant ID from request
func (h *Handlers) getTenantID(c *gin.Context) (int64, error) {
    tenantIDStr := c.GetHeader("X-Tenant-ID")
    if tenantIDStr == "" {
        return 0, fmt.Errorf("missing tenant ID")
    }
    return strconv.ParseInt(tenantIDStr, 10, 64)
}

// Example handler
func (h *Handlers) GetRevenueSummary(c *gin.Context) {
    tenantID, err := h.getTenantID(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tenant ID"})
        return
    }
    
    summary, err := h.analyticsSvc.GetRevenueSummary(tenantID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, summary)
}
```

### 4.3 Router Setup

```go
// internal/api/server.go
package api

import (
    "github.com/convin/crae/internal/api/handlers"
    "github.com/convin/crae/internal/services"
    "github.com/gin-gonic/gin"
    "github.com/jmoiron/sqlx"
)

type Server struct {
    Router *gin.Engine
    db     *sqlx.DB
}

func NewServer(db *sqlx.DB, cfg interface{}) *Server {
    router := gin.Default()
    
    // Initialize services
    analyticsSvc := services.NewAnalyticsService(db)
    attributionSvc := services.NewAttributionService(db)
    
    // Initialize handlers
    h := handlers.NewHandlers(analyticsSvc, attributionSvc)
    
    // API routes
    v1 := router.Group("/v1")
    {
        // Analytics
        v1.GET("/analytics/revenue", h.GetRevenueSummary)
        v1.GET("/analytics/agents", h.GetAgentPerformance)
        
        // Attribution
        v1.POST("/attribution/runs", h.CreateAttributionRun)
        v1.GET("/attribution/runs/:id", h.GetAttributionRun)
        
        // Customers
        v1.GET("/customers/:id/journey", h.GetCustomerJourney)
    }
    
    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    return &Server{Router: router, db: db}
}
```

### 4.4 Error Handling

```go
// internal/api/handlers/errors.go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

func handleError(c *gin.Context, statusCode int, message string, err error) {
    apiErr := APIError{
        Code:    statusCode,
        Message: message,
    }
    
    if err != nil {
        apiErr.Details = err.Error()
    }
    
    c.JSON(statusCode, apiErr)
}

// Usage in handlers
func (h *Handlers) GetCustomer(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        handleError(c, http.StatusBadRequest, "Invalid customer ID", err)
        return
    }
    
    customer, err := h.customerSvc.GetByID(id)
    if err != nil {
        handleError(c, http.StatusInternalServerError, "Failed to fetch customer", err)
        return
    }
    
    c.JSON(http.StatusOK, customer)
}
```

---

## 5. Frontend Development

### 5.1 Project Structure

```
frontend/
├── public/
│   └── index.html
├── src/
│   ├── components/
│   │   ├── Layout.js          # Main layout with sidebar
│   │   └── Logo.js            # Branding component
│   ├── pages/
│   │   ├── Dashboard.js       # Overview page
│   │   ├── AgentsDashboard.js
│   │   ├── MMMDashboard.js
│   │   └── ABMDashboard.js
│   ├── services/
│   │   └── api.js             # API client
│   ├── App.js                 # Main app component
│   └── index.js               # Entry point
└── package.json
```

### 5.2 API Client Setup

```javascript
// src/services/api.js
import axios from 'axios';

const API_BASE = process.env.REACT_APP_API_URL || 'http://localhost:8080/v1';
const TENANT_ID = '1'; // In production, get from auth context

// Create axios instance
const apiClient = axios.create({
  baseURL: API_BASE,
  headers: {
    'Content-Type': 'application/json',
    'X-Tenant-ID': TENANT_ID,
  },
});

// Request interceptor for auth tokens
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Redirect to login
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// API methods
export const analyticsAPI = {
  getRevenueSummary: () => apiClient.get('/analytics/revenue'),
  getAgentPerformance: (params) => apiClient.get('/analytics/agents', { params }),
};

export const attributionAPI = {
  createRun: (data) => apiClient.post('/attribution/runs', data),
  getRun: (id) => apiClient.get(`/attribution/runs/${id}`),
};

export const customerAPI = {
  getJourney: (id) => apiClient.get(`/customers/${id}/journey`),
};

export default apiClient;
```

### 5.3 Dashboard Implementation

```javascript
// src/pages/Dashboard.js
import React, { useState, useEffect } from 'react';
import { Box, Grid, Card, Typography } from '@mui/material';
import { analyticsAPI } from '../services/api';

function Dashboard() {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const response = await analyticsAPI.getRevenueSummary();
      setData(response.data);
    } catch (err) {
      setError(err.message);
      console.error('Error loading data:', err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <Typography>Loading...</Typography>;
  if (error) return <Typography color="error">Error: {error}</Typography>;

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Welcome back, User
      </Typography>
      
      <Grid container spacing={3}>
        <Grid item xs={12} md={3}>
          <Card>
            <Typography variant="h6">Total Revenue</Typography>
            <Typography variant="h3">
              ${data?.totalRevenue?.toLocaleString()}
            </Typography>
          </Card>
        </Grid>
        {/* More cards... */}
      </Grid>
    </Box>
  );
}

export default Dashboard;
```

### 5.4 State Management Pattern

**For Simple State:**
```javascript
// Use React useState
const [data, setData] = useState([]);
```

**For Complex State (Optional):**
```javascript
// Use React Context
import React, { createContext, useContext, useState } from 'react';

const AppContext = createContext();

export function AppProvider({ children }) {
  const [user, setUser] = useState(null);
  const [tenant, setTenant] = useState(null);

  return (
    <AppContext.Provider value={{ user, setUser, tenant, setTenant }}>
      {children}
    </AppContext.Provider>
  );
}

export const useApp = () => useContext(AppContext);
```

---

## 6. External Integrations

### 6.1 Integration Architecture

```
┌──────────────────┐
│ Convin Economics │
└────────┬─────────┘
         │
    ┌────┴────┐
    │ Integration │
    │   Manager   │
    └────┬────┘
         │
    ┌────┴─────────┬──────────┬─────────┐
    │              │          │         │
┌───▼───┐    ┌────▼────┐  ┌──▼──┐  ┌──▼────┐
│Salesforce│  │ HubSpot │  │Google│  │LinkedIn│
└─────────┘    └─────────┘  └Ads──┘  └Ads────┘
```

### 6.2 Integration Service

```go
// internal/services/integrations.go
package services

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type IntegrationService struct {
    db *sqlx.DB
}

type Integration struct {
    ID           int64     `db:"id" json:"id"`
    TenantID     int64     `db:"tenant_id" json:"tenant_id"`
    PlatformName string    `db:"platform_name" json:"platform_name"`
    Type         string    `db:"integration_type" json:"integration_type"`
    Config       JSONB     `db:"config" json:"config"`
    Status       string    `db:"status" json:"status"`
    LastSyncAt   time.Time `db:"last_sync_at" json:"last_sync_at"`
}

// Sync with external platform
func (s *IntegrationService) SyncIntegration(tenantID, integrationID int64) error {
    // Get integration config
    integration, err := s.GetIntegration(tenantID, integrationID)
    if err != nil {
        return err
    }
    
    // Route to appropriate sync handler
    switch integration.PlatformName {
    case "salesforce":
        return s.syncSalesforce(integration)
    case "hubspot":
        return s.syncHubSpot(integration)
    case "google_ads":
        return s.syncGoogleAds(integration)
    default:
        return fmt.Errorf("unsupported platform: %s", integration.PlatformName)
    }
}

// Example: Sync Salesforce
func (s *IntegrationService) syncSalesforce(integration *Integration) error {
    // Get access token from config
    var config map[string]interface{}
    json.Unmarshal([]byte(integration.Config), &config)
    accessToken := config["access_token"].(string)
    
    // Make API call
    client := &http.Client{}
    req, _ := http.NewRequest("GET", 
        "https://yourinstance.salesforce.com/services/data/v56.0/query?q=SELECT+Id,Name+FROM+Account", 
        nil)
    req.Header.Add("Authorization", "Bearer "+accessToken)
    
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    // Process response and sync to database
    // ... implementation ...
    
    // Update last sync time
    s.updateLastSync(integration.ID)
    
    return nil
}
```

### 6.3 OAuth Flow Implementation

```go
// internal/api/handlers/oauth.go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// OAuth callback handler
func (h *Handlers) HandleOAuthCallback(c *gin.Context) {
    platform := c.Param("platform")
    code := c.Query("code")
    state := c.Query("state")
    
    // Verify state token
    if !h.verifyState(state) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
        return
    }
    
    // Exchange code for access token
    token, err := h.exchangeCode(platform, code)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get token"})
        return
    }
    
    // Save integration
    tenantID, _ := h.getTenantID(c)
    err = h.integrationSvc.SaveIntegration(tenantID, platform, token)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save integration"})
        return
    }
    
    c.Redirect(http.StatusFound, "/integrations?success=true")
}
```

### 6.4 Webhook Handlers

```go
// internal/api/handlers/webhooks.go
package handlers

func (h *Handlers) HandleSalesforceWebhook(c *gin.Context) {
    // Verify webhook signature
    signature := c.GetHeader("X-Salesforce-Signature")
    if !h.verifySalesforceSignature(c.Request.Body, signature) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
        return
    }
    
    // Parse webhook payload
    var payload SalesforceWebhookPayload
    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
        return
    }
    
    // Process webhook event
    switch payload.EventType {
    case "account.created":
        h.handleAccountCreated(payload.Data)
    case "opportunity.updated":
        h.handleOpportunityUpdated(payload.Data)
    }
    
    c.JSON(http.StatusOK, gin.H{"status": "processed"})
}
```

---

## 7. Testing Strategy

### 7.1 Backend Testing

**Unit Tests:**
```go
// internal/services/analytics_test.go
package services

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestGetRevenueSummary(t *testing.T) {
    // Setup test database
    db := setupTestDB()
    defer db.Close()
    
    // Seed test data
    seedTestData(db)
    
    // Create service
    svc := NewAnalyticsService(db)
    
    // Test
    summary, err := svc.GetRevenueSummary(1)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, summary)
    assert.Greater(t, summary.TotalRevenue, 0.0)
}
```

**Integration Tests:**
```bash
# Run tests
go test ./... -v

# With coverage
go test ./... -cover
```

### 7.2 API Testing

```bash
# test_apis.sh
#!/bin/bash

BASE_URL="http://localhost:8080/v1"
TENANT_ID="1"

echo "Testing Analytics API..."
curl -H "X-Tenant-ID: $TENANT_ID" \
     "$BASE_URL/analytics/revenue" | jq

echo "Testing Attribution API..."
curl -X POST \
     -H "X-Tenant-ID: $TENANT_ID" \
     -H "Content-Type: application/json" \
     -d '{"model_type":"linear","start_date":"2024-01-01"}' \
     "$BASE_URL/attribution/runs" | jq
```

### 7.3 Frontend Testing

```javascript
// src/pages/Dashboard.test.js
import { render, screen, waitFor } from '@testing-library/react';
import Dashboard from './Dashboard';
import { analyticsAPI } from '../services/api';

jest.mock('../services/api');

test('renders revenue summary', async () => {
  analyticsAPI.getRevenueSummary.mockResolvedValue({
    data: { totalRevenue: 100000 }
  });
  
  render(<Dashboard />);
  
  await waitFor(() => {
    expect(screen.getByText(/\$100,000/)).toBeInTheDocument();
  });
});
```

---

## 8. Deployment

### 8.1 Environment Setup

**Backend Environment Variables:**
```bash
# .env
DATABASE_URL=postgres://user:pass@host:5432/dbname
PORT=8080
JWT_SECRET=your-secret-key
CORS_ORIGINS=http://localhost:3000,https://yourdomain.com
```

**Frontend Environment Variables:**
```bash
# .env.production
REACT_APP_API_URL=https://api.yourdomain.com/v1
REACT_APP_ENV=production
```

### 8.2 Docker Deployment

**Backend Dockerfile:**
```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
```

**Frontend Dockerfile:**
```dockerfile
# Dockerfile
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

**Docker Compose:**
```yaml
# docker-compose.yml
version: '3.8'

services:
  db:
    image: postgres:14-alpine
    environment:
      POSTGRES_DB: convin_crae
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./database/schema.sql:/docker-entrypoint-initdb.d/schema.sql
    ports:
      - "5432:5432"

  backend:
    build: ./backend
    environment:
      DATABASE_URL: postgres://postgres:password@db:5432/convin_crae?sslmode=disable
      PORT: 8080
    ports:
      - "8080:8080"
    depends_on:
      - db

  frontend:
    build: ./frontend
    ports:
      - "80:80"
    depends_on:
      - backend

volumes:
  postgres_data:
```

**Deploy:**
```bash
docker-compose up -d
```

---

## 9. Scaling & Optimization

### 9.1 Database Optimization

**Connection Pooling:**
```go
// Configure database pool
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

**Query Optimization:**
```sql
-- Use EXPLAIN ANALYZE
EXPLAIN ANALYZE
SELECT * FROM interactions 
WHERE tenant_id = 1 
AND interaction_date >= '2024-01-01';

-- Add covering indexes
CREATE INDEX idx_interactions_covering 
ON interactions(tenant_id, interaction_date) 
INCLUDE (customer_id, channel_id);
```

### 9.2 Caching Strategy

```go
// Use Redis for caching
import "github.com/go-redis/redis/v8"

type CachedService struct {
    db    *sqlx.DB
    redis *redis.Client
}

func (s *CachedService) GetRevenueSummary(tenantID int64) (*RevenueSummary, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("revenue:summary:%d", tenantID)
    cached, err := s.redis.Get(ctx, cacheKey).Result()
    if err == nil {
        var summary RevenueSummary
        json.Unmarshal([]byte(cached), &summary)
        return &summary, nil
    }
    
    // Cache miss - query database
    summary, err := s.queryRevenueSummary(tenantID)
    if err != nil {
        return nil, err
    }
    
    // Store in cache
    data, _ := json.Marshal(summary)
    s.redis.Set(ctx, cacheKey, data, 5*time.Minute)
    
    return summary, nil
}
```

### 9.3 Load Balancing

```nginx
# nginx.conf
upstream backend {
    server backend1:8080;
    server backend2:8080;
    server backend3:8080;
}

server {
    listen 80;
    
    location /api/ {
        proxy_pass http://backend;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

---

## 10. Security Best Practices

### 10.1 Authentication & Authorization

**JWT Implementation:**
```go
// internal/middleware/auth.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(401, gin.H{"error": "Missing token"})
            c.Abort()
            return
        }
        
        // Validate token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })
        
        if err != nil || !token.Valid {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        // Set user context
        claims := token.Claims.(jwt.MapClaims)
        c.Set("user_id", claims["user_id"])
        c.Set("tenant_id", claims["tenant_id"])
        
        c.Next()
    }
}
```

### 10.2 SQL Injection Prevention

**Always use parameterized queries:**
```go
// ❌ BAD - Vulnerable to SQL injection
query := fmt.Sprintf("SELECT * FROM customers WHERE email = '%s'", email)

// ✅ GOOD - Safe parameterized query
query := "SELECT * FROM customers WHERE email = $1"
db.Query(query, email)
```

### 10.3 CORS Configuration

```go
// Configure CORS
import "github.com/gin-contrib/cors"

router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"https://yourdomain.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Tenant-ID"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}))
```

### 10.4 Rate Limiting

```go
// internal/middleware/ratelimit.go
import "github.com/ulule/limiter/v3"

func RateLimitMiddleware() gin.HandlerFunc {
    rate := limiter.Rate{
        Period: 1 * time.Minute,
        Limit:  100,
    }
    
    store := memory.NewStore()
    limiter := limiter.New(store, rate)
    
    return func(c *gin.Context) {
        context, err := limiter.Get(c, c.ClientIP())
        if err != nil {
            c.JSON(500, gin.H{"error": "Internal error"})
            c.Abort()
            return
        }
        
        if context.Reached {
            c.JSON(429, gin.H{"error": "Rate limit exceeded"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

---

## 📚 Additional Resources

### Documentation
- [PostgreSQL Docs](https://www.postgresql.org/docs/)
- [Go Documentation](https://go.dev/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [React Documentation](https://react.dev/)
- [Material-UI](https://mui.com/)

### Tools
- **Postman**: API testing
- **pgAdmin**: Database management
- **Docker**: Containerization
- **Git**: Version control

### Best Practices
- Follow RESTful API conventions
- Write tests for critical paths
- Use environment variables for config
- Implement proper logging
- Monitor performance metrics
- Regular security audits

---

## 🎯 Quick Start Checklist

- [ ] Set up PostgreSQL database
- [ ] Apply schema migrations
- [ ] Configure environment variables
- [ ] Start backend server
- [ ] Start frontend application
- [ ] Test API endpoints
- [ ] Load test data
- [ ] Configure integrations
- [ ] Set up monitoring
- [ ] Deploy to production

---

**Need help? Check the existing code in the repository for working examples!**

*Last Updated: November 2024*
*Version: 1.0.0*

