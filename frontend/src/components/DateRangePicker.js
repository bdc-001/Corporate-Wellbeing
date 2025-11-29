import React, { useState } from 'react';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { TextField, Popover, Box } from '@mui/material';

function DateRangePicker({ value, onChange, label = 'Date Range', ...props }) {
  const [startDate, endDate] = value || [null, null];
  const [anchorEl, setAnchorEl] = useState(null);
  const open = Boolean(anchorEl);

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleStartDateChange = (newValue) => {
    onChange([newValue, endDate]);
  };

  const handleEndDateChange = (newValue) => {
    onChange([startDate, newValue]);
  };

  const formatDateRange = () => {
    if (!startDate && !endDate) return '';
    const formatDate = (date) => {
      if (!date) return '';
      return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
    };
    return `${formatDate(startDate)} - ${formatDate(endDate)}`;
  };

  return (
    <LocalizationProvider dateAdapter={AdapterDateFns}>
      <TextField
        label={label}
        value={formatDateRange()}
        fullWidth
        onClick={handleClick}
        InputProps={{
          readOnly: true,
        }}
      />
      <Popover
        open={open}
        anchorEl={anchorEl}
        onClose={handleClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'left',
        }}
      >
        <Box sx={{ p: 2, display: 'flex', gap: 2, minWidth: 500 }}>
          <DatePicker
            label="From Date"
            value={startDate}
            onChange={handleStartDateChange}
            slotProps={{ textField: { size: 'small' } }}
            {...props}
          />
          <DatePicker
            label="To Date"
            value={endDate}
            onChange={handleEndDateChange}
            slotProps={{ textField: { size: 'small' } }}
            {...props}
          />
        </Box>
      </Popover>
    </LocalizationProvider>
  );
}

export default DateRangePicker;

