import React from 'react';
import { Box, Typography, Grid, Card, CardContent, Button } from '@mui/material';
import { useAuth } from '../contexts/AuthContext';
import AssessmentIcon from '@mui/icons-material/Assessment';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import TimelineIcon from '@mui/icons-material/Timeline';
import FolderOpenIcon from '@mui/icons-material/FolderOpen';
import BarChartIcon from '@mui/icons-material/BarChart';
import VerifiedUserIcon from '@mui/icons-material/VerifiedUser';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import UploadFileIcon from '@mui/icons-material/UploadFile';
import SearchIcon from '@mui/icons-material/Search';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';
import EditIcon from '@mui/icons-material/Edit';
import MicIcon from '@mui/icons-material/Mic';
import VolumeUpIcon from '@mui/icons-material/VolumeUp';
import FavoriteIcon from '@mui/icons-material/Favorite';
import AttachMoneyIcon from '@mui/icons-material/AttachMoney';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import CampaignIcon from '@mui/icons-material/Campaign';

function Dashboard() {
  const { user } = useAuth();
  
  // Extract first name from user's name
  const getFirstName = () => {
    if (!user) return 'User';
    if (user.first_name) return user.first_name;
    if (user.name) {
      const nameParts = user.name.trim().split(' ');
      return nameParts[0] || 'User';
    }
    return 'User';
  };

  const firstName = getFirstName();

  const quickActions = [
    {
      title: 'Run Attribution',
      description: 'Create or execute attribution analysis',
      icon: AssessmentIcon,
      iconColor: '#6366F1',
      iconBg: 'rgba(99, 102, 241, 0.1)',
      buttonText: 'Start Analysis',
      buttonVariant: 'gradient',
    },
    {
      title: 'Upload Data',
      description: 'Import interaction or conversion data',
      icon: CloudUploadIcon,
      iconColor: '#3B82F6',
      iconBg: 'rgba(59, 130, 246, 0.1)',
      buttonText: 'Upload File',
      buttonVariant: 'outlined',
    },
    {
      title: 'View Analytics',
      description: 'Generate revenue insights',
      icon: TimelineIcon,
      iconColor: '#8B5CF6',
      iconBg: 'rgba(139, 92, 246, 0.1)',
      buttonText: 'Generate',
      buttonVariant: 'outlined',
    },
    {
      title: 'Data Library',
      description: 'Browse historical data',
      icon: FolderOpenIcon,
      iconColor: '#10B981',
      iconBg: 'rgba(16, 185, 129, 0.1)',
      buttonText: 'Explore',
      buttonVariant: 'outlined',
    },
  ];

  const stats = [
    { title: 'Total Revenue', value: '$2.4M', icon: AttachMoneyIcon, color: '#6366F1' },
    { title: 'Conversions', value: '1,247', icon: TrendingUpIcon, color: '#3B82F6' },
    { title: 'Attribution Runs', value: '94', icon: AssessmentIcon, color: '#10B981' },
    { title: 'Active Campaigns', value: '32', icon: CampaignIcon, color: '#F59E0B' },
  ];

  const recentActivity = [
    { 
      title: "Created 'Q4 Attribution Run'", 
      time: '2 hours ago', 
      icon: EditIcon, 
      iconColor: '#6366F1',
      iconBg: 'rgba(99, 102, 241, 0.1)',
    },
    { 
      title: 'Uploaded conversion data', 
      time: '5 hours ago', 
      icon: MicIcon, 
      iconColor: '#3B82F6',
      iconBg: 'rgba(59, 130, 246, 0.1)',
    },
    { 
      title: 'Generated MMM report', 
      time: '1 day ago', 
      icon: VolumeUpIcon, 
      iconColor: '#10B981',
      iconBg: 'rgba(16, 185, 129, 0.1)',
    },
  ];

  return (
    <Box>
      {/* Welcome Section */}
      <Box sx={{ mb: 5 }}>
        <Typography
          variant="h1"
          sx={{
            fontSize: { xs: '1.5rem', sm: '1.75rem', md: '2rem' },
            fontWeight: 700,
            mb: 1,
            color: '#1F2937',
            lineHeight: 1.2,
            '& .username': {
              background: 'linear-gradient(to right, #6366F1, #8B5CF6)',
              WebkitBackgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
              backgroundClip: 'text',
              fontWeight: 700,
            },
          }}
        >
          Welcome back, <span className="username">{firstName}</span>
        </Typography>
        <Typography
          variant="body1"
          sx={{
            color: '#9CA3AF',
            fontSize: '1rem',
            fontWeight: 400,
            lineHeight: 1.5,
          }}
        >
          Here's what's happening with your revenue attribution platform
        </Typography>
      </Box>

      {/* Stats Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        {stats.map((stat, index) => {
          const IconComponent = stat.icon;
          return (
            <Grid item xs={12} sm={6} lg={3} key={index}>
              <Card
                sx={{
                  height: '100%',
                  borderRadius: '12px',
                  border: '1px solid',
                  borderColor: 'rgba(229, 231, 235, 0.8)',
                  backgroundColor: '#FFFFFF',
                  boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
                  transition: 'all 0.2s ease-in-out',
                  '&:hover': {
                    boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06)',
                    borderColor: 'rgba(209, 213, 219, 1)',
                    transform: 'translateY(-2px)',
                  },
                }}
              >
                <CardContent sx={{ p: 3 }}>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                    <Box>
                      <Typography
                        variant="body2"
                        sx={{
                          color: '#9CA3AF',
                          fontSize: '0.875rem',
                          fontWeight: 500,
                          mb: 1,
                        }}
                      >
                        {stat.title}
                      </Typography>
                      <Typography
                        variant="h4"
                        sx={{
                          fontWeight: 700,
                          color: '#1F2937',
                          fontSize: { xs: '1.75rem', sm: '2rem' },
                          lineHeight: 1,
                        }}
                      >
                        {stat.value}
                      </Typography>
                    </Box>
                    <Box
                      sx={{
                        p: 1.5,
                        borderRadius: '10px',
                        backgroundColor: `${stat.color}15`,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                      }}
                    >
                      <IconComponent sx={{ color: stat.color, fontSize: 24 }} />
                    </Box>
                  </Box>
                </CardContent>
              </Card>
            </Grid>
          );
        })}
      </Grid>

      {/* Quick Actions and Recent Activity */}
      <Grid container spacing={3} sx={{ mb: 5 }}>
        {/* Quick Actions - Left Side */}
        <Grid item xs={12} lg={8}>
          <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
            <Typography
              variant="h2"
              sx={{
                fontWeight: 700,
                color: '#1F2937',
                mb: 3,
                fontSize: { xs: '1.25rem', sm: '1.5rem' },
                lineHeight: 1.2,
              }}
            >
              Quick Actions
            </Typography>
            
            <Grid container spacing={3} sx={{ flex: 1 }}>
            {quickActions.map((action, index) => {
              const IconComponent = action.icon;
              return (
                <Grid item xs={12} sm={6} key={index}>
                  <Card
                    sx={{
                      height: '100%',
                      borderRadius: '12px',
                      border: '1px solid',
                      borderColor: 'rgba(229, 231, 235, 0.8)',
                      backgroundColor: '#FFFFFF',
                      boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
                      transition: 'all 0.2s ease-in-out',
                      '&:hover': {
                        boxShadow: '0 8px 16px -4px rgba(0, 0, 0, 0.1)',
                        borderColor: 'rgba(209, 213, 219, 1)',
                        transform: 'translateY(-2px)',
                      },
                    }}
                  >
                    <CardContent sx={{ p: 3 }}>
                      <Box sx={{ display: 'flex', alignItems: 'flex-start', mb: 3 }}>
                        <Box
                          sx={{
                            p: 1.5,
                            borderRadius: '10px',
                            backgroundColor: action.iconBg,
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            mr: 2,
                          }}
                        >
                          <IconComponent sx={{ color: action.iconColor, fontSize: 24 }} />
                        </Box>
                        <Box sx={{ flex: 1 }}>
                          <Typography
                            variant="h6"
                            sx={{
                              fontWeight: 700,
                              color: '#1F2937',
                              fontSize: '1.125rem',
                              mb: 0.5,
                              lineHeight: 1.3,
                            }}
                          >
                            {action.title}
                          </Typography>
                          <Typography
                            variant="body2"
                            sx={{
                              color: '#9CA3AF',
                              fontSize: '0.875rem',
                              lineHeight: 1.5,
                            }}
                          >
                            {action.description}
                          </Typography>
                        </Box>
                      </Box>
                      
                      {action.buttonVariant === 'gradient' ? (
                        <Button
                          fullWidth
                          startIcon={<PlayArrowIcon />}
                          sx={{
                            background: 'linear-gradient(135deg, #6366F1 0%, #8B5CF6 100%)',
                            color: 'white',
                            borderRadius: '8px',
                            textTransform: 'none',
                            fontWeight: 600,
                            fontSize: '0.875rem',
                            py: 1.25,
                            transition: 'all 0.2s ease-in-out',
                            '&:hover': {
                              background: 'linear-gradient(135deg, #5558E3 0%, #7C3AED 100%)',
                              transform: 'translateY(-1px)',
                              boxShadow: '0 4px 12px rgba(99, 102, 241, 0.3)',
                            },
                          }}
                        >
                          {action.buttonText}
                        </Button>
                      ) : (
                        <Button
                          fullWidth
                          startIcon={
                            action.buttonText === 'Upload File' ? <UploadFileIcon /> : 
                            action.buttonText === 'Explore' ? <SearchIcon /> : 
                            <PlayArrowIcon />
                          }
                          variant="outlined"
                          sx={{
                            borderRadius: '8px',
                            textTransform: 'none',
                            fontWeight: 600,
                            fontSize: '0.875rem',
                            borderColor: 'rgba(229, 231, 235, 1)',
                            color: '#1F2937',
                            py: 1.25,
                            transition: 'all 0.2s ease-in-out',
                            '&:hover': {
                              borderColor: action.iconColor,
                              backgroundColor: action.iconBg,
                              transform: 'translateY(-1px)',
                            },
                          }}
                        >
                          {action.buttonText}
                        </Button>
                      )}
                    </CardContent>
                  </Card>
                </Grid>
              );
            })}
            </Grid>
          </Box>
        </Grid>

        {/* Recent Activity - Right Side */}
        <Grid item xs={12} lg={4}>
          <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
            <Typography
              variant="h2"
              sx={{
                fontWeight: 700,
                color: '#1F2937',
                mb: 3,
                fontSize: { xs: '1.25rem', sm: '1.5rem' },
                lineHeight: 1.2,
              }}
            >
              Recent Activity
            </Typography>
            
            <Card
              sx={{
                flex: 1,
                borderRadius: '12px',
                border: '1px solid',
                borderColor: 'rgba(229, 231, 235, 0.8)',
                backgroundColor: '#FFFFFF',
                boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
                display: 'flex',
                flexDirection: 'column',
              }}
            >
              <CardContent sx={{ p: 0, flex: 1, display: 'flex', flexDirection: 'column' }}>
                <Box sx={{ flex: 1 }}>
                  {recentActivity.map((activity, index) => {
                    const IconComponent = activity.icon;
                    return (
                      <Box
                        key={index}
                        sx={{
                          p: 3,
                          display: 'flex',
                          alignItems: 'center',
                          gap: 2,
                          borderBottom: index < recentActivity.length - 1 ? '1px solid rgba(229, 231, 235, 1)' : 'none',
                          transition: 'background-color 0.2s',
                          '&:hover': {
                            backgroundColor: 'rgba(249, 250, 251, 1)',
                          },
                        }}
                      >
                        <Box
                          sx={{
                            p: 1.5,
                            borderRadius: '10px',
                            backgroundColor: activity.iconBg,
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            flexShrink: 0,
                          }}
                        >
                          <IconComponent sx={{ color: activity.iconColor, fontSize: 20 }} />
                        </Box>
                        <Box sx={{ flex: 1, minWidth: 0 }}>
                          <Typography
                            variant="body2"
                            sx={{
                              fontWeight: 600,
                              color: '#1F2937',
                              fontSize: '0.875rem',
                              mb: 0.25,
                              lineHeight: 1.4,
                            }}
                          >
                            {activity.title}
                          </Typography>
                          <Typography
                            variant="caption"
                            sx={{
                              color: '#9CA3AF',
                              fontSize: '0.75rem',
                              fontWeight: 400,
                            }}
                          >
                            {activity.time}
                          </Typography>
                        </Box>
                      </Box>
                    );
                  })}
                </Box>
                
                <Box
                  sx={{
                    p: 2.5,
                    textAlign: 'center',
                    borderTop: '1px solid rgba(229, 231, 235, 1)',
                    mt: 'auto',
                  }}
                >
                <Button
                  endIcon={<ArrowForwardIcon />}
                  sx={{
                    textTransform: 'none',
                    color: '#1F2937',
                    fontWeight: 600,
                    fontSize: '0.875rem',
                    borderRadius: '8px',
                    px: 2,
                    transition: 'all 0.2s ease-in-out',
                    '&:hover': {
                      backgroundColor: 'rgba(249, 250, 251, 1)',
                      '& .MuiSvgIcon-root': {
                        transform: 'translateX(4px)',
                      },
                    },
                    '& .MuiSvgIcon-root': {
                      transition: 'transform 0.2s',
                    },
                  }}
                >
                  View All Activity
                </Button>
              </Box>
            </CardContent>
          </Card>
          </Box>
        </Grid>
      </Grid>

      {/* Bottom Large Cards */}
      <Grid container spacing={3}>
        {/* Analytics & Insights */}
        <Grid item xs={12} md={6}>
          <Card
            sx={{
              borderRadius: '12px',
              border: '1px solid',
              borderColor: 'rgba(229, 231, 235, 0.8)',
              backgroundColor: '#FFFFFF',
              boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
              transition: 'all 0.2s ease-in-out',
              '&:hover': {
                boxShadow: '0 8px 16px -4px rgba(0, 0, 0, 0.1)',
                transform: 'translateY(-2px)',
              },
            }}
          >
            <CardContent sx={{ p: 4 }}>
              <Box
                sx={{
                  p: 1.5,
                  borderRadius: '10px',
                  backgroundColor: 'rgba(59, 130, 246, 0.1)',
                  display: 'inline-flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  mb: 3,
                }}
              >
                <BarChartIcon sx={{ color: '#3B82F6', fontSize: 28 }} />
              </Box>
              
              <Typography
                variant="h5"
                sx={{
                  fontWeight: 700,
                  color: '#1F2937',
                  mb: 1,
                  fontSize: '1.25rem',
                }}
              >
                Analytics & Insights
              </Typography>
              
              <Typography
                variant="body2"
                sx={{
                  color: '#9CA3AF',
                  fontSize: '0.875rem',
                  mb: 3,
                  lineHeight: 1.6,
                }}
              >
                Track performance metrics and revenue patterns
              </Typography>
              
              <Button
                endIcon={<ArrowForwardIcon />}
                sx={{
                  textTransform: 'none',
                  color: '#1F2937',
                  fontWeight: 600,
                  fontSize: '0.875rem',
                  borderRadius: '8px',
                  px: 0,
                  transition: 'all 0.2s ease-in-out',
                  '&:hover': {
                    backgroundColor: 'transparent',
                    '& .MuiSvgIcon-root': {
                      transform: 'translateX(4px)',
                    },
                  },
                  '& .MuiSvgIcon-root': {
                    transition: 'transform 0.2s',
                  },
                }}
              >
                View Dashboard
              </Button>
            </CardContent>
          </Card>
        </Grid>

        {/* Data Quality Manager */}
        <Grid item xs={12} md={6}>
          <Card
            sx={{
              borderRadius: '12px',
              border: '1px solid',
              borderColor: 'rgba(229, 231, 235, 0.8)',
              backgroundColor: '#FFFFFF',
              boxShadow: '0 1px 2px 0 rgba(0, 0, 0, 0.05)',
              transition: 'all 0.2s ease-in-out',
              '&:hover': {
                boxShadow: '0 8px 16px -4px rgba(0, 0, 0, 0.1)',
                transform: 'translateY(-2px)',
              },
            }}
          >
            <CardContent sx={{ p: 4 }}>
              <Box
                sx={{
                  p: 1.5,
                  borderRadius: '10px',
                  backgroundColor: 'rgba(16, 185, 129, 0.1)',
                  display: 'inline-flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  mb: 3,
                }}
              >
                <VerifiedUserIcon sx={{ color: '#10B981', fontSize: 28 }} />
              </Box>
              
              <Typography
                variant="h5"
                sx={{
                  fontWeight: 700,
                  color: '#1F2937',
                  mb: 1,
                  fontSize: '1.25rem',
                }}
              >
                Data Quality Manager
              </Typography>
              
              <Typography
                variant="body2"
                sx={{
                  color: '#9CA3AF',
                  fontSize: '0.875rem',
                  mb: 3,
                  lineHeight: 1.6,
                }}
              >
                Monitor data integrity and accuracy
              </Typography>
              
              <Button
                endIcon={<ArrowForwardIcon />}
                sx={{
                  textTransform: 'none',
                  color: '#1F2937',
                  fontWeight: 600,
                  fontSize: '0.875rem',
                  borderRadius: '8px',
                  px: 0,
                  transition: 'all 0.2s ease-in-out',
                  '&:hover': {
                    backgroundColor: 'transparent',
                    '& .MuiSvgIcon-root': {
                      transform: 'translateX(4px)',
                    },
                  },
                  '& .MuiSvgIcon-root': {
                    transition: 'transform 0.2s',
                  },
                }}
              >
                Manage
              </Button>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
}

export default Dashboard;
