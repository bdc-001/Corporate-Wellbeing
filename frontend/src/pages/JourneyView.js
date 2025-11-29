import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import {
  Box,
  Typography,
  Card,
  CardContent,
  Chip,
  Grid,
  Divider,
  Stack,
  Avatar,
} from '@mui/material';
import PhoneIcon from '@mui/icons-material/Phone';
import ShoppingCartIcon from '@mui/icons-material/ShoppingCart';
import EmailIcon from '@mui/icons-material/Email';
import client from '../api/client';

function JourneyView() {
  const { customerId } = useParams();
  const [journey, setJourney] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchJourney();
  }, [customerId]);

  const fetchJourney = async () => {
    setLoading(true);
    try {
      const response = await client.get(`/customers/${customerId}/journey`);
      setJourney(response.data);
    } catch (error) {
      console.error('Error fetching journey:', error);
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (amount, currency) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency || 'USD',
    }).format(amount);
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleString();
  };

  if (loading) {
    return <Typography>Loading...</Typography>;
  }

  if (!journey) {
    return <Typography>Customer journey not found</Typography>;
  }

  // Combine interactions and conversions into a single timeline
  const timelineItems = [];

  journey.interactions?.forEach((interaction) => {
    timelineItems.push({
      type: 'interaction',
      data: interaction,
      timestamp: interaction.started_at,
    });
  });

  journey.conversion_events?.forEach((conversion) => {
    timelineItems.push({
      type: 'conversion',
      data: conversion,
      timestamp: conversion.occurred_at,
    });
  });

  timelineItems.sort((a, b) => new Date(a.timestamp) - new Date(b.timestamp));

  return (
    <Box>
      <Typography variant="h4" gutterBottom sx={{ fontWeight: 700, mb: 4 }}>
        Customer Journey
      </Typography>

      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Customer ID: {journey.customer_id}
          </Typography>
          <Box sx={{ mt: 2 }}>
            <Typography variant="subtitle2" color="text.secondary">
              Identifiers:
            </Typography>
            <Box sx={{ mt: 1, display: 'flex', gap: 1, flexWrap: 'wrap' }}>
              {journey.identifiers?.map((ident, idx) => (
                <Chip
                  key={idx}
                  label={`${ident.type}: ${ident.value}`}
                  size="small"
                  variant="outlined"
                />
              ))}
            </Box>
          </Box>
        </CardContent>
      </Card>

      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Timeline
          </Typography>
          <Stack spacing={2} sx={{ mt: 2 }}>
            {timelineItems.map((item, idx) => (
              <Box key={idx}>
                <Box sx={{ display: 'flex', gap: 2 }}>
                  <Avatar
                    sx={{
                      bgcolor: item.type === 'conversion' ? 'success.main' : 'primary.main',
                      width: 40,
                      height: 40,
                    }}
                  >
                    {item.type === 'conversion' ? (
                      <ShoppingCartIcon />
                    ) : (
                      <PhoneIcon />
                    )}
                  </Avatar>
                  <Box sx={{ flex: 1 }}>
                    <Card variant="outlined">
                      <CardContent>
                        <Typography variant="subtitle2" color="text.secondary">
                          {formatDate(item.timestamp)}
                        </Typography>
                        {item.type === 'interaction' ? (
                          <Box>
                            <Typography variant="h6" sx={{ mt: 1 }}>
                              {item.data.channel_name} Interaction
                            </Typography>
                            {item.data.primary_intent && (
                              <Chip
                                label={item.data.primary_intent}
                                size="small"
                                sx={{ mt: 1 }}
                              />
                            )}
                            {item.data.duration_seconds && (
                              <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                                Duration: {Math.floor(item.data.duration_seconds / 60)}m{' '}
                                {item.data.duration_seconds % 60}s
                              </Typography>
                            )}
                          </Box>
                        ) : (
                          <Box>
                            <Typography variant="h6" sx={{ mt: 1 }}>
                              {item.data.event_type} - {formatCurrency(item.data.amount_decimal, item.data.currency_code)}
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                              Source: {item.data.event_source_name}
                            </Typography>
                          </Box>
                        )}
                      </CardContent>
                    </Card>
                  </Box>
                </Box>
                {idx < timelineItems.length - 1 && (
                  <Box sx={{ display: 'flex', justifyContent: 'center', my: 1 }}>
                    <Divider orientation="vertical" sx={{ height: 20, borderWidth: 2 }} />
                  </Box>
                )}
              </Box>
            ))}
          </Stack>
        </CardContent>
      </Card>
    </Box>
  );
}

export default JourneyView;

