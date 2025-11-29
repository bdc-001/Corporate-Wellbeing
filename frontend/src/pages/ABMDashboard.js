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
  Tabs,
  Tab,
} from '@mui/material';
import {
  Business,
  TrendingUp,
  Star,
  Timeline,
} from '@mui/icons-material';
import axios from 'axios';

const API_BASE = 'http://localhost:8080/v1';

function ABMDashboard() {
  const [accounts, setAccounts] = useState([]);
  const [insights, setInsights] = useState(null);
  const [loading, setLoading] = useState(true);
  const [tabValue, setTabValue] = useState(0);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const [accountsRes, insightsRes] = await Promise.all([
        axios.get(`${API_BASE}/abm/accounts`, {
          headers: { 'X-Tenant-ID': '1' },
        }),
        axios.get(`${API_BASE}/abm/insights/target-accounts`, {
          headers: { 'X-Tenant-ID': '1' },
        }),
      ]);
      
      setAccounts(accountsRes.data.accounts || []);
      setInsights(insightsRes.data);
      setLoading(false);
    } catch (error) {
      console.error('Error loading ABM data:', error);
      setLoading(false);
    }
  };

  const getHealthColor = (health) => {
    if (health >= 80) return 'success';
    if (health >= 50) return 'warning';
    return 'error';
  };

  return (
    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
      <Typography variant="h4" component="h1" fontWeight="bold" mb={3}>
        Account-Based Marketing (ABM)
      </Typography>

      {loading && <LinearProgress sx={{ mb: 3 }} />}

      {/* Key Metrics */}
      <Grid container spacing={3} mb={3}>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <Business color="primary" sx={{ mr: 1 }} />
                <Typography variant="h6">{accounts.length}</Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Total Accounts
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <Star color="warning" sx={{ mr: 1 }} />
                <Typography variant="h6">
                  {insights?.high_value_accounts || 0}
                </Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                High Value Accounts
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <TrendingUp color="success" sx={{ mr: 1 }} />
                <Typography variant="h6">
                  {insights?.avg_engagement_score?.toFixed(1) || 0}
                </Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Avg Engagement Score
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <Timeline color="info" sx={{ mr: 1 }} />
                <Typography variant="h6">
                  ${insights?.pipeline_value?.toLocaleString() || 0}
                </Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Pipeline Value
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Tabs */}
      <Card>
        <Tabs value={tabValue} onChange={(e, v) => setTabValue(v)}>
          <Tab label="All Accounts" />
          <Tab label="Target Accounts" />
          <Tab label="Hot Accounts" />
        </Tabs>

        <CardContent>
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Account Name</TableCell>
                  <TableCell>Lifecycle Stage</TableCell>
                  <TableCell align="right">Contacts</TableCell>
                  <TableCell align="right">Health Score</TableCell>
                  <TableCell align="right">Engagement</TableCell>
                  <TableCell>Status</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {accounts
                  .filter((acc) => {
                    if (tabValue === 0) return true;
                    if (tabValue === 1) return acc.account_score >= 70;
                    if (tabValue === 2) return acc.health_score >= 80;
                    return true;
                  })
                  .map((account) => (
                    <TableRow key={account.id} hover>
                      <TableCell>
                        <Typography variant="body2" fontWeight="bold">
                          {account.account_name}
                        </Typography>
                        <Typography variant="caption" color="text.secondary">
                          {account.domain}
                        </Typography>
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={account.lifecycle_stage || 'Prospect'}
                          size="small"
                          variant="outlined"
                        />
                      </TableCell>
                      <TableCell align="right">
                        {account.contacts_count || 0}
                      </TableCell>
                      <TableCell align="right">
                        <Chip
                          label={`${account.health_score || 0}%`}
                          size="small"
                          color={getHealthColor(account.health_score)}
                        />
                      </TableCell>
                      <TableCell align="right">
                        {account.engagement_score || 0}
                      </TableCell>
                      <TableCell>
                        <Chip
                          label={account.account_score >= 70 ? 'Target' : 'Monitor'}
                          size="small"
                          color={account.account_score >= 70 ? 'primary' : 'default'}
                        />
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

export default ABMDashboard;

