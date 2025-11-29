// Standardized Typography Constants
// Based on Design System: Figtree (Primary), Inter (Fallback), JetBrains Mono (Code)

// Font Families
export const fonts = {
  primary: '"Figtree", "Inter", sans-serif',
  monospace: '"JetBrains Mono", "Fira Code", "Consolas", "Monaco", "Courier New", monospace',
  serif: 'Georgia, serif', // Reserved for special cases
};

// Font Sizes (matching design system scale)
export const fontSizes = {
  xs: '0.75rem',      // 12px - Captions, labels, small text
  sm: '0.875rem',     // 14px - Body text, descriptions, secondary content
  base: '1rem',      // 16px - Default body text
  lg: '1.125rem',    // 18px - Emphasized body text
  xl: '1.25rem',     // 20px - Subheadings
  '2xl': '1.5rem',   // 24px - Section headings
  '3xl': '1.875rem', // 30px - Page headings
  '4xl': '2.25rem',  // 36px - Hero subheadings
  '5xl': '3rem',     // 48px - Hero headings (mobile)
  '6xl': '3.75rem',  // 60px - Hero headings (desktop)
  '7xl': '4.5rem',   // 72px - Large hero headings
};

// Font Weights
export const fontWeights = {
  light: 300,      // Rarely used, decorative text
  regular: 400,    // Body text, descriptions
  medium: 500,     // Emphasized body text
  semibold: 600,   // Subheadings, labels
  bold: 700,       // Headings, important text
  extrabold: 800,  // Hero headings, strong emphasis
  black: 900,      // Rarely used, maximum emphasis
};

// Line Heights
export const lineHeights = {
  tight: 1.25,     // Headings, short lines
  normal: 1.5,     // Body text, paragraphs
  relaxed: 1.75,   // Long-form content, descriptions
};

// Letter Spacing
export const letterSpacing = {
  normal: '0em',      // Default for most text
  wide: '0.05em',     // Uppercase labels, badges
  tight: '-0.02em',   // Large headings
};

// Typography Styles
export const typography = {
  // Hero Heading
  heroHeading: {
    fontFamily: fonts.primary,
    fontSize: { xs: fontSizes['4xl'], sm: fontSizes['5xl'], lg: fontSizes['6xl'], xl: fontSizes['7xl'] },
    fontWeight: fontWeights.bold,
    lineHeight: lineHeights.tight,
    letterSpacing: letterSpacing.tight,
    color: '#1F2937',
  },
  
  // Page Headers (h1)
  pageTitle: {
    fontFamily: fonts.primary,
    fontSize: { xs: fontSizes['2xl'], sm: fontSizes['3xl'], md: fontSizes['4xl'] },
    fontWeight: fontWeights.bold,
    lineHeight: lineHeights.tight,
    letterSpacing: letterSpacing.tight,
    color: '#1F2937',
  },
  
  // Section Headers (h2)
  sectionHeader: {
    fontFamily: fonts.primary,
    fontSize: { xs: fontSizes.xl, sm: fontSizes['2xl'] },
    fontWeight: fontWeights.bold,
    lineHeight: lineHeights.tight,
    color: '#1F2937',
  },
  
  // Card Titles (h3/h5/h6)
  cardTitle: {
    fontFamily: fonts.primary,
    fontSize: { xs: fontSizes.lg, sm: fontSizes.xl },
    fontWeight: fontWeights.bold,
    lineHeight: lineHeights.normal,
    color: '#1F2937',
  },
  
  // Card Subtitle
  cardSubtitle: {
    fontFamily: fonts.primary,
    fontSize: fontSizes.sm,
    fontWeight: fontWeights.medium,
    lineHeight: lineHeights.normal,
    color: '#9CA3AF',
  },
  
  // Stat Value (large numbers)
  statValue: {
    fontFamily: fonts.primary,
    fontSize: { xs: fontSizes['3xl'], sm: fontSizes['4xl'] },
    fontWeight: fontWeights.bold,
    lineHeight: 1,
    color: '#1F2937',
  },
  
  // Body Text
  bodyText: {
    fontFamily: fonts.primary,
    fontSize: { xs: fontSizes.base, sm: fontSizes.lg },
    fontWeight: fontWeights.regular,
    lineHeight: lineHeights.relaxed,
    color: '#1F2937',
  },
  
  // Body Text Small
  bodyTextSmall: {
    fontFamily: fonts.primary,
    fontSize: fontSizes.sm,
    fontWeight: fontWeights.regular,
    lineHeight: lineHeights.normal,
    color: '#6B7280',
  },
  
  // Caption Text
  caption: {
    fontFamily: fonts.primary,
    fontSize: fontSizes.xs,
    fontWeight: fontWeights.regular,
    lineHeight: 1.4,
    color: '#9CA3AF',
  },
  
  // Table Header
  tableHeader: {
    fontFamily: fonts.primary,
    fontSize: fontSizes.sm,
    fontWeight: fontWeights.semibold,
    lineHeight: lineHeights.normal,
    letterSpacing: letterSpacing.wide,
    textTransform: 'uppercase',
    color: '#374151',
  },
  
  // Table Cell
  tableCell: {
    fontFamily: fonts.primary,
    fontSize: fontSizes.sm,
    fontWeight: fontWeights.regular,
    lineHeight: lineHeights.normal,
    color: '#1F2937',
  },
  
  // Code/API Text
  code: {
    fontFamily: fonts.monospace,
    fontSize: fontSizes.sm,
    fontWeight: fontWeights.regular,
    lineHeight: lineHeights.normal,
    color: '#1F2937',
  },
};

// Spacing Scale (4px base unit)
export const spacing = {
  0: '0px',
  1: '4px',    // Tight spacing, icon padding
  2: '8px',    // Small gaps, button padding
  3: '12px',   // Component internal spacing
  4: '16px',   // Standard spacing, card padding
  6: '24px',   // Section spacing, card gaps
  8: '32px',   // Large component spacing
  12: '48px',  // Section gaps
  16: '64px',  // Major section spacing
  20: '80px',  // Hero section padding
  24: '96px',  // Extra large spacing
  32: '128px', // Maximum spacing
};

// Navigation Drawer Width
export const DRAWER_WIDTH = 240; // Reduced from 280px for better space utilization

