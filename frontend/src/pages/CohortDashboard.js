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
  LinearProgress,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material';
import {
  People,
  TrendingUp,
  TrendingDown,
  Timeline,
} from '@mui/icons-material';
import axios from 'axios';

const API_BASE = 'http://localhost:8080/v1';

function CohortDashboard() {
  const [cohortData, setCohortData] = useState([]);
  const [retentionData, setRetentionData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [selectedSegment, setSelectedSegment] = useState(1);
  const [stats, setStats] = useState({
    totalCustomers: 0,
    retentionRate: 0,
    churnRate: 0,
    avgLifetime: 0,
  });

  useEffect(() => {
    loadData();
  }, [selectedSegment]);

  const loadData = async () => {
    try {
      const [retentionRes] = await Promise.all([
        axios.get(`${API_BASE}/cohorts/segments/${selectedSegment}/retention`, {
          headers: { 'X-Tenant-ID': '1' },
        }),
      ]);
      
      setRetentionData(retentionRes.data.retention || []);
      
      // Calculate stats from retention data
      const totalCustomers = retentionData.reduce((sum, r) => sum + (r.customers || 0), 0);
      const avgRetention = retentionData.length > 0
        ? retentionData.reduce((sum, r) => sum + (r.retention_rate || 0), 0) / retentionData.length
        : 0;
      
      setStats({
        totalCustomers,
        retentionRate: avgRetention.toFixed(1),
        churnRate: (100 - avgRetention).toFixed(1),
        avgLifetime: 8.4, // Mock
      });
      
      setLoading(false);
    } catch (error) {
      console.error('Error loading cohort data:', error);
      setLoading(false);
    }
  };

  const getRetentionColor = (rate) => {
    if (rate >= 80) return 'success';
    if (rate >= 50) return 'warning';
    return 'error';
  };

  return (
    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" component="h1" fontWeight="bold">
          Cohort Analysis
        </Typography>
        <FormControl sx={{ minWidth: 200 }}>
          <InputLabel>Segment</InputLabel>
          <Select
            value={selectedSegment}
            label="Segment"
            onChange={(e) => setSelectedSegment(e.target.value)}
          >
            <MenuItem value={1}>High Value</MenuItem>
            <MenuItem value={2}>Enterprise</MenuItem>
            <MenuItem value={3}>SMB</MenuItem>
            <MenuItem value={4}>Trial Users</MenuItem>
          </Select>
        </FormControl>
      </Box>

      {loading && <LinearProgress sx={{ mb: 3 }} />}

      {/* Key Metrics */}
      <Grid container spacing={3} mb={3}>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <People color="primary" sx={{ mr: 1 }} />
                <Typography variant="h6">{stats.totalCustomers}</Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Total Customers
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <TrendingUp color="success" sx={{ mr: 1 }} />
                <Typography variant="h6">{stats.retentionRate}%</Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Retention Rate
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <TrendingDown color="error" sx={{ mr: 1 }} />
                <Typography variant="h6">{stats.churnRate}%</Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Churn Rate
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <Timeline color="info" sx={{ mr: 1 }} />
                <Typography variant="h6">{stats.avgLifetime}m</Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Avg Customer Lifetime
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Retention Curve */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Retention Curve
          </Typography>
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Period</TableCell>
                  <TableCell align="right">Customers</TableCell>
                  <TableCell align="right">Retention Rate</TableCell>
                  <TableCell>Trend</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {retentionData.map((row, idx) => (
                  <TableRow key={idx}>
                    <TableCell>Period {row.period || idx + 1}</TableCell>
                    <TableCell align="right">{row.customers || 0}</TableCell>
                    <TableCell align="right">
                      <Box display="flex" alignItems="center" justifyContent="flex-end">
                        <Typography
                          variant="body2"
                          fontWeight="bold"
                          color={`${getRetentionColor(row.retention_rate)}.main`}
                          sx={{ mr: 1 }}
                        >
                          {row.retention_rate?.toFixed(1)}%
                        </Typography>
                        <LinearProgress
                          variant="determinate"
                          value={row.retention_rate || 0}
                          sx={{ width: 100 }}
                          color={getRetentionColor(row.retention_rate)}
                        />
                      </Box>
                    </TableCell>
                    <TableCell>
                      {idx > 0 && (
                        <>
                          {row.retention_rate > retentionData[idx - 1]?.retention_rate ? (
                            <TrendingUp color="success" />
                          ) : (
                            <TrendingDown color="error" />
                          )}
                        </>
                      )}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </CardContent>
      </Card>

      {/* Cohort Matrix */}
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Cohort Analysis Matrix
          </Typography>
          <Typography variant="body2" color="text.secondary" mb={2}>
            Retention rates by cohort over time
          </Typography>
          {/* Add cohort matrix visualization here */}
          <Box
            sx={{
              p: 4,
              textAlign: 'center',
              bgcolor: 'grey.100',
              borderRadius: 1,
            }}
          >
            <Typography variant="body2" color="text.secondary">
              Cohort matrix visualization coming soon
            </Typography>
          </Box>
        </CardContent>
      </Card>
    </Container>
  );
}

export default CohortDashboard;

