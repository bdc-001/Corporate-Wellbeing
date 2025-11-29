import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Card,
  CardContent,
  Grid,
  TextField,
  MenuItem,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
} from '@mui/material';
import DateRangePicker from '../components/DateRangePicker';
import client from '../api/client';

function IntentsDashboard() {
  const [intents, setIntents] = useState([]);
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
    fetchIntents();
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

  const fetchIntents = async () => {
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
      console.log('Fetching intents with filters:', params);
      const response = await client.get('/analytics/intents/revenue', { params });
      console.log('API Response:', response.data);
      setIntents(response.data.intents || []);
    } catch (error) {
      console.error('Error fetching intents:', error);
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

  const formatDuration = (seconds) => {
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}m ${secs}s`;
  };

  const getProfitabilityColor = (score) => {
    if (score > 100) return 'success';
    if (score > 50) return 'warning';
    return 'error';
  };

  return (
    <Box>
      <Typography variant="h4" gutterBottom sx={{ fontWeight: 700, mb: 4 }}>
        Intent-Level Profitability
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

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell><strong>Intent</strong></TableCell>
              <TableCell align="right"><strong>Total Revenue</strong></TableCell>
              <TableCell align="right"><strong>Conversions</strong></TableCell>
              <TableCell align="right"><strong>Avg Handle Time</strong></TableCell>
              <TableCell align="right"><strong>Profitability Score</strong></TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={5} align="center">
                  Loading...
                </TableCell>
              </TableRow>
            ) : intents.length === 0 ? (
              <TableRow>
                <TableCell colSpan={5} align="center">
                  No data available
                </TableCell>
              </TableRow>
            ) : (
              intents.map((intent) => (
                <TableRow key={intent.intent_code} hover>
                  <TableCell>
                    <Chip label={intent.intent_code} color="primary" variant="outlined" />
                  </TableCell>
                  <TableCell align="right">
                    <Typography variant="body1" fontWeight={600} color="primary">
                      {formatCurrency(intent.total_attributed_amount)}
                    </Typography>
                  </TableCell>
                  <TableCell align="right">{intent.total_conversions}</TableCell>
                  <TableCell align="right">
                    {formatDuration(intent.avg_handle_time_seconds)}
                  </TableCell>
                  <TableCell align="right">
                    <Chip
                      label={intent.profitability_score.toFixed(2)}
                      color={getProfitabilityColor(intent.profitability_score)}
                      size="small"
                    />
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </TableContainer>
    </Box>
  );
}

export default IntentsDashboard;

