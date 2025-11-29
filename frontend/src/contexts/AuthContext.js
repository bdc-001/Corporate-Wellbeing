import React, { createContext, useState, useContext, useEffect } from 'react';
import api from '../api/client';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Check if user is logged in (from localStorage)
    const storedUser = localStorage.getItem('convin_user');
    if (storedUser) {
      try {
        setUser(JSON.parse(storedUser));
      } catch (e) {
        console.error('Error parsing stored user:', e);
        localStorage.removeItem('convin_user');
      }
    }
    setLoading(false);
  }, []);

  const login = async (email, password) => {
    try {
      const response = await api.post('/users/login', { email, password });
      const userData = response.data.user || response.data;
      
      if (!userData || !userData.id) {
        return { success: false, error: 'Invalid response from server' };
      }

      localStorage.setItem('convin_user', JSON.stringify(userData));
      setUser(userData);
      return { success: true };
    } catch (error) {
      const errorMessage = error.response?.data?.error || error.message || 'Login failed';
      return { success: false, error: errorMessage };
    }
  };

  const signup = async (firstName, lastName, email, password) => {
    try {
      // For signup, we'll create a user and then log them in
      // Note: You may want to create a separate signup endpoint
      const response = await api.post('/users', {
        email,
        first_name: firstName,
        last_name: lastName || null,
        password,
        user_type: 'product_user',
      });

      // Check if user was created successfully
      if (!response.data || !response.data.user) {
        return { success: false, error: 'Failed to create user account' };
      }

      // After signup, log them in
      return await login(email, password);
    } catch (error) {
      // Handle different error status codes
      if (error.response?.status === 409) {
        return { success: false, error: 'An account with this email already exists. Please sign in instead.' };
      }
      const errorMessage = error.response?.data?.error || error.message || 'Signup failed. Please try again.';
      return { success: false, error: errorMessage };
    }
  };

  const logout = () => {
    localStorage.removeItem('convin_user');
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, login, signup, logout, loading }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

