# Testing Notifications

This guide explains how to test the notification center feature.

## Prerequisites

1. **Backend server running**: `http://localhost:8080`
2. **Frontend running**: `http://localhost:3000`
3. **Logged in user**: You need to be logged in to the application

## Method 1: Using the Web Interface (Easiest)

1. Open `test_notifications.html` in your browser
2. Enter your **Auth Token** (get it from browser localStorage after logging in)
3. Enter your **Tenant ID** (usually `1` for default tenant)
4. Click any preset button (Error, Warning, Info, Success, Critical) to load a test notification
5. Click **"Create Notification"** to create a single notification
6. Or click **"Create All Presets"** to create all 5 test notifications at once

### Getting Your Auth Token

1. Open your browser's Developer Tools (F12)
2. Go to **Application** tab (Chrome) or **Storage** tab (Firefox)
3. Navigate to **Local Storage** → `http://localhost:3000`
4. Find the `token` key and copy its value

## Method 2: Using the Command Line Script

1. **Get your auth token** (see above)

2. **Run the script**:
   ```bash
   TOKEN=your_token_here ./test_notifications.sh
   ```

   Or set environment variables:
   ```bash
   export TOKEN=your_token_here
   export TENANT_ID=1
   export API_URL=http://localhost:8080/v1
   ./test_notifications.sh
   ```

## Method 3: Using cURL Directly

### Create a Test Notification

```bash
curl -X POST "http://localhost:8080/v1/realtime/alerts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "X-Tenant-ID: 1" \
  -d '{
    "alert_type": "system_error",
    "severity": "error",
    "title": "Database Connection Failed",
    "description": "Unable to connect to primary database.",
    "entity_type": "system"
  }'
```

### Get All Notifications

```bash
curl -X GET "http://localhost:8080/v1/realtime/alerts?limit=20&unresolved=true" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "X-Tenant-ID: 1"
```

### Acknowledge a Notification

```bash
curl -X POST "http://localhost:8080/v1/realtime/alerts/1/acknowledge?by=user@example.com" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "X-Tenant-ID: 1"
```

## Test Notification Types

### 1. Error Alert
- **Severity**: `error`
- **Type**: System errors, critical issues
- **Icon**: Red error icon

### 2. Warning Alert
- **Severity**: `warning`
- **Type**: Data quality issues, performance warnings
- **Icon**: Yellow warning icon

### 3. Info Alert
- **Severity**: `info`
- **Type**: General information, status updates
- **Icon**: Blue info icon

### 4. Success Alert
- **Severity**: `success`
- **Type**: Successful operations, completions
- **Icon**: Green success icon

## What to Check

After creating notifications:

1. **Notification Icon Badge**: 
   - The bell icon in the top right should show a red badge with the unread count

2. **Notification Popover**:
   - Click the notification icon
   - You should see all unread notifications
   - Each notification should have:
     - Appropriate icon (error/warning/info/success)
     - Title and description
     - Timestamp (e.g., "Just now", "2m ago")
     - Blue background for unread items
     - Dot indicator for unread items

3. **Acknowledge Functionality**:
   - Click on an unread notification
   - It should be marked as acknowledged
   - The badge count should decrease
   - The notification background should change

4. **Real-time Updates**:
   - Create a new notification
   - Wait up to 30 seconds (polling interval)
   - The notification should appear automatically

5. **View All Link**:
   - Click "View All Notifications" at the bottom
   - Should navigate to `/realtime` dashboard

## Troubleshooting

### No notifications showing?

1. **Check backend logs**: Ensure the backend is running and receiving requests
2. **Check browser console**: Look for API errors
3. **Verify token**: Make sure your auth token is valid
4. **Check tenant ID**: Ensure you're using the correct tenant ID

### Notifications not updating?

1. **Check polling**: The frontend polls every 30 seconds
2. **Refresh manually**: Click the notification icon to force a refresh
3. **Check API response**: Verify `/v1/realtime/alerts` returns data

### Can't acknowledge notifications?

1. **Check API endpoint**: Verify `POST /v1/realtime/alerts/:id/acknowledge` works
2. **Check permissions**: Ensure your user has proper permissions
3. **Check console**: Look for JavaScript errors

## API Endpoints

- `POST /v1/realtime/alerts` - Create alert
- `GET /v1/realtime/alerts` - Get alerts (with filters)
- `POST /v1/realtime/alerts/:id/acknowledge` - Acknowledge alert

## Example Test Flow

1. Open the application and log in
2. Open `test_notifications.html` in a new tab
3. Copy your auth token from browser localStorage
4. Paste token in the test page
5. Click "Create All Presets"
6. Go back to the main application
7. Click the notification icon (should show badge with count)
8. Verify all 5 notifications appear
9. Click on each notification to acknowledge
10. Verify badge count decreases
11. Create a new notification and wait for auto-update

