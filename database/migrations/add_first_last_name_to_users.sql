-- Add first_name and last_name columns to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS first_name VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_name VARCHAR(255);

-- Migrate existing data: split name into first_name and last_name
UPDATE users 
SET 
  first_name = CASE 
    WHEN name IS NULL OR name = '' THEN NULL
    WHEN position(' ' in name) = 0 THEN name
    ELSE substring(name from 1 for position(' ' in name) - 1)
  END,
  last_name = CASE 
    WHEN name IS NULL OR name = '' THEN NULL
    WHEN position(' ' in name) = 0 THEN NULL
    ELSE substring(name from position(' ' in name) + 1)
  END
WHERE first_name IS NULL OR last_name IS NULL;

-- Make first_name NOT NULL after migration (optional, can keep nullable)
-- ALTER TABLE users ALTER COLUMN first_name SET NOT NULL;

COMMENT ON COLUMN users.first_name IS 'User first name';
COMMENT ON COLUMN users.last_name IS 'User last name';

