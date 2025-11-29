import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  AppBar,
  Box,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Toolbar,
  Typography,
  IconButton,
  useTheme,
  useMediaQuery,
  Divider,
  ListSubheader,
  InputBase,
  Avatar,
  Menu,
  MenuItem,
} from '@mui/material';
import MenuIcon from '@mui/icons-material/Menu';
import SearchIcon from '@mui/icons-material/Search';
import NotificationsIcon from '@mui/icons-material/Notifications';
import SettingsIcon from '@mui/icons-material/Settings';
import Badge from '@mui/material/Badge';
import Popover from '@mui/material/Popover';
import ListItemAvatar from '@mui/material/ListItemAvatar';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import ErrorIcon from '@mui/icons-material/Error';
import WarningIcon from '@mui/icons-material/Warning';
import InfoIcon from '@mui/icons-material/Info';
import CircleIcon from '@mui/icons-material/Circle';
import LogoutIcon from '@mui/icons-material/Logout';
import DeleteIcon from '@mui/icons-material/Delete';
import DoneAllIcon from '@mui/icons-material/DoneAll';
import Tooltip from '@mui/material/Tooltip';
import PersonIcon from '@mui/icons-material/Person';
import DashboardIcon from '@mui/icons-material/Dashboard';
import PeopleIcon from '@mui/icons-material/People';
import BusinessIcon from '@mui/icons-material/Business';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import AssessmentIcon from '@mui/icons-material/Assessment';
import StarIcon from '@mui/icons-material/Star';
import SpeedIcon from '@mui/icons-material/Speed';
import GroupWorkIcon from '@mui/icons-material/GroupWork';
import CloudIcon from '@mui/icons-material/Cloud';
import Logo from './Logo';
import { useAuth } from '../contexts/AuthContext';
import { DRAWER_WIDTH, typography } from '../theme/typography';
// import api from '../api/client'; // Uncomment when implementing real notification fetching

const drawerWidth = DRAWER_WIDTH;

const menuItems = [
  { 
    section: 'Core Analytics',
    items: [
      { text: 'Overview', icon: <DashboardIcon />, path: '/' },
      { text: 'Agents', icon: <PeopleIcon />, path: '/agents' },
      { text: 'Vendors', icon: <BusinessIcon />, path: '/vendors' },
      { text: 'Intents', icon: <TrendingUpIcon />, path: '/intents' },
    ]
  },
  {
    section: 'Advanced Analytics',
    items: [
      { text: 'Mix Modeling', icon: <AssessmentIcon />, path: '/mmm' },
      { text: 'Account Marketing', icon: <BusinessIcon />, path: '/abm' },
      { text: 'Lead Scoring', icon: <StarIcon />, path: '/lead-scoring' },
      { text: 'Cohort Analysis', icon: <GroupWorkIcon />, path: '/cohorts' },
      { text: 'Real Time', icon: <SpeedIcon />, path: '/realtime' },
    ]
  },
  {
    section: 'Platform',
    items: [
      { text: 'Settings', icon: <SettingsIcon />, path: '/settings' },
      { text: 'Integrations', icon: <CloudIcon />, path: '/integrations' },
    ]
  }
];

