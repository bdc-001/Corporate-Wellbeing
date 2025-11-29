-- ============================================
-- UPDATE USER TYPES
-- ============================================
-- Add product_user and observer user types

-- Update existing users with 'standard' type to 'product_user'
UPDATE users SET user_type = 'product_user' WHERE user_type = 'standard';

-- Update existing users with 'admin' type to 'product_user' (admins are product users with admin role)
UPDATE users SET user_type = 'product_user' WHERE user_type = 'admin';

-- Update existing users with 'auditor' type to 'observer'
UPDATE users SET user_type = 'observer' WHERE user_type = 'auditor';

-- Add comment to clarify user types
COMMENT ON COLUMN users.user_type IS 'User type: product_user (can access based on role permissions) or observer (no functional access)';

