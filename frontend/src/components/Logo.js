import React from 'react';
import { Box, Typography } from '@mui/material';

function Logo({ size = 'medium', showSubtitle = true }) {
  const sizes = {
    small: { icon: 28, title: '0.95rem', subtitle: '0.65rem' },
    medium: { icon: 36, title: '1.1rem', subtitle: '0.7rem' },
    large: { icon: 64, title: '2rem', subtitle: '1rem' },
  };

  const currentSize = sizes[size];

  return (
    <Box sx={{ display: 'flex', alignItems: 'flex-start', gap: 1 }}>
      {/* Icon */}
      <Box
        sx={{
          width: currentSize.icon,
          height: currentSize.icon,
          borderRadius: 2.5,
          background: 'linear-gradient(135deg, #1A62F2 0%, #8B5CF6 100%)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          boxShadow: '0 2px 8px rgba(26, 98, 242, 0.2)',
          position: 'relative',
        }}
      >
        <svg
          width={currentSize.icon * 0.7}
          height={currentSize.icon * 0.7}
          viewBox="0 0 24 24"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          {/* Trending up revenue chart - Economics theme */}
          <path
            d="M4 18L8 13L12 15L20 7"
            stroke="white"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            fill="none"
          />
          {/* Data points on chart */}
          <circle cx="8" cy="13" r="1.5" fill="white" />
          <circle cx="12" cy="15" r="1.5" fill="white" />
          <circle cx="20" cy="7" r="1.5" fill="white" />
          {/* Currency symbol accent (₹) - Finance theme, positioned elegantly */}
          <path
            d="M15 5C15 5 15.4 4.6 16.2 4.6C17 4.6 17.4 5 17.4 5.8C17.4 6.6 17 7 16.2 7C15.4 7 15 7.4 15 8.2M15 9.8V10.6M15 4.2V5"
            stroke="white"
            strokeWidth="1.2"
            strokeLinecap="round"
            strokeLinejoin="round"
            opacity="0.9"
          />
        </svg>
      </Box>

      {/* Text */}
      <Box sx={{ pt: 0.3, flex: 1, minWidth: 0, overflow: 'hidden' }}>
        <Typography
          variant="h6"
          component="div"
          sx={{
            fontWeight: 700,
            fontSize: currentSize.title,
            lineHeight: 1.2,
            color: '#1F2937',
            fontFamily: '"Figtree", sans-serif',
            letterSpacing: '-0.01em',
            mb: showSubtitle ? 0.2 : 0,
            whiteSpace: 'nowrap',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
          }}
        >
          Economics
        </Typography>
        {showSubtitle && (
          <Typography
            variant="caption"
            component="div"
            sx={{
              fontSize: currentSize.subtitle,
              color: '#9CA3AF',
              fontWeight: 400,
              lineHeight: 1.3,
              fontFamily: '"Figtree", sans-serif',
              display: 'block',
              whiteSpace: 'nowrap',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
            }}
          >
            Revenue Attribution Platform
          </Typography>
        )}
      </Box>
    </Box>
  );
}

export default Logo;

