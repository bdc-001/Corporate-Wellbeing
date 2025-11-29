import React, { useState, useEffect } from 'react';
import {
  Box,
  Container,
  Typography,
  Grid,
  Card,
  CardContent,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  LinearProgress,
  Button,
  Alert,
  IconButton,
} from '@mui/material';
import {
  Speed,
  TrendingUp,
  Notifications,
  NotificationImportant,
  Check,
  Refresh,
} from '@mui/icons-material';
import axios from 'axios';

const API_BASE = 'http://localhost:8080/v1';

function RealtimeDashboard() {
  const [metrics, setMetrics] = useState(null);
  const [alerts, setAlerts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [autoRefresh, setAutoRefresh] = useState(true);

  useEffect(() => {
    loadData();
    
    const interval = setInterval(() => {
      if (autoRefresh) {
        loadData();
      }
    }, 5000); // Refresh every 5 seconds

    return () => clearInterval(interval);
  }, [autoRefresh]);

  const loadData = async () => {
    try {
      const [metricsRes, alertsRes] = await Promise.all([
        axios.get(`${API_BASE}/analytics/realtime/metrics`, {
          headers: { 'X-Tenant-ID': '1' },
          params: { window: 15 },
        }),
        axios.get(`${API_BASE}/realtime/alerts`, {
          headers: { 'X-Tenant-ID': '1' },
        }),
      ]);
      
      setMetrics(metricsRes.data);
      setAlerts(alertsRes.data.alerts || []);
      setLoading(false);
    } catch (error) {
      console.error('Error loading realtime data:', error);
      setLoading(false);
    }
  };

  const acknowledgeAlert = async (alertId) => {
    try {
      await axios.post(
        `${API_BASE}/realtime/alerts/${alertId}/acknowledge`,
        {},
        { headers: { 'X-Tenant-ID': '1' } }
      );
      loadData();
    } catch (error) {
      console.error('Error acknowledging alert:', error);
    }
  };

  const getSeverityColor = (severity) => {
    switch (severity) {
      case 'critical': return 'error';
      case 'warning': return 'warning';
      case 'info': return 'info';
      default: return 'default';
    }
  };

  return (
    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" component="h1" fontWeight="bold">
          Real-Time Analytics
        </Typography>
        <Box>
          <Chip
            icon={<Speed />}
            label={autoRefresh ? 'Live' : 'Paused'}
            color={autoRefresh ? 'success' : 'default'}
            onClick={() => setAutoRefresh(!autoRefresh)}
            sx={{ mr: 2 }}
          />
          <IconButton onClick={loadData}>
            <Refresh />
          </IconButton>
        </Box>
      </Box>

      {loading && <LinearProgress sx={{ mb: 3 }} />}

      {/* Real-time Metrics */}
      <Grid container spacing={3} mb={3}>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <Speed color="primary" sx={{ mr: 1 }} />
                <Typography variant="h6">{metrics?.active_users || 0}</Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Active Users (15m)
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <TrendingUp color="success" sx={{ mr: 1 }} />
                <Typography variant="h6">{metrics?.events_per_minute || 0}</Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Events/Min
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <Notifications color="warning" sx={{ mr: 1 }} />
                <Typography variant="h6">
                  {alerts.filter(a => a.status === 'active').length}
                </Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Active Alerts
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <NotificationImportant color="error" sx={{ mr: 1 }} />
                <Typography variant="h6">
                  {alerts.filter(a => a.severity === 'critical').length}
                </Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Critical Alerts
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Active Alerts */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Active Alerts
          </Typography>
          {alerts.length === 0 ? (
            <Alert severity="success">No active alerts</Alert>
          ) : (
            <TableContainer>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Severity</TableCell>
                    <TableCell>Message</TableCell>
                    <TableCell>Triggered At</TableCell>
                    <TableCell>Status</TableCell>
                    <TableCell>Action</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {alerts.map((alert) => (
                    <TableRow key={alert.id} hover>
                      <TableCell>
                        <Chip
                          label={alert.severity}
                          size="small"
                          color={getSeverityColor(alert.severity)}
                        />
                      </TableCell>
                      <TableCell>
                        <Typography variant="body2">{alert.alert_message}</Typography>
                        <Typography variant="caption" color="text.secondary">
                          Rule: {alert.rule_name}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        {new Date(alert.triggered_at).toLocaleString()}
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={alert.status}
                          size="small"
                          variant="outlined"
                        />
                      </TableCell>
                      <TableCell>
                        {alert.status === 'active' && (
                          <Button
                            size="small"
                            startIcon={<Check />}
                            onClick={() => acknowledgeAlert(alert.id)}
                          >
                            Acknowledge
                          </Button>
                        )}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          )}
        </CardContent>
      </Card>

      {/* Recent Events Stream */}
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Event Stream (Last 15 minutes)
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Total Events: {metrics?.total_events || 0}
          </Typography>
          <Box mt={2}>
            <Typography variant="caption" color="text.secondary">
              Streaming live data...
            </Typography>
          </Box>
        </CardContent>
      </Card>
    </Container>
  );
}

export default RealtimeDashboard;

