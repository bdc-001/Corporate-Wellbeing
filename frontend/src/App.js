import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { AuthProvider } from './contexts/AuthContext';
import ProtectedRoute from './components/ProtectedRoute';
import Layout from './components/Layout';
import Login from './pages/Login';
import Signup from './pages/Signup';
import Dashboard from './pages/Dashboard';
import AgentsDashboard from './pages/AgentsDashboard';
import VendorsDashboard from './pages/VendorsDashboard';
import IntentsDashboard from './pages/IntentsDashboard';
import JourneyView from './pages/JourneyView';
import MMMDashboard from './pages/MMMDashboard';
import ABMDashboard from './pages/ABMDashboard';
import LeadScoringDashboard from './pages/LeadScoringDashboard';
import RealtimeDashboard from './pages/RealtimeDashboard';
import CohortDashboard from './pages/CohortDashboard';
import IntegrationsDashboard from './pages/IntegrationsDashboard';
import Settings from './pages/Settings';
import Profile from './pages/Profile';

// Convin Design System Theme
const theme = createTheme({
  palette: {
    primary: {
      main: '#1A62F2', // Primary Blue
    },
    secondary: {
      main: '#F030FE', // Accent Purple
    },
    success: {
      main: '#1AC468', // Success Green
    },
    error: {
      main: '#F93739', // Error Red
    },
    warning: {
      main: '#F8AA0D', // Accent Yellow
    },
    background: {
      default: '#FFFFFF',
      paper: '#FFFFFF',
    },
    text: {
      primary: '#333333',
      secondary: '#666666',
    },
  },
  typography: {
    fontFamily: '"Figtree", "Inter", sans-serif',
    h1: {
      fontWeight: 800,
    },
    h2: {
      fontWeight: 700,
    },
    h3: {
      fontWeight: 600,
    },
    button: {
      textTransform: 'none',
      fontWeight: 500,
    },
  },
  shape: {
    borderRadius: 12,
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          padding: '8px 16px',
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 12,
          boxShadow: '0px 2px 4px rgba(0,0,0,0.05)',
        },
      },
    },
  },
});

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <AuthProvider>
        <Router>
          <Routes>
            {/* Public Routes */}
            <Route path="/login" element={<Login />} />
            <Route path="/signup" element={<Signup />} />

            {/* Protected Routes */}
            <Route
              path="/*"
              element={
                <ProtectedRoute>
                  <Layout>
                    <Routes>
                      <Route path="/" element={<Dashboard />} />
                      <Route path="/agents" element={<AgentsDashboard />} />
                      <Route path="/vendors" element={<VendorsDashboard />} />
                      <Route path="/intents" element={<IntentsDashboard />} />
                      <Route path="/journey/:customerId" element={<JourneyView />} />
                      <Route path="/mmm" element={<MMMDashboard />} />
                      <Route path="/abm" element={<ABMDashboard />} />
                      <Route path="/lead-scoring" element={<LeadScoringDashboard />} />
                      <Route path="/realtime" element={<RealtimeDashboard />} />
                      <Route path="/cohorts" element={<CohortDashboard />} />
                      <Route path="/settings" element={<Settings />} />
                      <Route path="/profile" element={<Profile />} />
                      <Route path="/integrations" element={<IntegrationsDashboard />} />
                    </Routes>
                  </Layout>
                </ProtectedRoute>
              }
            />
          </Routes>
        </Router>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;

