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
  Button,
  LinearProgress,
  Avatar,
} from '@mui/material';
import {
  Cloud,
  Sync,
  CheckCircle,
  Error,
  Refresh,
} from '@mui/icons-material';
import axios from 'axios';

const API_BASE = 'http://localhost:8080/v1';

function IntegrationsDashboard() {
  const [integrations, setIntegrations] = useState([]);
  const [loading, setLoading] = useState(true);
  const [syncing, setSyncing] = useState({});

  useEffect(() => {
    loadIntegrations();
  }, []);

  const loadIntegrations = async () => {
    try {
      const response = await axios.get(`${API_BASE}/integrations`, {
        headers: { 'X-Tenant-ID': '1' },
      });
      setIntegrations(response.data.integrations || []);
      setLoading(false);
    } catch (error) {
      console.error('Error loading integrations:', error);
      setLoading(false);
    }
  };

  const syncIntegration = async (integrationId) => {
    setSyncing({ ...syncing, [integrationId]: true });
    try {
      await axios.post(
        `${API_BASE}/integrations/${integrationId}/sync`,
        {},
        { headers: { 'X-Tenant-ID': '1' } }
      );
      await loadIntegrations();
    } catch (error) {
      console.error('Error syncing integration:', error);
    } finally {
      setSyncing({ ...syncing, [integrationId]: false });
    }
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'connected': return 'success';
      case 'error': return 'error';
      case 'syncing': return 'info';
      default: return 'default';
    }
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case 'connected': return <CheckCircle color="success" />;
      case 'error': return <Error color="error" />;
      case 'syncing': return <Sync color="info" />;
      default: return <Cloud />;
    }
  };

  const getPlatformIcon = (platform) => {
    const icons = {
      salesforce: '☁️',
      hubspot: '🔶',
      marketo: 'M',
      google_ads: 'G',
      facebook_ads: 'F',
      linkedin_ads: 'in',
      ga4: 'GA',
    };
    return icons[platform] || '🔌';
  };

  return (
    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
      <Typography variant="h4" component="h1" fontWeight="bold" mb={3}>
        Integrations
      </Typography>

      {loading && <LinearProgress sx={{ mb: 3 }} />}

      {/* Stats */}
      <Grid container spacing={3} mb={3}>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <Cloud color="primary" sx={{ mr: 1 }} />
                <Typography variant="h6">{integrations.length}</Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Total Integrations
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <CheckCircle color="success" sx={{ mr: 1 }} />
                <Typography variant="h6">
                  {integrations.filter(i => i.status === 'connected').length}
                </Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Connected
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <Sync color="info" sx={{ mr: 1 }} />
                <Typography variant="h6">
                  {integrations.filter(i => i.status === 'syncing').length}
                </Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Syncing
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <Error color="error" sx={{ mr: 1 }} />
                <Typography variant="h6">
                  {integrations.filter(i => i.status === 'error').length}
                </Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Errors
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Integrations Table */}
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Active Integrations
          </Typography>
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Platform</TableCell>
                  <TableCell>Type</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell>Last Sync</TableCell>
                  <TableCell>Data Synced</TableCell>
                  <TableCell>Action</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {integrations.map((integration) => (
                  <TableRow key={integration.id} hover>
                    <TableCell>
                      <Box display="flex" alignItems="center">
                        <Avatar sx={{ mr: 2, bgcolor: 'primary.light' }}>
                          {getPlatformIcon(integration.platform_name)}
                        </Avatar>
                        <Box>
                          <Typography variant="body2" fontWeight="bold">
                            {integration.platform_name}
                          </Typography>
                          <Typography variant="caption" color="text.secondary">
                            {integration.integration_type}
                          </Typography>
                        </Box>
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={integration.integration_type}
                        size="small"
                        variant="outlined"
                      />
                    </TableCell>
                    <TableCell>
                      <Chip
                        icon={getStatusIcon(integration.status)}
                        label={integration.status}
                        size="small"
                        color={getStatusColor(integration.status)}
                      />
                    </TableCell>
                    <TableCell>
                      {integration.last_sync_at
                        ? new Date(integration.last_sync_at).toLocaleString()
                        : 'Never'}
                    </TableCell>
                    <TableCell>
                      <Typography variant="caption" color="text.secondary">
                        {Math.floor(Math.random() * 10000)} records
                      </Typography>
                    </TableCell>
                    <TableCell>
                      <Button
                        size="small"
                        startIcon={syncing[integration.id] ? <Sync /> : <Refresh />}
                        onClick={() => syncIntegration(integration.id)}
                        disabled={syncing[integration.id] || integration.status === 'syncing'}
                      >
                        {syncing[integration.id] ? 'Syncing...' : 'Sync Now'}
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </CardContent>
      </Card>
    </Container>
  );
}

export default IntegrationsDashboard;

