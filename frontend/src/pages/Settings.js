import React, { useState, useEffect } from 'react';
import {
  Box,
  Tabs,
  Tab,
  Typography,
  Paper,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Chip,
  Alert,
  Snackbar,
  Checkbox,
  FormControlLabel,
  FormGroup,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  InputAdornment,
  Collapse,
} from '@mui/material';
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  ExpandMore as ExpandMoreIcon,
  Upload as UploadIcon,
  Download as DownloadIcon,
  Visibility,
  VisibilityOff,
  KeyboardArrowDown,
  KeyboardArrowUp,
  PersonAdd as PersonAddIcon,
} from '@mui/icons-material';
import api from '../api/client';
import { typography } from '../theme/typography';

function TabPanel({ children, value, index }) {
  return (
    <div role="tabpanel" hidden={value !== index}>
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  );
}

function Settings() {
  const [tabValue, setTabValue] = useState(0);
  const [users, setUsers] = useState([]);
  const [roles, setRoles] = useState([]);
  const [permissionGroups, setPermissionGroups] = useState([]);
  const [teams, setTeams] = useState([]);
  const [useCases, setUseCases] = useState([]);
  const [vendors, setVendors] = useState([]);
  const [loading, setLoading] = useState(false);
  const [snackbar, setSnackbar] = useState({ open: false, message: '', severity: 'success' });

  // User Manager State
  const [userDialogOpen, setUserDialogOpen] = useState(false);
  const [userForm, setUserForm] = useState({
    email: '',
    name: '',
    phone: '',
    role_id: '',
    manager_id: '',
    auditor_id: '',
    team_id: '',
    user_type: 'standard',
    location: '',
    password: '',
  });
  const [showPassword, setShowPassword] = useState(false);
  const [editingUser, setEditingUser] = useState(null);

  // Role Manager State
  const [roleDialogOpen, setRoleDialogOpen] = useState(false);
  const [roleForm, setRoleForm] = useState({
    name: '',
    description: '',
    code_names: [],
    allowed_team_ids: [],
  });
  const [editingRole, setEditingRole] = useState(null);

  // Team Manager State
  const [teamDialogOpen, setTeamDialogOpen] = useState(false);
  const [subteamDialogOpen, setSubteamDialogOpen] = useState(false);
  const [teamForm, setTeamForm] = useState({
    name: '',
    description: '',
    vendor_id: '',
    group_id: null,
  });
  const [editingTeam, setEditingTeam] = useState(null);
  const [parentTeamForSubteam, setParentTeamForSubteam] = useState(null);
  const [expandedTeams, setExpandedTeams] = useState({}); // Track which teams have expanded sub-teams
  const [assignUserDialogOpen, setAssignUserDialogOpen] = useState(false);
  const [selectedTeamForUser, setSelectedTeamForUser] = useState(null);
  const [selectedSubteamForUser, setSelectedSubteamForUser] = useState(null);

  useEffect(() => {
    if (tabValue === 0) {
      fetchUsers();
    } else if (tabValue === 1) {
      fetchRoles();
      fetchPermissions();
    } else if (tabValue === 2) {
      fetchUsers();
      fetchTeams();
      fetchUseCases();
      fetchVendors();
    }
  }, [tabValue]);

  const fetchUsers = async () => {
    try {
      setLoading(true);
      const response = await api.get('/users');
      setUsers(response.data.users || []);
    } catch (error) {
      showSnackbar('Failed to fetch users', 'error');
    } finally {
      setLoading(false);
    }
  };

  const fetchRoles = async () => {
    try {
      setLoading(true);
      const response = await api.get('/roles');
      setRoles(response.data.roles || []);
    } catch (error) {
      showSnackbar('Failed to fetch roles', 'error');
    } finally {
      setLoading(false);
    }
  };

  const fetchPermissions = async () => {
    try {
      setLoading(true);
      const response = await api.get('/permissions');
      setPermissionGroups(response.data.permission_groups || []);
    } catch (error) {
      showSnackbar('Failed to fetch permissions', 'error');
    } finally {
      setLoading(false);
    }
  };

  const fetchTeams = async () => {
    try {
      const response = await api.get('/teams', { params: { include_subteams: true } });
      setTeams(response.data.teams || []);
    } catch (error) {
      showSnackbar('Failed to fetch teams', 'error');
    }
  };

  const fetchUseCases = async () => {
    try {
      const response = await api.get('/use-cases');
      setUseCases(response.data.use_cases || []);
    } catch (error) {
      showSnackbar('Failed to fetch use cases', 'error');
    }
  };

  const fetchVendors = async () => {
    try {
      const response = await api.get('/vendors');
      setVendors(response.data.vendors || []);
    } catch (error) {
      console.error('Error fetching vendors:', error);
    }
  };

  const handleAssignTeam = async (userId, teamId) => {
    try {
      const user = users.find(u => u.id === userId);
      if (!user) return;

      await api.put(`/users/${userId}`, {
        name: user.name,
        phone: user.phone,
        role_id: user.role_id,
        manager_id: user.manager_id,
        auditor_id: user.auditor_id,
        team_id: teamId ? parseInt(teamId) : null,
        user_type: user.user_type,
        location: user.location,
      });
      showSnackbar('Team assignment updated successfully');
      fetchUsers();
    } catch (error) {
      showSnackbar(error.response?.data?.error || 'Failed to update team assignment', 'error');
    }
  };

  const handleSaveTeam = async () => {
    try {
      const payload = {
        name: teamForm.name,
        description: teamForm.description,
        vendor_id: teamForm.vendor_id && teamForm.vendor_id !== '' ? parseInt(teamForm.vendor_id) : null,
        group_id: teamForm.group_id,
        members: [],
        subteams: [],
      };

      if (editingTeam) {
        await api.put(`/teams/${editingTeam.id}`, payload);
        showSnackbar('Team updated successfully');
      } else {
        await api.post('/teams', payload);
        showSnackbar('Team created successfully');
      }
      setTeamDialogOpen(false);
      fetchTeams();
      fetchUsers();
    } catch (error) {
      showSnackbar(error.response?.data?.error || 'Failed to save team', 'error');
    }
  };

  const handleSaveSubteam = async () => {
    try {
      const payload = {
        name: teamForm.name,
        description: teamForm.description,
        vendor_id: teamForm.vendor_id && teamForm.vendor_id !== '' ? parseInt(teamForm.vendor_id) : null,
        group_id: teamForm.group_id,
        members: [],
        subteams: [],
      };

      await api.post('/teams', payload);
      showSnackbar('Sub-team created successfully');
      setSubteamDialogOpen(false);
      fetchTeams();
      fetchUsers();
    } catch (error) {
      showSnackbar(error.response?.data?.error || 'Failed to create sub-team', 'error');
    }
  };

  const handleDeleteTeam = async (teamId) => {
    const team = teams.find(t => t.id === teamId);
    const memberCount = users.filter(u => u.team_id === teamId).length;
    
    if (memberCount > 0) {
      const transferTo = window.prompt(
        `This team has ${memberCount} members. Enter a team ID to transfer them to (or cancel to abort):`
      );
      if (!transferTo) return;

      try {
        await api.delete(`/teams/${teamId}`, {
          data: { transfer_to_team_id: parseInt(transferTo) },
        });
        showSnackbar('Team deleted successfully');
        fetchTeams();
        fetchUsers();
      } catch (error) {
        showSnackbar(error.response?.data?.error || 'Failed to delete team', 'error');
      }
    } else {
      if (window.confirm('Are you sure you want to delete this team?')) {
        try {
          await api.delete(`/teams/${teamId}`);
          showSnackbar('Team deleted successfully');
          fetchTeams();
        } catch (error) {
          showSnackbar(error.response?.data?.error || 'Failed to delete team', 'error');
        }
      }
    }
  };

  const showSnackbar = (message, severity = 'success') => {
    setSnackbar({ open: true, message, severity });
  };

  const handleTabChange = (event, newValue) => {
    setTabValue(newValue);
  };

  // User Manager Functions
  const handleCreateUser = () => {
    setEditingUser(null);
    setUserForm({
      email: '',
      name: '',
      phone: '',
      role_id: '',
      manager_id: '',
      auditor_id: '',
      team_id: '',
      user_type: 'standard',
      location: '',
      password: '',
    });
    setShowPassword(false);
    setUserDialogOpen(true);
  };

  const handleEditUser = (user) => {
    setEditingUser(user);
    setUserForm({
      email: user.email,
      name: user.name,
      phone: user.phone || '',
      role_id: user.role_id || '',
      manager_id: user.manager_id || '',
      auditor_id: user.auditor_id || '',
      team_id: user.team_id || '',
      user_type: user.user_type || 'standard',
      location: user.location || '',
      password: '',
    });
    setShowPassword(false);
    setUserDialogOpen(true);
  };

  const handleSaveUser = async () => {
    try {
      const payload = {
        ...userForm,
        role_id: userForm.role_id ? parseInt(userForm.role_id) : null,
        manager_id: userForm.manager_id ? parseInt(userForm.manager_id) : null,
        auditor_id: userForm.auditor_id ? parseInt(userForm.auditor_id) : null,
        team_id: userForm.team_id ? parseInt(userForm.team_id) : null,
        phone: userForm.phone || null,
        location: userForm.location || null,
      };

      // Only include password if it's provided (for editing) or if creating new user
      if (editingUser) {
        // For editing, only include password if it's not empty
        if (userForm.password) {
          payload.password = userForm.password;
        } else {
          delete payload.password;
        }
        await api.put(`/users/${editingUser.id}`, payload);
        showSnackbar('User updated successfully');
      } else {
        // For new users, password is optional (will be auto-generated if not provided)
        if (userForm.password) {
          payload.password = userForm.password;
        }
        const response = await api.post('/users', payload);
        showSnackbar(`User created successfully. Password: ${response.data.password}`);
      }
      setUserDialogOpen(false);
      fetchUsers();
    } catch (error) {
      showSnackbar(error.response?.data?.error || 'Failed to save user', 'error');
    }
  };

  const handleDeleteUser = async (userId) => {
    if (window.confirm('Are you sure you want to delete this user?')) {
      try {
        await api.delete(`/users/${userId}`);
        showSnackbar('User deleted successfully');
        fetchUsers();
      } catch (error) {
        showSnackbar(error.response?.data?.error || 'Failed to delete user', 'error');
      }
    }
  };

  const handleBulkUpload = async (event) => {
    const file = event.target.files[0];
    if (!file) return;

    // Simple CSV parsing (in production, use a proper CSV library)
    const text = await file.text();
    const lines = text.split('\n');
    const headers = lines[0].split(',').map(h => h.trim());
    const users = [];

    for (let i = 1; i < lines.length; i++) {
      if (!lines[i].trim()) continue;
      const values = lines[i].split(',').map(v => v.trim());
      const user = {};
      headers.forEach((header, idx) => {
        user[header.toLowerCase().replace(/\s+/g, '_')] = values[idx] || '';
      });
      users.push({
        email: user.email,
        name: user.name || user.email.split('@')[0],
        phone: user.phone || null,
        role_id: user.role_id ? parseInt(user.role_id) : null,
        manager_id: user.manager_id ? parseInt(user.manager_id) : null,
        auditor_id: user.auditor_id ? parseInt(user.auditor_id) : null,
        team_id: user.team_id ? parseInt(user.team_id) : null,
        user_type: user.user_type || 'product_user',
        location: user.location || null,
      });
    }

    try {
      await api.post('/users/bulk', { users });
      showSnackbar(`Successfully uploaded ${users.length} users`);
      fetchUsers();
    } catch (error) {
      showSnackbar(error.response?.data?.error || 'Failed to upload users', 'error');
    }
  };

  // Role Manager Functions
  const handleCreateRole = () => {
    setEditingRole(null);
    setRoleForm({
      name: '',
      description: '',
      code_names: [],
      allowed_team_ids: [],
    });
    setRoleDialogOpen(true);
  };

  const handleEditRole = (role) => {
    setEditingRole(role);
    setRoleForm({
      name: role.name,
      description: role.description || '',
      code_names: role.code_names || [],
      allowed_team_ids: role.allowed_team_ids || [],
    });
    setRoleDialogOpen(true);
  };

  const handleSaveRole = async () => {
    try {
      const payload = {
        ...roleForm,
        description: roleForm.description || null,
      };

      if (editingRole) {
        await api.put(`/roles/${editingRole.id}`, payload);
        showSnackbar('Role updated successfully');
      } else {
        await api.post('/roles', payload);
        showSnackbar('Role created successfully');
      }
      setRoleDialogOpen(false);
      fetchRoles();
    } catch (error) {
      showSnackbar(error.response?.data?.error || 'Failed to save role', 'error');
    }
  };

  const handleDeleteRole = async (roleId) => {
    if (window.confirm('Are you sure you want to delete this role?')) {
      try {
        await api.delete(`/roles/${roleId}`);
        showSnackbar('Role deleted successfully');
        fetchRoles();
      } catch (error) {
        showSnackbar(error.response?.data?.error || 'Failed to delete role', 'error');
      }
    }
  };

  const togglePermission = (codeName) => {
    setRoleForm(prev => ({
      ...prev,
      code_names: prev.code_names.includes(codeName)
        ? prev.code_names.filter(cn => cn !== codeName)
        : [...prev.code_names, codeName],
    }));
  };

  return (
    <Box sx={{ width: '100%' }}>
      <Typography variant="h4" sx={{ ...typography.pageTitle, mb: 3 }}>
        Settings
      </Typography>

      <Paper sx={{ mb: 3 }}>
        <Tabs value={tabValue} onChange={handleTabChange}>
          <Tab label="User Manager" />
          <Tab label="Role Manager" />
          <Tab label="Team Manager" />
        </Tabs>
      </Paper>

      {/* User Manager Tab */}
      <TabPanel value={tabValue} index={0}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
          <Typography variant="h6" sx={{ ...typography.sectionHeader, fontSize: '1.125rem' }}>Users</Typography>
          <Box>
            <input
              accept=".csv,.xlsx"
              style={{ display: 'none' }}
              id="bulk-upload"
              type="file"
              onChange={handleBulkUpload}
            />
            <label htmlFor="bulk-upload">
              <Button
                component="span"
                variant="outlined"
                startIcon={<UploadIcon />}
                sx={{ mr: 2 }}
              >
                Bulk Upload
              </Button>
            </label>
            <Button
              variant="contained"
              startIcon={<AddIcon />}
              onClick={handleCreateUser}
            >
              Add User
            </Button>
          </Box>
        </Box>

        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>Email</TableCell>
                <TableCell>Role</TableCell>
                <TableCell>Team</TableCell>
                <TableCell>Type</TableCell>
                <TableCell>Status</TableCell>
                <TableCell align="right">Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {users.map((user) => (
                <TableRow key={user.id}>
                  <TableCell>{user.name}</TableCell>
                  <TableCell>{user.email}</TableCell>
                  <TableCell>{user.role_name || '-'}</TableCell>
                  <TableCell>
                    {(() => {
                      const userTeam = teams.find(t => t.id === user.team_id);
                      if (!userTeam) return '-';
                      // If it's a subteam (has group_id), show parent team name
                      if (userTeam.group_id) {
                        const parentTeam = teams.find(t => t.id === userTeam.group_id);
                        return parentTeam ? parentTeam.name : '-';
                      }
                      return userTeam.name;
                    })()}
                  </TableCell>
                  <TableCell>
                    {(() => {
                      const userTeam = teams.find(t => t.id === user.team_id);
                      if (!userTeam || !userTeam.group_id) return '-';
                      return userTeam.name;
                    })()}
                  </TableCell>
                  <TableCell>
                    <Chip
                      label={user.user_type === 'product_user' ? 'Product User' : 'Observer'}
                      color={user.user_type === 'product_user' ? 'primary' : 'default'}
                      size="small"
                    />
                  </TableCell>
                  <TableCell>
                    <Chip
                      label={user.is_active ? 'Active' : 'Inactive'}
                      color={user.is_active ? 'success' : 'default'}
                      size="small"
                    />
                  </TableCell>
                  <TableCell align="right">
                    <IconButton size="small" onClick={() => handleEditUser(user)}>
                      <EditIcon />
                    </IconButton>
                    <IconButton size="small" onClick={() => handleDeleteUser(user.id)}>
                      <DeleteIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </TabPanel>

      {/* Role Manager Tab */}
      <TabPanel value={tabValue} index={1}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
          <Typography variant="h6" sx={{ ...typography.sectionHeader, fontSize: '1.125rem' }}>Roles</Typography>
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={handleCreateRole}
          >
            Create Role
          </Button>
        </Box>

        <TableContainer component={Paper} sx={{ mb: 3 }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>Description</TableCell>
                <TableCell>Permissions</TableCell>
                <TableCell>Users</TableCell>
                <TableCell>Status</TableCell>
                <TableCell align="right">Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {roles.map((role) => (
                <TableRow key={role.id}>
                  <TableCell>
                    {role.name}
                    {role.is_default && (
                      <Chip label="Default" size="small" sx={{ ml: 1 }} />
                    )}
                  </TableCell>
                  <TableCell>{role.description || '-'}</TableCell>
                  <TableCell>{role.code_names?.length || 0} permissions</TableCell>
                  <TableCell>{role.user_count || 0}</TableCell>
                  <TableCell>
                    <Chip
                      label={role.can_be_edited ? 'Editable' : 'Locked'}
                      color={role.can_be_edited ? 'default' : 'warning'}
                      size="small"
                    />
                  </TableCell>
                  <TableCell align="right">
                    <IconButton
                      size="small"
                      onClick={() => handleEditRole(role)}
                      disabled={!role.can_be_edited}
                    >
                      <EditIcon />
                    </IconButton>
                    <IconButton
                      size="small"
                      onClick={() => handleDeleteRole(role.id)}
                      disabled={!role.can_be_edited || role.is_default}
                    >
                      <DeleteIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </TabPanel>

      {/* Team Manager Tab */}
      <TabPanel value={tabValue} index={2}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
          <Typography variant="h6" sx={{ ...typography.sectionHeader, fontSize: '1.125rem' }}>Team Management</Typography>
          <Button
            variant="contained"
            startIcon={<AddIcon />}
            onClick={() => {
              setTeamDialogOpen(true);
              setEditingTeam(null);
              setTeamForm({
                name: '',
                description: '',
                vendor_id: '',
                group_id: null,
              });
            }}
          >
            Create Team
          </Button>
        </Box>

        {/* Teams Tree View */}
        <TableContainer component={Paper} sx={{ mb: 3 }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Team Name</TableCell>
                <TableCell>Vendor</TableCell>
                <TableCell>Members</TableCell>
                <TableCell>Subteams</TableCell>
                <TableCell align="right">Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {teams.filter(t => !t.group_id).map((team) => {
                const subteams = teams.filter(t => t.group_id === team.id);
                const teamMembers = users.filter(u => u.team_id === team.id);
                const isExpanded = expandedTeams[team.id] || false;
                
                return (
                  <React.Fragment key={team.id}>
                    <TableRow>
                      <TableCell>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          {subteams.length > 0 && (
                            <IconButton
                              size="small"
                              onClick={() => {
                                setExpandedTeams({ ...expandedTeams, [team.id]: !isExpanded });
                              }}
                            >
                              {isExpanded ? <KeyboardArrowUp /> : <KeyboardArrowDown />}
                            </IconButton>
                          )}
                          <Box>
                            <Typography variant="body1" fontWeight={600}>
                              {team.name}
                            </Typography>
                            {team.description && (
                              <Typography variant="caption" color="text.secondary">
                                {team.description}
                              </Typography>
                            )}
                          </Box>
                        </Box>
                      </TableCell>
                      <TableCell>
                        {vendors.find(v => v.id === team.vendor_id)?.name || '-'}
                      </TableCell>
                      <TableCell>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                          <Typography variant="body2">
                            {teamMembers.length} members
                          </Typography>
                          <IconButton
                            size="small"
                            onClick={() => {
                              setSelectedTeamForUser(team);
                              setSelectedSubteamForUser(null);
                              setAssignUserDialogOpen(true);
                            }}
                            title="Add User to Team"
                          >
                            <PersonAddIcon fontSize="small" />
                          </IconButton>
                        </Box>
                      </TableCell>
                      <TableCell>
                        {subteams.length} subteams
                      </TableCell>
                      <TableCell align="right">
                        <IconButton
                          size="small"
                          onClick={() => {
                            setEditingTeam(team);
                            setTeamForm({
                              name: team.name,
                              description: team.description || '',
                              vendor_id: team.vendor_id,
                              group_id: team.group_id,
                            });
                            setTeamDialogOpen(true);
                          }}
                        >
                          <EditIcon />
                        </IconButton>
                        <IconButton
                          size="small"
                          onClick={() => {
                            setParentTeamForSubteam(team);
                            setSubteamDialogOpen(true);
                            setTeamForm({
                              name: '',
                              description: '',
                              vendor_id: team.vendor_id,
                              group_id: team.id,
                            });
                          }}
                          title="Create Sub-team"
                        >
                          <AddIcon />
                        </IconButton>
                        <IconButton
                          size="small"
                          onClick={() => handleDeleteTeam(team.id)}
                        >
                          <DeleteIcon />
                        </IconButton>
                      </TableCell>
                    </TableRow>
                    {/* Subteams - Collapsible */}
                    {subteams.length > 0 && (
                      <TableRow>
                        <TableCell colSpan={5} sx={{ py: 0, border: 0 }}>
                          <Collapse in={isExpanded} timeout="auto" unmountOnExit>
                            <Box sx={{ bgcolor: 'grey.50', pl: 4, pr: 2, py: 1 }}>
                              {subteams.map((subteam) => {
                                const subteamMembers = users.filter(u => u.team_id === subteam.id);
                                return (
                                  <Box key={subteam.id} sx={{ display: 'flex', alignItems: 'center', py: 1, borderBottom: '1px solid', borderColor: 'divider' }}>
                                    <Box sx={{ flex: 1 }}>
                                      <Typography variant="body2" fontWeight={500}>
                                        └ {subteam.name}
                                      </Typography>
                                      {subteam.description && (
                                        <Typography variant="caption" color="text.secondary">
                                          {subteam.description}
                                        </Typography>
                                      )}
                                    </Box>
                                    <Box sx={{ minWidth: 150, textAlign: 'center' }}>
                                      {vendors.find(v => v.id === subteam.vendor_id)?.name || '-'}
                                    </Box>
                                    <Box sx={{ minWidth: 150, display: 'flex', alignItems: 'center', gap: 1, justifyContent: 'center' }}>
                                      <Typography variant="body2">
                                        {subteamMembers.length} members
                                      </Typography>
                                      <IconButton
                                        size="small"
                                        onClick={() => {
                                          setSelectedTeamForUser(team);
                                          setSelectedSubteamForUser(subteam);
                                          setAssignUserDialogOpen(true);
                                        }}
                                        title="Add User to Sub-team"
                                      >
                                        <PersonAddIcon fontSize="small" />
                                      </IconButton>
                                    </Box>
                                    <Box sx={{ minWidth: 100 }}></Box>
                                    <Box sx={{ minWidth: 150, display: 'flex', justifyContent: 'flex-end', gap: 0.5 }}>
                                      <IconButton
                                        size="small"
                                        onClick={() => {
                                          setEditingTeam(subteam);
                                          setTeamForm({
                                            name: subteam.name,
                                            description: subteam.description || '',
                                            vendor_id: subteam.vendor_id,
                                            group_id: subteam.group_id,
                                          });
                                          setTeamDialogOpen(true);
                                        }}
                                      >
                                        <EditIcon fontSize="small" />
                                      </IconButton>
                                      <IconButton
                                        size="small"
                                        onClick={() => handleDeleteTeam(subteam.id)}
                                      >
                                        <DeleteIcon fontSize="small" />
                                      </IconButton>
                                    </Box>
                                  </Box>
                                );
                              })}
                            </Box>
                          </Collapse>
                        </TableCell>
                      </TableRow>
                    )}
                  </React.Fragment>
                );
              })}
            </TableBody>
          </Table>
        </TableContainer>

      </TabPanel>

      {/* User Dialog */}
      <Dialog open={userDialogOpen} onClose={() => setUserDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle sx={{ pb: 1 }}>{editingUser ? 'Edit User' : 'Create User'}</DialogTitle>
        <DialogContent sx={{ pt: 2 }}>
          <TextField
            fullWidth
            label="Email"
            type="email"
            value={userForm.email}
            onChange={(e) => setUserForm({ ...userForm, email: e.target.value })}
            margin="normal"
            required
            disabled={!!editingUser}
          />
          <TextField
            fullWidth
            label="Name"
            value={userForm.name}
            onChange={(e) => setUserForm({ ...userForm, name: e.target.value })}
            margin="normal"
            required
          />
          <TextField
            fullWidth
            label="Phone"
            value={userForm.phone}
            onChange={(e) => setUserForm({ ...userForm, phone: e.target.value })}
            margin="normal"
          />
          <FormControl fullWidth margin="normal" sx={{ mt: 2 }}>
            <InputLabel id="user-type-label" shrink>User Type</InputLabel>
            <Select
              labelId="user-type-label"
              value={userForm.user_type}
              label="User Type"
              onChange={(e) => setUserForm({ ...userForm, user_type: e.target.value })}
            >
              <MenuItem value="product_user">Product User</MenuItem>
              <MenuItem value="observer">Observer</MenuItem>
            </Select>
          </FormControl>
          {userForm.user_type === 'observer' && (
            <Alert severity="info" sx={{ mt: 1, mb: 1 }}>
              Observers have no functional access to the product, regardless of role permissions.
            </Alert>
          )}
          {userForm.user_type === 'product_user' && (
            <Alert severity="info" sx={{ mt: 1, mb: 1 }}>
              Product users can access features based on their role permissions.
            </Alert>
          )}
          <FormControl fullWidth margin="normal" sx={{ mt: 2 }}>
            <InputLabel id="role-label" shrink={!!userForm.role_id}>Role</InputLabel>
            <Select
              labelId="role-label"
              value={userForm.role_id || ''}
              label="Role"
              onChange={(e) => setUserForm({ ...userForm, role_id: e.target.value })}
              disabled={userForm.user_type === 'observer'}
            >
              <MenuItem value="">None</MenuItem>
              {roles.map((role) => (
                <MenuItem key={role.id} value={role.id}>
                  {role.name}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          {userForm.user_type === 'observer' && (
            <Alert severity="warning" sx={{ mt: 1, mb: 1 }}>
              Role assignment is disabled for observers as they have no functional access.
            </Alert>
          )}
          <TextField
            fullWidth
            label="Location"
            value={userForm.location}
            onChange={(e) => setUserForm({ ...userForm, location: e.target.value })}
            margin="normal"
            sx={{ mt: 2 }}
          />
          {editingUser && (
            <TextField
              fullWidth
              label="New Password"
              type={showPassword ? 'text' : 'password'}
              value={userForm.password}
              onChange={(e) => setUserForm({ ...userForm, password: e.target.value })}
              margin="normal"
              sx={{ mt: 2 }}
              helperText="Leave empty to keep current password"
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton
                      onClick={() => setShowPassword(!showPassword)}
                      edge="end"
                      size="small"
                    >
                      {showPassword ? <VisibilityOff fontSize="small" /> : <Visibility fontSize="small" />}
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setUserDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleSaveUser} variant="contained">
            {editingUser ? 'Update' : 'Create'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Role Dialog */}
      <Dialog open={roleDialogOpen} onClose={() => setRoleDialogOpen(false)} maxWidth="md" fullWidth>
        <DialogTitle sx={{ pb: 1 }}>{editingRole ? 'Edit Role' : 'Create Role'}</DialogTitle>
        <DialogContent sx={{ pt: 2 }}>
          <TextField
            fullWidth
            label="Role Name"
            value={roleForm.name}
            onChange={(e) => setRoleForm({ ...roleForm, name: e.target.value })}
            margin="normal"
            required
          />
          <TextField
            fullWidth
            label="Description"
            value={roleForm.description}
            onChange={(e) => setRoleForm({ ...roleForm, description: e.target.value })}
            margin="normal"
            multiline
            rows={3}
          />
          <Typography variant="subtitle2" sx={{ mt: 2, mb: 1 }}>
            Permissions
          </Typography>
          {permissionGroups.map((group) => (
            <Accordion key={group.id}>
              <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                <Typography>{group.name}</Typography>
              </AccordionSummary>
              <AccordionDetails>
                <FormGroup>
                  {group.permissions?.map((permission) => (
                    <FormControlLabel
                      key={permission.id}
                      control={
                        <Checkbox
                          checked={roleForm.code_names.includes(permission.code_name)}
                          onChange={() => togglePermission(permission.code_name)}
                        />
                      }
                      label={
                        <Box>
                          <Typography variant="body2" fontWeight={600}>
                            {permission.name}
                          </Typography>
                          {permission.description && (
                            <Typography variant="caption" color="text.secondary">
                              {permission.description}
                            </Typography>
                          )}
                        </Box>
                      }
                    />
                  ))}
                </FormGroup>
              </AccordionDetails>
            </Accordion>
          ))}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setRoleDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleSaveRole} variant="contained">
            {editingRole ? 'Update' : 'Create'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Team Dialog */}
      <Dialog open={teamDialogOpen} onClose={() => setTeamDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle sx={{ pb: 1 }}>{editingTeam ? 'Edit Team' : 'Create Team'}</DialogTitle>
        <DialogContent sx={{ pt: 2 }}>
          <TextField
            fullWidth
            label="Team Name"
            value={teamForm.name}
            onChange={(e) => setTeamForm({ ...teamForm, name: e.target.value })}
            margin="normal"
            required
          />
          <TextField
            fullWidth
            label="Description"
            value={teamForm.description}
            onChange={(e) => setTeamForm({ ...teamForm, description: e.target.value })}
            margin="normal"
            multiline
            rows={2}
          />
          <FormControl fullWidth margin="normal" sx={{ mt: 2 }}>
            <InputLabel id="vendor-label">Vendor</InputLabel>
            <Select
              labelId="vendor-label"
              value={teamForm.vendor_id}
              label="Vendor"
              onChange={(e) => setTeamForm({ ...teamForm, vendor_id: e.target.value })}
            >
              <MenuItem value="">None</MenuItem>
              {vendors.map((vendor) => (
                <MenuItem key={vendor.id} value={vendor.id}>
                  {vendor.name}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
          {!editingTeam && (
            <FormControl fullWidth margin="normal" sx={{ mt: 2 }}>
              <InputLabel id="parent-team-label">Parent Team (Optional)</InputLabel>
              <Select
                labelId="parent-team-label"
                value={teamForm.group_id || ''}
                label="Parent Team (Optional)"
                onChange={(e) => setTeamForm({ ...teamForm, group_id: e.target.value ? parseInt(e.target.value) : null })}
              >
                <MenuItem value="">None (Top-level Team)</MenuItem>
                {teams.filter(t => !t.group_id).map((team) => (
                  <MenuItem key={team.id} value={team.id}>
                    {team.name}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setTeamDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleSaveTeam} variant="contained">
            {editingTeam ? 'Update' : 'Create'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Sub-team Dialog */}
      <Dialog open={subteamDialogOpen} onClose={() => setSubteamDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle sx={{ pb: 1 }}>Create Sub-team</DialogTitle>
        <DialogContent sx={{ pt: 2 }}>
          {parentTeamForSubteam && (
            <Alert severity="info" sx={{ mb: 2 }}>
              Creating sub-team under: <strong>{parentTeamForSubteam.name}</strong>
            </Alert>
          )}
          <TextField
            fullWidth
            label="Sub-team Name"
            value={teamForm.name}
            onChange={(e) => setTeamForm({ ...teamForm, name: e.target.value })}
            margin="normal"
            required
          />
          <TextField
            fullWidth
            label="Description"
            value={teamForm.description}
            onChange={(e) => setTeamForm({ ...teamForm, description: e.target.value })}
            margin="normal"
            multiline
            rows={2}
          />
          <FormControl fullWidth margin="normal" sx={{ mt: 2 }}>
            <InputLabel id="subteam-vendor-label">Vendor</InputLabel>
            <Select
              labelId="subteam-vendor-label"
              value={teamForm.vendor_id}
              label="Vendor"
              onChange={(e) => setTeamForm({ ...teamForm, vendor_id: e.target.value })}
            >
              <MenuItem value="">None</MenuItem>
              {vendors.map((vendor) => (
                <MenuItem key={vendor.id} value={vendor.id}>
                  {vendor.name}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setSubteamDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleSaveSubteam} variant="contained">
            Create Sub-team
          </Button>
        </DialogActions>
      </Dialog>

      {/* Assign User to Team/Sub-team Dialog */}
      <Dialog open={assignUserDialogOpen} onClose={() => setAssignUserDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle sx={{ pb: 1 }}>
          Assign User to {selectedSubteamForUser ? 'Sub-team' : 'Team'}
        </DialogTitle>
        <DialogContent sx={{ pt: 2 }}>
          {selectedTeamForUser && (
            <Alert severity="info" sx={{ mb: 2 }}>
              {selectedSubteamForUser ? (
                <>
                  Assigning to sub-team: <strong>{selectedSubteamForUser.name}</strong> 
                  (under team: <strong>{selectedTeamForUser.name}</strong>)
                </>
              ) : (
                <>Assigning to team: <strong>{selectedTeamForUser.name}</strong></>
              )}
            </Alert>
          )}
          <FormControl fullWidth margin="normal">
            <InputLabel id="user-select-label">Select User</InputLabel>
            <Select
              labelId="user-select-label"
              value=""
              label="Select User"
              onChange={async (e) => {
                const userId = e.target.value;
                if (!userId) return;

                try {
                  const targetTeamId = selectedSubteamForUser ? selectedSubteamForUser.id : selectedTeamForUser.id;
                  
                  // Get current user to update
                  const user = users.find(u => u.id === userId);
                  if (!user) return;

                  // Update user's team_id
                  await api.put(`/users/${userId}`, {
                    name: user.name,
                    phone: user.phone,
                    role_id: user.role_id,
                    manager_id: user.manager_id,
                    auditor_id: user.auditor_id,
                    team_id: parseInt(targetTeamId), // Assign to selected team/sub-team
                    user_type: user.user_type,
                    location: user.location,
                  });

                  showSnackbar(`User assigned to ${selectedSubteamForUser ? 'sub-team' : 'team'} successfully`);
                  setAssignUserDialogOpen(false);
                  fetchUsers();
                  fetchTeams();
                } catch (error) {
                  showSnackbar(error.response?.data?.error || 'Failed to assign user', 'error');
                }
              }}
            >
              {users
                .filter(u => u.user_type === 'product_user')
                .map((user) => (
                  <MenuItem key={user.id} value={user.id}>
                    {user.name} ({user.email})
                    {user.team_id && (
                      <Typography component="span" variant="caption" color="text.secondary" sx={{ ml: 1 }}>
                        - Currently in: {(() => {
                          const currentTeam = teams.find(t => t.id === user.team_id);
                          if (!currentTeam) return 'Unknown';
                          if (currentTeam.group_id) {
                            const parent = teams.find(t => t.id === currentTeam.group_id);
                            return parent ? `${parent.name} > ${currentTeam.name}` : currentTeam.name;
                          }
                          return currentTeam.name;
                        })()}
                      </Typography>
                    )}
                  </MenuItem>
                ))}
            </Select>
          </FormControl>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setAssignUserDialogOpen(false)}>Cancel</Button>
        </DialogActions>
      </Dialog>

      <Snackbar
        open={snackbar.open}
        autoHideDuration={6000}
        onClose={() => setSnackbar({ ...snackbar, open: false })}
      >
        <Alert severity={snackbar.severity} onClose={() => setSnackbar({ ...snackbar, open: false })}>
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Box>
  );
}

export default Settings;

