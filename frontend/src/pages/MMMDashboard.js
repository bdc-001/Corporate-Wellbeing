import React, { useState, useEffect } from 'react';
import {
  Box,
  Container,
  Typography,
  Grid,
  Card,
  CardContent,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
  LinearProgress,
  Alert,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  MenuItem,
} from '@mui/material';
import {
  TrendingUp,
  Assessment,
  EmojiEvents,
  Warning,
} from '@mui/icons-material';
import axios from 'axios';

const API_BASE = 'http://localhost:8080/v1';

function MMMDashboard() {
  const [models, setModels] = useState([]);
  const [selectedModel, setSelectedModel] = useState(null);
  const [loading, setLoading] = useState(true);
  const [openDialog, setOpenDialog] = useState(false);
  const [newModel, setNewModel] = useState({
    model_name: '',
    granularity: 'weekly',
    target_metric: 'revenue',
  });

  useEffect(() => {
    loadModels();
  }, []);

  const loadModels = async () => {
    try {
      const response = await axios.get(`${API_BASE}/mmm/models`, {
        headers: { 'X-Tenant-ID': '1' },
      });
      setModels(response.data.models || []);
      setLoading(false);
    } catch (error) {
      console.error('Error loading MMM models:', error);
      setLoading(false);
    }
  };

  const loadModelResults = async (modelId) => {
    try {
      const response = await axios.get(`${API_BASE}/mmm/models/${modelId}/results`, {
        headers: { 'X-Tenant-ID': '1' },
      });
      setSelectedModel(response.data);
    } catch (error) {
      console.error('Error loading model results:', error);
    }
  };

  const runNewModel = async () => {
    try {
      setLoading(true);
      const endDate = new Date();
      const startDate = new Date();
      startDate.setMonth(startDate.getMonth() - 3);

      await axios.post(
        `${API_BASE}/mmm/run`,
        {
          ...newModel,
          start_date: startDate.toISOString(),
          end_date: endDate.toISOString(),
          channels: [1, 2, 3, 4],
          include_seasons: true,
          include_trends: true,
        },
        { headers: { 'X-Tenant-ID': '1' } }
      );
      
      setOpenDialog(false);
      loadModels();
    } catch (error) {
      console.error('Error running MMM:', error);
      setLoading(false);
    }
  };

  return (
    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4" component="h1" fontWeight="bold">
          Marketing Mix Modeling (MMM)
        </Typography>
        <Button
          variant="contained"
          color="primary"
          onClick={() => setOpenDialog(true)}
        >
          Run New MMM Analysis
        </Button>
      </Box>

      {loading && <LinearProgress sx={{ mb: 3 }} />}

      {/* Key Metrics */}
      <Grid container spacing={3} mb={3}>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <Assessment color="primary" sx={{ mr: 1 }} />
                <Typography variant="h6">{models.length}</Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Total Models
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
                  {models.filter(m => m.status === 'completed').length}
                </Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Completed
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <EmojiEvents color="warning" sx={{ mr: 1 }} />
                <Typography variant="h6">
                  {selectedModel?.channel_effectiveness?.[0]?.channel_name || 'N/A'}
                </Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Top Performing Channel
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={3}>
          <Card>
            <CardContent>
              <Box display="flex" alignItems="center" mb={1}>
                <Warning color="error" sx={{ mr: 1 }} />
                <Typography variant="h6">
                  {selectedModel?.recommendations?.length || 0}
                </Typography>
              </Box>
              <Typography variant="body2" color="text.secondary">
                Recommendations
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Models List */}
      <Grid container spacing={3}>
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                MMM Models
              </Typography>
              <TableContainer>
                <Table size="small">
                  <TableHead>
                    <TableRow>
                      <TableCell>Model Name</TableCell>
                      <TableCell>Status</TableCell>
                      <TableCell>Action</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {models.map((model) => (
                      <TableRow key={model.id} hover>
                        <TableCell>{model.model_name}</TableCell>
                        <TableCell>
                          <Chip
                            label={model.status}
                            size="small"
                            color={model.status === 'completed' ? 'success' : 'default'}
                          />
                        </TableCell>
                        <TableCell>
                          <Button
                            size="small"
                            onClick={() => loadModelResults(model.id)}
                          >
                            View
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            </CardContent>
          </Card>
        </Grid>

        {/* Model Results */}
        <Grid item xs={12} md={8}>
          {selectedModel ? (
            <>
              <Card sx={{ mb: 3 }}>
                <CardContent>
                  <Typography variant="h6" gutterBottom>
                    Channel Effectiveness
                  </Typography>
                  <TableContainer>
                    <Table>
                      <TableHead>
                        <TableRow>
                          <TableCell>Channel</TableCell>
                          <TableCell align="right">Contribution %</TableCell>
                          <TableCell align="right">ROI</TableCell>
                          <TableCell align="right">Coefficient</TableCell>
                        </TableRow>
                      </TableHead>
                      <TableBody>
                        {selectedModel.channel_effectiveness?.map((channel, idx) => (
                          <TableRow key={idx}>
                            <TableCell>{channel.channel_name}</TableCell>
                            <TableCell align="right">
                              {channel.contribution_percentage.toFixed(1)}%
                            </TableCell>
                            <TableCell align="right">
                              {channel.roi.toFixed(2)}%
                            </TableCell>
                            <TableCell align="right">
                              {channel.coefficient.toFixed(3)}
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </TableContainer>
                </CardContent>
              </Card>

              <Card>
                <CardContent>
                  <Typography variant="h6" gutterBottom>
                    Recommendations
                  </Typography>
                  {selectedModel.recommendations?.map((rec, idx) => (
                    <Alert key={idx} severity="info" sx={{ mb: 1 }}>
                      {rec}
                    </Alert>
                  ))}
                </CardContent>
              </Card>
            </>
          ) : (
            <Card>
              <CardContent>
                <Typography variant="body1" color="text.secondary">
                  Select a model to view results
                </Typography>
              </CardContent>
            </Card>
          )}
        </Grid>
      </Grid>

      {/* New Model Dialog */}
      <Dialog open={openDialog} onClose={() => setOpenDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle sx={{ pb: 1 }}>Run New MMM Analysis</DialogTitle>
        <DialogContent sx={{ pt: 2 }}>
          <TextField
            fullWidth
            label="Model Name"
            value={newModel.model_name}
            onChange={(e) => setNewModel({ ...newModel, model_name: e.target.value })}
            margin="normal"
            required
          />
          <TextField
            fullWidth
            select
            label="Granularity"
            value={newModel.granularity}
            onChange={(e) => setNewModel({ ...newModel, granularity: e.target.value })}
            margin="normal"
            sx={{ mt: 2 }}
          >
            <MenuItem value="daily">Daily</MenuItem>
            <MenuItem value="weekly">Weekly</MenuItem>
            <MenuItem value="monthly">Monthly</MenuItem>
          </TextField>
          <TextField
            fullWidth
            select
            label="Target Metric"
            value={newModel.target_metric}
            onChange={(e) => setNewModel({ ...newModel, target_metric: e.target.value })}
            margin="normal"
            sx={{ mt: 2 }}
          >
            <MenuItem value="revenue">Revenue</MenuItem>
            <MenuItem value="conversions">Conversions</MenuItem>
            <MenuItem value="leads">Leads</MenuItem>
          </TextField>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDialog(false)}>Cancel</Button>
          <Button
            onClick={runNewModel}
            variant="contained"
            color="primary"
            disabled={!newModel.model_name}
          >
            Run Analysis
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
}

export default MMMDashboard;