function Layout({ children }) {
  const [mobileOpen, setMobileOpen] = useState(false);
  const [anchorEl, setAnchorEl] = useState(null);
  const [notificationAnchor, setNotificationAnchor] = useState(null);
  const [notifications, setNotifications] = useState([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const navigate = useNavigate();
  const location = useLocation();
  const { user, logout } = useAuth();

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  const handleNavigation = (path) => {
    navigate(path);
    if (isMobile) {
      setMobileOpen(false);
    }
  };

  const handleMenuOpen = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    logout();
    handleMenuClose();
    navigate('/login');
  };

  // Notification handlers
  const handleNotificationOpen = (event) => {
    setNotificationAnchor(event.currentTarget);
    // For design preview, skip API call
    // Uncomment for real API: fetchNotifications();
  };

  const handleNotificationClose = () => {
    setNotificationAnchor(null);
  };

  // Uncomment for real API integration:
  // const fetchNotifications = async () => {
  //   try {
  //     const response = await api.get('/realtime/alerts', {
  //       params: {
  //         limit: 20,
  //         unresolved: true,
  //       }
  //     });
  //     const alerts = response.data.alerts || [];
  //     setNotifications(alerts);
  //     setUnreadCount(alerts.filter(a => !a.acknowledged).length);
  //   } catch (error) {
  //     console.error('Failed to fetch notifications:', error);
  //   }
  // };

  const handleAcknowledgeNotification = async (alertId) => {
    // For design preview, update local state
    setNotifications(prev => 
      prev.map(n => 
        n.id === alertId ? { ...n, acknowledged: true } : n
      )
    );
    setUnreadCount(prev => Math.max(0, prev - 1));
    
    // Uncomment for real API call:
    // try {
    //   await api.post(`/realtime/alerts/${alertId}/acknowledge`, {
    //     by: user?.email || 'user'
    //   });
    //   fetchNotifications();
    // } catch (error) {
    //   console.error('Failed to acknowledge notification:', error);
    // }
  };

  const handleMarkAllAsRead = () => {
    setNotifications(prev => 
      prev.map(n => ({ ...n, acknowledged: true }))
    );
    setUnreadCount(0);
    
    // Uncomment for real API call:
    // notifications.forEach(n => {
    //   if (!n.acknowledged) {
    //     api.post(`/realtime/alerts/${n.id}/acknowledge`, {
    //       by: user?.email || 'user'
    //     });
    //   }
    // });
  };

  const handleRemoveNotification = (alertId, event) => {
    event.stopPropagation(); // Prevent triggering acknowledge
    setNotifications(prev => {
      const removed = prev.find(n => n.id === alertId);
      if (removed && !removed.acknowledged) {
        setUnreadCount(count => Math.max(0, count - 1));
      }
      return prev.filter(n => n.id !== alertId);
    });
    
    // Uncomment for real API call:
    // try {
    //   await api.delete(`/realtime/alerts/${alertId}`);
    //   fetchNotifications();
    // } catch (error) {
    //   console.error('Failed to remove notification:', error);
    // }
  };

  const handleClearAll = () => {
    setNotifications([]);
    setUnreadCount(0);
    
    // Uncomment for real API call:
    // notifications.forEach(n => {
    //   api.delete(`/realtime/alerts/${n.id}`);
    // });
  };

  // Dummy notifications for design preview
  const dummyNotifications = [
    {
      id: 1,
      alert_type: 'system_error',
      severity: 'error',
      title: 'Database Connection Failed',
      description: 'Unable to connect to primary database. Fallback to replica.',
      entity_type: 'system',
      entity_id: null,
      triggered_at: new Date().toISOString(),
      acknowledged: false,
    },
    {
      id: 2,
      alert_type: 'data_quality',
      severity: 'warning',
      title: 'Low Data Quality Score',
      description: 'Data quality score dropped below 80% for the last hour.',
      entity_type: 'data_quality',
      entity_id: null,
      triggered_at: new Date(Date.now() - 5 * 60000).toISOString(), // 5 minutes ago
      acknowledged: false,
    },
    {
      id: 3,
      alert_type: 'attribution_complete',
      severity: 'info',
      title: 'Attribution Run Completed',
      description: 'Q1 2024 attribution analysis has been completed successfully.',
      entity_type: 'attribution_run',
      entity_id: 123,
      triggered_at: new Date(Date.now() - 15 * 60000).toISOString(), // 15 minutes ago
      acknowledged: false,
    },
    {
      id: 4,
      alert_type: 'integration_success',
      severity: 'success',
      title: 'Salesforce Sync Successful',
      description: 'Successfully synced 250 contacts from Salesforce.',
      entity_type: 'integration',
      entity_id: 5,
      triggered_at: new Date(Date.now() - 30 * 60000).toISOString(), // 30 minutes ago
      acknowledged: true,
    },
    {
      id: 5,
      alert_type: 'fraud_detected',
      severity: 'error',
      title: 'Potential Fraud Detected',
      description: 'Unusual pattern detected in conversion events. Requires immediate review.',
      entity_type: 'fraud_incident',
      entity_id: 789,
      triggered_at: new Date(Date.now() - 60 * 60000).toISOString(), // 1 hour ago
      acknowledged: false,
    },
    {
      id: 6,
      alert_type: 'performance_alert',
      severity: 'warning',
      title: 'High API Latency',
      description: 'API response time exceeded 2 seconds for the last 10 minutes.',
      entity_type: 'system',
      entity_id: null,
      triggered_at: new Date(Date.now() - 2 * 3600000).toISOString(), // 2 hours ago
      acknowledged: false,
    },
  ];

  // Fetch notifications on mount and set up polling
  useEffect(() => {
    // For design preview, use dummy notifications
    // In production, uncomment the real fetch:
    // fetchNotifications();
    setNotifications(dummyNotifications);
    setUnreadCount(dummyNotifications.filter(n => !n.acknowledged).length);
    
    // Uncomment for real API polling:
    // const interval = setInterval(fetchNotifications, 30000);
    // return () => clearInterval(interval);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const getNotificationIcon = (severity) => {
    switch (severity?.toLowerCase()) {
      case 'error':
        return <ErrorIcon sx={{ color: '#F93739' }} />;
      case 'warning':
        return <WarningIcon sx={{ color: '#F8AA0D' }} />;
      case 'success':
        return <CheckCircleIcon sx={{ color: '#1AC468' }} />;
      default:
        return <InfoIcon sx={{ color: '#1A62F2' }} />;
    }
  };

  const formatNotificationTime = (timestamp) => {
    if (!timestamp) return 'Just now';
    const date = new Date(timestamp);
    const now = new Date();
    const diffMs = now - date;
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return 'Just now';
    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    if (diffDays < 7) return `${diffDays}d ago`;
    return date.toLocaleDateString();
  };

  const drawer = (
    <Box sx={{ overflow: 'auto' }}>
      <Toolbar sx={{ py: 3, px: 2.5, minHeight: 'auto !important', alignItems: 'center' }}>
        <Logo size="medium" showSubtitle={true} />
      </Toolbar>
      <Divider sx={{ mx: 2, mb: 1 }} />
      {menuItems.map((section, idx) => (
        <List
          key={section.section}
          subheader={
            <ListSubheader 
              component="div" 
              sx={{ 
                bgcolor: 'transparent', 
                ...typography.tableHeader,
                fontSize: '0.75rem',
                textTransform: 'uppercase',
                letterSpacing: '0.05em',
                py: 1.5,
                px: 2,
              }}
            >
              {section.section}
            </ListSubheader>
          }
        >
          {section.items.map((item) => (
            <ListItem key={item.text} disablePadding>
              <ListItemButton
                selected={location.pathname === item.path}
                onClick={() => handleNavigation(item.path)}
                sx={{
                  py: 1.25,
                  px: 2,
                  '&.Mui-selected': {
                    backgroundColor: 'rgba(26, 98, 242, 0.1)',
                    '&:hover': {
                      backgroundColor: 'rgba(26, 98, 242, 0.15)',
                    },
                  },
                }}
              >
                <ListItemIcon sx={{ color: location.pathname === item.path ? 'primary.main' : 'inherit' }}>
                  {item.icon}
                </ListItemIcon>
                <ListItemText 
                  primary={item.text}
                  primaryTypographyProps={{ 
                    ...typography.bodyTextSmall,
                    fontSize: '0.875rem',
                    fontWeight: location.pathname === item.path ? 600 : 400,
                  }}
                />
              </ListItemButton>
            </ListItem>
          ))}
        </List>
      ))}
    </Box>
  );

  return (
    <Box sx={{ display: 'flex' }}>
      <AppBar
        position="fixed"
        sx={{
          width: { md: `calc(100% - ${drawerWidth}px)` },
          ml: { md: `${drawerWidth}px` },
          backgroundColor: 'white',
          color: 'text.primary',
          boxShadow: '0px 1px 3px rgba(0,0,0,0.05)',
        }}
      >
        <Toolbar sx={{ gap: 2 }}>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            edge="start"
            onClick={handleDrawerToggle}
            sx={{ mr: 1, display: { md: 'none' } }}
          >
            <MenuIcon />
          </IconButton>

          {/* Search Bar */}
          <Box
            sx={{
              display: 'flex',
              alignItems: 'center',
              backgroundColor: '#F3F4F6',
              borderRadius: 2,
              px: 2,
              py: 0.75,
              maxWidth: { xs: '100%', md: '600px' },
              width: { xs: '100%', md: '600px' },
            }}
          >
            <SearchIcon sx={{ color: '#9CA3AF', mr: 1 }} />
            <InputBase
              placeholder="Search agents, revenue, campaigns..."
              sx={{
                flex: 1,
                color: 'text.primary',
                fontSize: '0.9rem',
                '& input::placeholder': {
                  color: '#9CA3AF',
                  opacity: 1,
                },
              }}
            />
          </Box>

          {/* Right side icons */}
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, ml: 'auto' }}>
            {/* Notification Center */}
            <IconButton 
              color="inherit" 
              size="small"
              onClick={handleNotificationOpen}
              sx={{ position: 'relative' }}
            >
              <Badge badgeContent={unreadCount} color="error" max={99}>
                <NotificationsIcon />
              </Badge>
            </IconButton>
            
            {/* Notification Popover */}
            <Popover
              open={Boolean(notificationAnchor)}
              anchorEl={notificationAnchor}
              onClose={handleNotificationClose}
              anchorOrigin={{
                vertical: 'bottom',
                horizontal: 'right',
              }}
              transformOrigin={{
                vertical: 'top',
                horizontal: 'right',
              }}
              PaperProps={{
                sx: {
                  mt: 1.5,
                  width: 380,
                  maxHeight: 500,
                  borderRadius: '12px',
                  boxShadow: '0 4px 20px rgba(0,0,0,0.1)',
                },
              }}
            >
              <Box sx={{ p: 2, borderBottom: '1px solid #E5E7EB' }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
                  <Typography variant="h6" sx={{ fontWeight: 600, fontSize: '1rem' }}>
                    Notifications
                  </Typography>
                  {notifications.length > 0 && (
                    <Box sx={{ display: 'flex', gap: 0.5 }}>
                      {unreadCount > 0 && (
                        <Tooltip title="Mark all as read">
                          <IconButton
                            size="small"
                            onClick={handleMarkAllAsRead}
                            sx={{ 
                              color: 'primary.main',
                              '&:hover': { backgroundColor: 'rgba(26, 98, 242, 0.08)' }
                            }}
                          >
                            <DoneAllIcon fontSize="small" />
                          </IconButton>
                        </Tooltip>
                      )}
                      <Tooltip title="Clear all">
                        <IconButton
                          size="small"
                          onClick={handleClearAll}
                          sx={{ 
                            color: '#9CA3AF',
                            '&:hover': { backgroundColor: 'rgba(0, 0, 0, 0.04)', color: '#F93739' }
                          }}
                        >
                          <DeleteIcon fontSize="small" />
                        </IconButton>
                      </Tooltip>
                    </Box>
                  )}
                </Box>
                {unreadCount > 0 && (
                  <Typography variant="caption" sx={{ color: '#9CA3AF' }}>
                    {unreadCount} unread
                  </Typography>
                )}
              </Box>
              
              <Box sx={{ maxHeight: 400, overflowY: 'auto' }}>
                {notifications.length === 0 ? (
                  <Box sx={{ p: 4, textAlign: 'center' }}>
                    <NotificationsIcon sx={{ fontSize: 48, color: '#D1D5DB', mb: 2 }} />
                    <Typography variant="body2" sx={{ color: '#9CA3AF' }}>
                      No notifications
                    </Typography>
                  </Box>
                ) : (
                  <List sx={{ p: 0 }}>
                    {notifications.map((notification) => (
                      <ListItem
                        key={notification.id}
                        sx={{
                          px: 2,
                          py: 1.5,
                          borderBottom: '1px solid #F3F4F6',
                          cursor: 'pointer',
                          backgroundColor: notification.acknowledged ? 'transparent' : 'rgba(26, 98, 242, 0.05)',
                          '&:hover': {
                            backgroundColor: 'rgba(26, 98, 242, 0.08)',
                            '& .notification-actions': {
                              opacity: 1,
                            },
                          },
                        }}
                        onClick={() => {
                          if (!notification.acknowledged) {
                            handleAcknowledgeNotification(notification.id);
                          }
                        }}
                      >
                        <ListItemAvatar>
                          {getNotificationIcon(notification.severity)}
                        </ListItemAvatar>
                        <ListItemText
                          primary={
                            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                              <Typography variant="body2" sx={{ fontWeight: 600, flex: 1 }}>
                                {notification.title || notification.alert_type}
                              </Typography>
                              {!notification.acknowledged && (
                                <CircleIcon sx={{ fontSize: 8, color: 'primary.main' }} />
                              )}
                            </Box>
                          }
                          secondary={
                            <Box>
                              <Typography variant="caption" sx={{ color: '#6B7280', display: 'block' }}>
                                {notification.description || notification.title}
                              </Typography>
                              <Typography variant="caption" sx={{ color: '#9CA3AF', fontSize: '0.7rem' }}>
                                {formatNotificationTime(notification.triggered_at)}
                              </Typography>
                            </Box>
                          }
                        />
                        <Box
                          className="notification-actions"
                          sx={{
                            display: 'flex',
                            gap: 0.5,
                            opacity: 0,
                            transition: 'opacity 0.2s',
                            ml: 1,
                          }}
                          onClick={(e) => e.stopPropagation()}
                        >
                          {!notification.acknowledged && (
                            <Tooltip title="Mark as read">
                              <IconButton
                                size="small"
                                onClick={() => handleAcknowledgeNotification(notification.id)}
                                sx={{
                                  color: 'primary.main',
                                  '&:hover': { backgroundColor: 'rgba(26, 98, 242, 0.1)' }
                                }}
                              >
                                <CheckCircleIcon fontSize="small" />
                              </IconButton>
                            </Tooltip>
                          )}
                          <Tooltip title="Remove">
                            <IconButton
                              size="small"
                              onClick={(e) => handleRemoveNotification(notification.id, e)}
                              sx={{
                                color: '#9CA3AF',
                                '&:hover': { backgroundColor: 'rgba(249, 55, 57, 0.1)', color: '#F93739' }
                              }}
                            >
                              <DeleteIcon fontSize="small" />
                            </IconButton>
                          </Tooltip>
                        </Box>
                      </ListItem>
                    ))}
                  </List>
                )}
              </Box>
              
              {notifications.length > 0 && (
                <Box sx={{ p: 2, borderTop: '1px solid #E5E7EB', textAlign: 'center' }}>
                  <Typography
                    variant="caption"
                    sx={{
                      color: 'primary.main',
                      cursor: 'pointer',
                      fontWeight: 500,
                      '&:hover': { textDecoration: 'underline' },
                    }}
                    onClick={() => navigate('/realtime')}
                  >
                    View All Notifications
                  </Typography>
                </Box>
              )}
            </Popover>

            {/* User Avatar */}
            <IconButton onClick={handleMenuOpen} sx={{ p: 0 }}>
              <Avatar
                sx={{
                  width: 32,
                  height: 32,
                  bgcolor: 'primary.main',
                  fontSize: '0.875rem',
                }}
              >
                {user?.name?.charAt(0).toUpperCase() || 'A'}
              </Avatar>
            </IconButton>
          </Box>

          {/* User Menu */}
          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleMenuClose}
            PaperProps={{
              sx: {
                mt: 1.5,
                minWidth: 200,
                borderRadius: '8px',
                boxShadow: '0 4px 20px rgba(0,0,0,0.1)',
              },
            }}
          >
            <Box sx={{ px: 2, py: 1.5, borderBottom: '1px solid #E5E7EB' }}>
              <Typography variant="body2" sx={{ fontWeight: 600, color: '#1F2937' }}>
                {user?.first_name && user?.last_name 
                  ? `${user.first_name} ${user.last_name}` 
                  : user?.first_name || user?.name || 'User'}
              </Typography>
              <Typography variant="caption" sx={{ color: '#9CA3AF' }}>
                {user?.email || 'user@example.com'}
              </Typography>
            </Box>
            <MenuItem 
              onClick={() => {
                handleMenuClose();
                navigate('/profile');
              }} 
              sx={{ py: 1.5 }}
            >
              <ListItemIcon>
                <PersonIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText>Profile</ListItemText>
            </MenuItem>
            <MenuItem 
              onClick={() => {
                handleMenuClose();
                navigate('/settings');
              }} 
              sx={{ py: 1.5 }}
            >
              <ListItemIcon>
                <SettingsIcon fontSize="small" />
              </ListItemIcon>
              <ListItemText>Settings</ListItemText>
            </MenuItem>
            <Divider />
            <MenuItem onClick={handleLogout} sx={{ py: 1.5, color: 'error.main' }}>
              <ListItemIcon>
                <LogoutIcon fontSize="small" color="error" />
              </ListItemIcon>
              <ListItemText>Logout</ListItemText>
            </MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>
      <Box
        component="nav"
        sx={{ width: { md: drawerWidth }, flexShrink: { md: 0 } }}
      >
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{
            keepMounted: true,
          }}
          sx={{
            display: { xs: 'block', md: 'none' },
            '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
          }}
        >
          {drawer}
        </Drawer>
        <Drawer
          variant="permanent"
          sx={{
            display: { xs: 'none', md: 'block' },
            '& .MuiDrawer-paper': { boxSizing: 'border-box', width: drawerWidth },
          }}
          open
        >
          {drawer}
        </Drawer>
      </Box>
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: 3,
          width: { md: `calc(100% - ${drawerWidth}px)` },
          mt: 8,
        }}
      >
        {children}
      </Box>
    </Box>
  );
}

export default Layout;

