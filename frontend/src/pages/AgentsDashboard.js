import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Card,
  CardContent,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  TextField,
  MenuItem,
  Grid,
  Chip,
} from '@mui/material';
import DateRangePicker from '../components/DateRangePicker';
import client from '../api/client';

function AgentsDashboard() {
  const [agents, setAgents] = useState([]);
  const [loading, setLoading] = useState(true);
  // Data range: July 30, 2025 to November 12, 2025 (based on seed data)
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
    fetchAgents();
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
      const response = await client.get('/users', { params: { role_id: 2 } }); // Assuming role_id 2 is for agents
      setAgentsList(response.data.users || []);
    } catch (error) {
      console.error('Error fetching agents list:', error);
    }
  };

  const fetchAgents = async () => {
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
      console.log('Fetching agents with filters:', params);
      const response = await client.get('/analytics/agents/revenue', { params });
      console.log('API Response:', response.data);
      setAgents(response.data.agents || []);
    } catch (error) {
      console.error('Error fetching agents:', error);
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
        Agent Revenue Performance
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
              <TableCell><strong>Agent</strong></TableCell>
              <TableCell><strong>Vendor</strong></TableCell>
              <TableCell><strong>Team</strong></TableCell>
              <TableCell align="right"><strong>Total Revenue</strong></TableCell>
              <TableCell align="right"><strong>Conversions</strong></TableCell>
              <TableCell align="right"><strong>Avg per Interaction</strong></TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={6} align="center">
                  Loading...
                </TableCell>
              </TableRow>
            ) : agents.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} align="center">
                  No data available
                </TableCell>
              </TableRow>
            ) : (
              agents.map((agent) => (
                <TableRow key={agent.agent_id} hover>
                  <TableCell>
                    <Box>
                      <Typography variant="body1" fontWeight={500}>
                        {agent.name}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        {agent.email}
                      </Typography>
                    </Box>
                  </TableCell>
                  <TableCell>
                    <Chip label={agent.vendor_name || 'N/A'} size="small" />
                  </TableCell>
                  <TableCell>{agent.team_name || 'N/A'}</TableCell>
                  <TableCell align="right">
                    <Typography variant="body1" fontWeight={600} color="primary">
                      {formatCurrency(agent.total_attributed_amount)}
                    </Typography>
                  </TableCell>
                  <TableCell align="right">{agent.total_conversions}</TableCell>
                  <TableCell align="right">
                    {formatCurrency(agent.avg_attributed_amount_per_interaction)}
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

export default AgentsDashboard;

