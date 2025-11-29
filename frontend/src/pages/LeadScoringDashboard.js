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
  Avatar,
} from '@mui/material';
import {
  Star,
  TrendingUp,
  People,
  AttachMoney,
} from '@mui/icons-material';
import axios from 'axios';

const API_BASE = 'http://localhost:8080/v1';

function LeadScoringDashboard() {
  const [leads, setLeads] = useState([]);
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState({
    totalLeads: 0,
    avgScore: 0,
    highValue: 0,
    conversionRate: 0,
  });

  useEffect(() => {
    loadLeads();
  }, []);

  const loadLeads = async () => {
    try {
      const response = await axios.get(`${API_BASE}/leads/high-value`, {
        headers: { 'X-Tenant-ID': '1' },
        params: { min_score: 60 },
      });
      
      const leadsData = response.data.leads || [];
      setLeads(leadsData);
      
      // Calculate stats
      const totalLeads = leadsData.length;
      const avgScore = totalLeads > 0
        ? leadsData.reduce((sum, l) => sum + (l.score || 0), 0) / totalLeads
        : 0;
      const highValue = leadsData.filter(l => l.score >= 80).length;
      
      setStats({
        totalLeads,
        avgScore: avgScore.toFixed(1),
        highValue,
        conversionRate: 68.5, // Mock data
      });
      
      setLoading(false);
    } catch (error) {
      console.error('Error loading leads:', error);
      setLoading(false);
    }
  };

  const getScoreColor = (score) => {
    if (score >= 80) return 'success';
    if (score >= 60) return 'warning';
    return 'default';
  };

  const getScoreLabel = (score) => {
    if (score >= 80) return 'Hot';
    if (score >= 60) return 'Warm';
    return 'Cold';
  };

  return (
    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
      <Typography variant="h4" component="h1" fontWeight="bold" mb={3}>
        Lead Scoring & Predictive Analytics
      </Typography>

      {loading && <LinearProgress sx={{ mb: 3 }} />}

      {/* Key Metrics */}
      <Grid container spacing={3} mb={3}>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <People color="primary" sx={{ mr: 1 }} />
                <Typography variant="h6">{stats.totalLeads}</Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Qualified Leads
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <Star color="warning" sx={{ mr: 1 }} />
                <Typography variant="h6">{stats.avgScore}</Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Average Score
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <TrendingUp color="success" sx={{ mr: 1 }} />
                <Typography variant="h6">{stats.highValue}</Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Hot Leads (80+)
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <AttachMoney color="info" sx={{ mr: 1 }} />
                <Typography variant="h6">{stats.conversionRate}%</Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Predicted Conversion
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Leads Table */}
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            High-Value Leads
          </Typography>
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Customer</TableCell>
                  <TableCell align="right">Lead Score</TableCell>
                  <TableCell>Temperature</TableCell>
                  <TableCell align="right">Predicted LTV</TableCell>
                  <TableCell align="right">Conversion Probability</TableCell>
                  <TableCell>Scoring Model</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {leads.map((lead, idx) => (
                  <TableRow key={lead.id || idx} hover>
                    <TableCell>
                      <Box display="flex" alignItems="center">
                        <Avatar sx={{ mr: 2, width: 32, height: 32 }}>
                          {lead.customer_id?.toString().slice(0, 1) || 'C'}
                        </Avatar>
                        <Box>
                          <Typography variant="body2" fontWeight="bold">
                            Customer {lead.customer_id}
                          </Typography>
                          <Typography variant="caption" color="text.secondary">
                            ID: {lead.id}
                          </Typography>
                        </Box>
                      </Box>
                    </TableCell>
                    <TableCell align="right">
                      <Box>
                        <Typography variant="h6" fontWeight="bold">
                          {lead.score}
                        </Typography>
                        <LinearProgress
                          variant="determinate"
                          value={lead.score}
                          sx={{ mt: 0.5 }}
                          color={getScoreColor(lead.score)}
                        />
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={getScoreLabel(lead.score)}
                        size="small"
                        color={getScoreColor(lead.score)}
                      />
                    </TableCell>
                    <TableCell align="right">
                      ${((lead.score * 1000) + 5000).toLocaleString()}
                    </TableCell>
                    <TableCell align="right">
                      {(lead.score * 0.85).toFixed(1)}%
                    </TableCell>
                    <TableCell>
                      <Typography variant="caption" color="text.secondary">
                        Model {lead.model_id || 1}
                      </Typography>
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

export default LeadScoringDashboard;

