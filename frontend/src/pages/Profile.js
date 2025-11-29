import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Card,
  CardContent,
  TextField,
  Button,
  Avatar,
  Grid,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  InputAdornment,
  IconButton,
  Alert,
  Snackbar,
  MenuItem,
  Select,
  FormControl,
  InputLabel,
} from '@mui/material';
import {
  Home as HomeIcon,
  Visibility,
  VisibilityOff,
  Edit as EditIcon,
} from '@mui/icons-material';
import { useAuth } from '../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';
import api from '../api/client';

function Profile() {
  const { user: authUser } = useAuth();
  const navigate = useNavigate();
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState({
    firstName: false,
    lastName: false,
  });
  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
    email: '',
    location: '',
    phone: '',
    alternatePhone: '',
    whatsappPhone: '',
  });
  const [passwordDialogOpen, setPasswordDialogOpen] = useState(false);
  const [timezoneDialogOpen, setTimezoneDialogOpen] = useState(false);
  const [passwordData, setPasswordData] = useState({
    currentPassword: '',
    newPassword: '',
    confirmPassword: '',
  });
  const [timezone, setTimezone] = useState('UTC');
  const [showPassword, setShowPassword] = useState({
    current: false,
    new: false,
    confirm: false,
  });
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' });

  useEffect(() => {
    if (authUser && authUser.id) {
      fetchUserProfile();
    }
  }, [authUser]);

  const fetchUserProfile = async () => {
    try {
      setLoading(true);
      const response = await api.get(`/users/${authUser.id}`);
      const userData = response.data.user || response.data;
      
      if (!userData || !userData.id) {
        console.error('Invalid user data received:', response.data);
        showSnackbar('Invalid user data received', 'error');
        return;
      }
      
      // Use first_name and last_name from API, fallback to splitting name if not available
      const firstName = userData.first_name || (userData.name ? userData.name.split(' ')[0] : '') || '';
      const lastName = userData.last_name || (userData.name ? userData.name.split(' ').slice(1).join(' ') : '') || '';

      setUser(userData);
      setFormData({
        firstName,
        lastName,
        email: userData.email || '',
        location: userData.location || '',
        phone: userData.phone || '',
        alternatePhone: '',
        whatsappPhone: '',
      });
      setTimezone(userData.timezone || 'UTC');
    } catch (error) {
      console.error('Error fetching user profile:', error);
      const errorMessage = error.response?.data?.error || error.message || 'Failed to load profile';
      showSnackbar(errorMessage, 'error');
      // Don't set user to null if there's an error, keep showing loading or error state
    } finally {
      setLoading(false);
    }
  };

  const handleSaveName = async () => {
    try {
      await api.put(`/users/${user.id}`, {
        first_name: formData.firstName,
        last_name: formData.lastName || null,
        phone: formData.phone,
        location: formData.location,
        role_id: user.role_id,
        manager_id: user.manager_id,
        auditor_id: user.auditor_id,
        team_id: user.team_id,
        user_type: user.user_type,
      });
      
      setEditing({ firstName: false, lastName: false });
      showSnackbar('Name updated successfully');
      fetchUserProfile();
    } catch (error) {
      showSnackbar(error.response?.data?.error || 'Failed to update name', 'error');
    }
  };

  const handleChangePassword = async () => {
    if (passwordData.newPassword !== passwordData.confirmPassword) {
      showSnackbar('New passwords do not match', 'error');
      return;
    }

    if (passwordData.newPassword.length < 6) {
      showSnackbar('Password must be at least 6 characters', 'error');
      return;
    }

    try {
      // Note: You'll need to implement a password change endpoint
      // For now, we'll use the update user endpoint with a new password field
      await api.put(`/users/${user.id}`, {
        name: user.name,
        phone: user.phone,
        location: user.location,
        role_id: user.role_id,
        manager_id: user.manager_id,
        auditor_id: user.auditor_id,
        team_id: user.team_id,
        user_type: user.user_type,
        new_password: passwordData.newPassword, // Backend should handle this
      });

      setPasswordDialogOpen(false);
      setPasswordData({ currentPassword: '', newPassword: '', confirmPassword: '' });
      showSnackbar('Password changed successfully');
    } catch (error) {
      showSnackbar(error.response?.data?.error || 'Failed to change password', 'error');
    }
  };

  const handleChangeTimezone = async () => {
    try {
      await api.put(`/users/${user.id}`, {
        name: user.name,
        phone: user.phone,
        location: user.location,
        role_id: user.role_id,
        manager_id: user.manager_id,
        auditor_id: user.auditor_id,
        team_id: user.team_id,
        user_type: user.user_type,
        timezone: timezone,
      });

      setTimezoneDialogOpen(false);
      showSnackbar('Timezone updated successfully');
    } catch (error) {
      showSnackbar(error.response?.data?.error || 'Failed to update timezone', 'error');
    }
  };

  const showSnackbar = (message, severity = 'success') => {
    setSnackbar({ open: true, message, severity });
  };

  const getInitials = (name) => {
    if (!name) return 'U';
    const parts = name.trim().split(' ');
    if (parts.length >= 2) {
      return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
    }
    return name[0].toUpperCase();
  };

  const formatDate = (dateString) => {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return date.toLocaleDateString('en-GB', { day: 'numeric', month: 'short', year: 'numeric' });
  };

  if (loading) {
    return (
      <Box sx={{ p: 3 }}>
        <Typography>Loading profile...</Typography>
      </Box>
    );
  }

  if (!user) {
    return (
      <Box sx={{ p: 3 }}>
        <Typography variant="h6" sx={{ mb: 2, color: 'error.main' }}>
          User not found
        </Typography>
        <Typography variant="body2" sx={{ mb: 2, color: 'text.secondary' }}>
          Unable to load user profile. Please try logging in again.
        </Typography>
        <Typography variant="body2" sx={{ color: 'text.secondary', fontSize: '0.75rem' }}>
          User ID: {authUser?.id || 'N/A'}
        </Typography>
      </Box>
    );
  }

  return (
    <Box>
      {/* Header with Avatar and Navigation */}
      <Card
        sx={{
          mb: 3,
          borderRadius: '12px',
          border: '1px solid',
          borderColor: 'rgba(229, 231, 235, 0.8)',
          backgroundColor: '#FFFFFF',
          boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
        }}
      >
        <CardContent sx={{ p: 3 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 3 }}>
            <Avatar
              sx={{
                width: 64,
                height: 64,
                bgcolor: '#3B82F6',
                fontSize: '1.5rem',
                fontWeight: 600,
              }}
            >
              {getInitials(user)}
            </Avatar>
            <Box sx={{ flex: 1 }}>
              <Typography
                variant="h5"
                sx={{
                  fontWeight: 700,
                  color: '#1F2937',
                  mb: 1,
                }}
              >
                {user.first_name && user.last_name 
                  ? `${user.first_name} ${user.last_name}` 
                  : user.first_name || user.name || 'User'}
              </Typography>
              <Box
                sx={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: 1,
                  color: '#6B7280',
                  cursor: 'pointer',
                  '&:hover': {
                    color: '#3B82F6',
                  },
                }}
                onClick={() => navigate('/')}
              >
                <HomeIcon fontSize="small" />
                <Typography variant="body2">Overview</Typography>
              </Box>
            </Box>
          </Box>
        </CardContent>
      </Card>

      {/* Primary Details */}
      <Card
        sx={{
          mb: 3,
          borderRadius: '12px',
          border: '1px solid',
          borderColor: 'rgba(229, 231, 235, 0.8)',
          backgroundColor: '#FFFFFF',
          boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
        }}
      >
        <CardContent sx={{ p: 3 }}>
          <Typography
            variant="h6"
            sx={{
              fontWeight: 700,
              color: '#1F2937',
              mb: 3,
            }}
          >
            Primary Details
          </Typography>

          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <Box sx={{ mb: 2 }}>
                <Typography variant="body2" sx={{ color: '#6B7280', mb: 1, fontWeight: 500 }}>
                  First Name
                </Typography>
                {editing.firstName ? (
                  <Box sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
                    <TextField
                      fullWidth
                      size="small"
                      value={formData.firstName}
                      onChange={(e) => setFormData({ ...formData, firstName: e.target.value })}
                      autoFocus
                    />
                    <Button
                      size="small"
                      variant="contained"
                      onClick={handleSaveName}
                    >
                      Save
                    </Button>
                    <Button
                      size="small"
                      onClick={() => {
                        setEditing({ ...editing, firstName: false });
                        const nameParts = (user.name || '').split(' ');
                        setFormData({ ...formData, firstName: nameParts[0] || '' });
                      }}
                    >
                      Cancel
                    </Button>
                  </Box>
                ) : (
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <Typography variant="body1" sx={{ color: '#1F2937' }}>
                      {formData.firstName || '-'}
                    </Typography>
                    <IconButton
                      size="small"
                      onClick={() => setEditing({ ...editing, firstName: true })}
                    >
                      <EditIcon fontSize="small" />
                    </IconButton>
                  </Box>
                )}
              </Box>

              <Box sx={{ mb: 2 }}>
                <Typography variant="body2" sx={{ color: '#6B7280', mb: 1, fontWeight: 500 }}>
                  Last Name
                </Typography>
                {editing.lastName ? (
                  <Box sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
                    <TextField
                      fullWidth
                      size="small"
                      value={formData.lastName}
                      onChange={(e) => setFormData({ ...formData, lastName: e.target.value })}
                      autoFocus
                    />
                    <Button
                      size="small"
                      variant="contained"
                      onClick={handleSaveName}
                    >
                      Save
                    </Button>
                    <Button
                      size="small"
                      onClick={() => {
                        setEditing({ ...editing, lastName: false });
                        const nameParts = (user.name || '').split(' ');
                        setFormData({ ...formData, lastName: nameParts.slice(1).join(' ') || '' });
                      }}
                    >
                      Cancel
                    </Button>
                  </Box>
                ) : (
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <Typography variant="body1" sx={{ color: '#1F2937' }}>
                      {formData.lastName || '-'}
                    </Typography>
                    <IconButton
                      size="small"
                      onClick={() => setEditing({ ...editing, lastName: true })}
                    >
                      <EditIcon fontSize="small" />
                    </IconButton>
                  </Box>
                )}
              </Box>

              <Box sx={{ mb: 2 }}>
                <Typography variant="body2" sx={{ color: '#6B7280', mb: 1, fontWeight: 500 }}>
                  Email Address
                </Typography>
                <Typography variant="body1" sx={{ color: '#1F2937' }}>
                  {formData.email || '-'}
                </Typography>
              </Box>

              <Box sx={{ mb: 2 }}>
                <Typography variant="body2" sx={{ color: '#6B7280', mb: 1, fontWeight: 500 }}>
                  Location
                </Typography>
                <Typography variant="body1" sx={{ color: '#1F2937' }}>
                  {formData.location || '-'}
                </Typography>
              </Box>
            </Grid>

            <Grid item xs={12} md={6}>
              <Box sx={{ mb: 2 }}>
                <Typography variant="body2" sx={{ color: '#6B7280', mb: 1, fontWeight: 500 }}>
                  Contact Number
                </Typography>
                <Typography variant="body1" sx={{ color: '#1F2937' }}>
                  {formData.phone || '-'}
                </Typography>
              </Box>

              <Box sx={{ mb: 2 }}>
                <Typography variant="body2" sx={{ color: '#6B7280', mb: 1, fontWeight: 500 }}>
                  Alternate Contact Number
                </Typography>
                <Typography variant="body1" sx={{ color: '#1F2937' }}>
                  {formData.alternatePhone || '-'}
                </Typography>
              </Box>

              <Box sx={{ mb: 2 }}>
                <Typography variant="body2" sx={{ color: '#6B7280', mb: 1, fontWeight: 500 }}>
                  WhatsApp Contact Number
                </Typography>
                <Typography variant="body1" sx={{ color: '#1F2937' }}>
                  {formData.whatsappPhone || '-'}
                </Typography>
              </Box>

              <Box sx={{ mb: 2 }}>
                <Typography variant="body2" sx={{ color: '#6B7280', mb: 1, fontWeight: 500 }}>
                  User Id
                </Typography>
                <Typography variant="body1" sx={{ color: '#1F2937' }}>
                  {user.id || '-'}
                </Typography>
              </Box>
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      {/* Organizational Details */}
      <Card
        sx={{
          mb: 3,
          borderRadius: '12px',
          border: '1px solid',
          borderColor: 'rgba(229, 231, 235, 0.8)',
          backgroundColor: '#FFFFFF',
          boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
        }}
      >
        <CardContent sx={{ p: 3 }}>
          <Typography
            variant="h6"
            sx={{
              fontWeight: 700,
              color: '#1F2937',
              mb: 3,
            }}
          >
            Organizational Details
          </Typography>

          <Grid container spacing={3} sx={{ mb: 3 }}>
            <Grid item xs={12} md={6}>
              <Box sx={{ mb: 2 }}>
                <Typography variant="body2" sx={{ color: '#6B7280', mb: 1, fontWeight: 500 }}>
                  Team
                </Typography>
                <Typography variant="body1" sx={{ color: '#1F2937' }}>
                  {user.team_name || '-'}
                </Typography>
              </Box>

              <Box sx={{ mb: 2 }}>
                <Typography variant="body2" sx={{ color: '#6B7280', mb: 1, fontWeight: 500 }}>
                  User Type
                </Typography>
                <Typography variant="body1" sx={{ color: '#1F2937' }}>
                  {user.user_type === 'product_user' ? 'Product User' : 'Observer'}
                </Typography>
              </Box>

              <Box sx={{ mb: 2 }}>
                <Typography variant="body2" sx={{ color: '#6B7280', mb: 1, fontWeight: 500 }}>
                  License Type
                </Typography>
                <Typography variant="body1" sx={{ color: '#1F2937' }}>
                  Free
                </Typography>
              </Box>
            </Grid>

            <Grid item xs={12} md={6}>
              <Box sx={{ mb: 2 }}>
                <Typography variant="body2" sx={{ color: '#6B7280', mb: 1, fontWeight: 500 }}>
                  Role
                </Typography>
                <Typography variant="body1" sx={{ color: '#1F2937' }}>
                  {user.role_name || '-'}
                </Typography>
              </Box>

              <Box sx={{ mb: 2 }}>
                <Typography variant="body2" sx={{ color: '#6B7280', mb: 1, fontWeight: 500 }}>
                  Manager
                </Typography>
                <Typography variant="body1" sx={{ color: '#1F2937' }}>
                  {user.manager_name || '-'}
                </Typography>
              </Box>

              <Box sx={{ mb: 2 }}>
                <Typography variant="body2" sx={{ color: '#6B7280', mb: 1, fontWeight: 500 }}>
                  Join Date
                </Typography>
                <Typography variant="body1" sx={{ color: '#1F2937' }}>
                  {formatDate(user.created_at)}
                </Typography>
              </Box>
            </Grid>
          </Grid>

          {/* Action Buttons */}
          <Box sx={{ display: 'flex', gap: 2, mt: 3 }}>
            <Button
              variant="contained"
              onClick={() => setPasswordDialogOpen(true)}
              sx={{
                borderRadius: '8px',
                textTransform: 'none',
                fontWeight: 600,
                px: 3,
                py: 1.5,
              }}
            >
              Change Password
            </Button>
            <Button
              variant="outlined"
              onClick={() => setTimezoneDialogOpen(true)}
              sx={{
                borderRadius: '8px',
                textTransform: 'none',
                fontWeight: 600,
                px: 3,
                py: 1.5,
                borderColor: 'rgba(229, 231, 235, 1)',
                color: '#1F2937',
              }}
            >
              Change Timezone
            </Button>
          </Box>
        </CardContent>
      </Card>

      {/* Change Password Dialog */}
      <Dialog open={passwordDialogOpen} onClose={() => setPasswordDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle sx={{ pb: 1 }}>Change Password</DialogTitle>
        <DialogContent sx={{ pt: 2 }}>
          <TextField
            fullWidth
            label="Current Password"
            type={showPassword.current ? 'text' : 'password'}
            value={passwordData.currentPassword}
            onChange={(e) => setPasswordData({ ...passwordData, currentPassword: e.target.value })}
            margin="normal"
            InputProps={{
              endAdornment: (
                <InputAdornment position="end">
                  <IconButton
                    onClick={() => setShowPassword({ ...showPassword, current: !showPassword.current })}
                    edge="end"
                  >
                    {showPassword.current ? <VisibilityOff /> : <Visibility />}
                  </IconButton>
                </InputAdornment>
              ),
            }}
          />
          <TextField
            fullWidth
            label="New Password"
            type={showPassword.new ? 'text' : 'password'}
            value={passwordData.newPassword}
            onChange={(e) => setPasswordData({ ...passwordData, newPassword: e.target.value })}
            margin="normal"
            InputProps={{
              endAdornment: (
                <InputAdornment position="end">
                  <IconButton
                    onClick={() => setShowPassword({ ...showPassword, new: !showPassword.new })}
                    edge="end"
                  >
                    {showPassword.new ? <VisibilityOff /> : <Visibility />}
                  </IconButton>
                </InputAdornment>
              ),
            }}
          />
          <TextField
            fullWidth
            label="Confirm New Password"
            type={showPassword.confirm ? 'text' : 'password'}
            value={passwordData.confirmPassword}
            onChange={(e) => setPasswordData({ ...passwordData, confirmPassword: e.target.value })}
            margin="normal"
            InputProps={{
              endAdornment: (
                <InputAdornment position="end">
                  <IconButton
                    onClick={() => setShowPassword({ ...showPassword, confirm: !showPassword.confirm })}
                    edge="end"
                  >
                    {showPassword.confirm ? <VisibilityOff /> : <Visibility />}
                  </IconButton>
                </InputAdornment>
              ),
            }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setPasswordDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleChangePassword} variant="contained">
            Change Password
          </Button>
        </DialogActions>
      </Dialog>

      {/* Change Timezone Dialog */}
      <Dialog open={timezoneDialogOpen} onClose={() => setTimezoneDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle sx={{ pb: 1 }}>Change Timezone</DialogTitle>
        <DialogContent sx={{ pt: 2 }}>
          <FormControl fullWidth margin="normal">
            <InputLabel>Timezone</InputLabel>
            <Select
              value={timezone}
              label="Timezone"
              onChange={(e) => setTimezone(e.target.value)}
            >
              <MenuItem value="UTC">UTC</MenuItem>
              <MenuItem value="America/New_York">America/New_York (EST/EDT)</MenuItem>
              <MenuItem value="America/Chicago">America/Chicago (CST/CDT)</MenuItem>
              <MenuItem value="America/Denver">America/Denver (MST/MDT)</MenuItem>
              <MenuItem value="America/Los_Angeles">America/Los_Angeles (PST/PDT)</MenuItem>
              <MenuItem value="Europe/London">Europe/London (GMT/BST)</MenuItem>
              <MenuItem value="Europe/Paris">Europe/Paris (CET/CEST)</MenuItem>
              <MenuItem value="Asia/Dubai">Asia/Dubai (GST)</MenuItem>
              <MenuItem value="Asia/Kolkata">Asia/Kolkata (IST)</MenuItem>
              <MenuItem value="Asia/Singapore">Asia/Singapore (SGT)</MenuItem>
              <MenuItem value="Asia/Tokyo">Asia/Tokyo (JST)</MenuItem>
              <MenuItem value="Australia/Sydney">Australia/Sydney (AEDT/AEST)</MenuItem>
            </Select>
          </FormControl>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setTimezoneDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleChangeTimezone} variant="contained">
            Save Timezone
          </Button>
        </DialogActions>
      </Dialog>

      <Snackbar
        open={snackbar.open}
        autoHideDuration={6000}
        onClose={() => setSnackbar({ ...snackbar, open: false })}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      >
        <Alert onClose={() => setSnackbar({ ...snackbar, open: false })} severity={snackbar.severity} sx={{ width: '100%' }}>
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Box>
  );
}

export default Profile;

