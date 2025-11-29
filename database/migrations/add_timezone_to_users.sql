-- Add timezone column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS timezone VARCHAR(100) DEFAULT 'UTC';

COMMENT ON COLUMN users.timezone IS 'User timezone (e.g., UTC, America/New_York, Asia/Kolkata)';

