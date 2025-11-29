import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Card,
  CardContent,
  Grid,
  TextField,
  MenuItem,
} from '@mui/material';
import DateRangePicker from '../components/DateRangePicker';
import client from '../api/client';

function VendorsDashboard() {
  const [vendors, setVendors] = useState([]);
  const [loading, setLoading] = useState(true);
  const [dateRange, setDateRange] = useState([new Date('2025-07-30'), new Date('2025-11-12')]);
  const [modelCode, setModelCode] = useState('AI_WEIGHTED');
  const [teamId, setTeamId] = useState('');
  const [agentId, setAgentId] = useState('');
  const [teams, setTeams] = useState([]);
  const [agentsList, setAgentsList] = useState([]);

  useEffect(() => {
    fetchTeams();
    fetchAgentsList();
  }, []);

  useEffect(() => {
    fetchVendors();
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [dateRange, modelCode, teamId, agentId]);

  const fetchTeams = async () => {
    try {
      const response = await client.get('/teams');
      setTeams(response.data.teams || []);
    } catch (error) {
      console.error('Error fetching teams:', error);
    }
  };

  const fetchAgentsList = async () => {
    try {
      const response = await client.get('/users', { params: { role_id: 2 } });
      setAgentsList(response.data.users || []);
    } catch (error) {
      console.error('Error fetching agents list:', error);
    }
  };

  const fetchVendors = async () => {
    setLoading(true);
    try {
      const [fromDate, toDate] = dateRange || [null, null];
      if (!fromDate || !toDate) return;

      const params = {
        from: fromDate.toISOString(),
        to: toDate.toISOString(),
        model_code: modelCode,
      };
      if (teamId) {
        params.team_id = teamId;
      }
      if (agentId) {
        params.agent_id = agentId;
      }
      console.log('Fetching vendors with filters:', params);
      const response = await client.get('/analytics/vendors/comparison', { params });
      console.log('API Response:', response.data);
      setVendors(response.data.vendors || []);
    } catch (error) {
      console.error('Error fetching vendors:', error);
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(amount);
  };

  return (
    <Box>
      <Typography variant="h4" gutterBottom sx={{ fontWeight: 700, mb: 4 }}>
        Vendor Performance Comparison
      </Typography>

      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Grid container spacing={2} alignItems="center">
            <Grid item xs={12} sm={6} md={3}>
              <DateRangePicker
                label="Date Range"
                value={dateRange}
                onChange={setDateRange}
              />
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <TextField
                select
                label="Attribution Model"
                value={modelCode}
                onChange={(e) => setModelCode(e.target.value)}
                fullWidth
              >
                <MenuItem value="FIRST_TOUCH">First Touch</MenuItem>
                <MenuItem value="LAST_TOUCH">Last Touch</MenuItem>
                <MenuItem value="LINEAR">Linear</MenuItem>
                <MenuItem value="TIME_DECAY">Time Decay</MenuItem>
                <MenuItem value="AI_WEIGHTED">AI Weighted</MenuItem>
              </TextField>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <TextField
                select
                label="Team"
                value={teamId}
                onChange={(e) => setTeamId(e.target.value)}
                fullWidth
              >
                <MenuItem value="">All Teams</MenuItem>
                {teams.map((team) => (
                  <MenuItem key={team.id} value={team.id}>
                    {team.name}
                  </MenuItem>
                ))}
              </TextField>
            </Grid>
            <Grid item xs={12} sm={6} md={3}>
              <TextField
                select
                label="Agent"
                value={agentId}
                onChange={(e) => setAgentId(e.target.value)}
                fullWidth
              >
                <MenuItem value="">All Agents</MenuItem>
                {agentsList.map((agent) => (
                  <MenuItem key={agent.id} value={agent.id}>
                    {agent.name}
                  </MenuItem>
                ))}
              </TextField>
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {loading ? (
        <Typography>Loading...</Typography>
      ) : vendors.length === 0 ? (
        <Typography>No data available</Typography>
      ) : (
        <Grid container spacing={3}>
          {vendors.map((vendor) => (
            <Grid item xs={12} md={6} key={vendor.vendor_id}>
              <Card>
                <CardContent>
                  <Typography variant="h6" gutterBottom>
                    {vendor.name}
                  </Typography>
                  <Box sx={{ mt: 2 }}>
                    <Typography variant="body2" color="text.secondary">
                      Total Revenue
                    </Typography>
                    <Typography variant="h5" color="primary" fontWeight={700}>
                      {formatCurrency(vendor.total_attributed_amount)}
                    </Typography>
                  </Box>
                  <Box sx={{ mt: 2 }}>
                    <Typography variant="body2" color="text.secondary">
                      Total Conversions
                    </Typography>
                    <Typography variant="h6">{vendor.total_conversions}</Typography>
                  </Box>
                  <Box sx={{ mt: 2 }}>
                    <Typography variant="body2" color="text.secondary">
                      Avg Conversion Value
                    </Typography>
                    <Typography variant="h6">
                      {formatCurrency(vendor.avg_conversion_value)}
                    </Typography>
                  </Box>
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}
    </Box>
  );
}

export default VendorsDashboard;

