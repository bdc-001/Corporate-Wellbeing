import React, { useState } from 'react';
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
import LogoutIcon from '@mui/icons-material/Logout';
import PersonIcon from '@mui/icons-material/Person';
import DashboardIcon from '@mui/icons-material/Dashboard';
import PeopleIcon from '@mui/icons-material/People';
import BusinessIcon from '@mui/icons-material/Business';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import TimelineIcon from '@mui/icons-material/Timeline';
import AssessmentIcon from '@mui/icons-material/Assessment';
import StarIcon from '@mui/icons-material/Star';
import SpeedIcon from '@mui/icons-material/Speed';
import GroupWorkIcon from '@mui/icons-material/GroupWork';
import CloudIcon from '@mui/icons-material/Cloud';
import Logo from './Logo';
import { useAuth } from '../contexts/AuthContext';

const drawerWidth = 280;

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

  const drawer = (
    <Box sx={{ overflow: 'auto' }}>
      <Toolbar sx={{ py: 2, px: 2, minHeight: '80px !important' }}>
        <Logo size="medium" showSubtitle={true} />
      </Toolbar>
      <Divider sx={{ mx: 2 }} />
      {menuItems.map((section, idx) => (
        <List
          key={section.section}
          subheader={
            <ListSubheader component="div" sx={{ bgcolor: 'transparent', fontWeight: 600 }}>
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
                  primaryTypographyProps={{ fontSize: '0.875rem' }}
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
            <IconButton color="inherit" size="small">
              <NotificationsIcon />
            </IconButton>
            <IconButton color="inherit" size="small">
              <SettingsIcon />
            </IconButton>
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
            <MenuItem onClick={handleMenuClose} sx={{ py: 1.5 }}>
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

