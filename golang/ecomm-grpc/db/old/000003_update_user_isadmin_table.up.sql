-- Migration 3: Update users table

-- Step 1: Add new role column
ALTER TABLE users
ADD COLUMN role VARCHAR(50) DEFAULT 'Employee';

-- Step 2: Update the new role column based on the existing is_admin column
UPDATE users
SET role = CASE
    WHEN is_admin = true THEN 'Admin'
    ELSE 'Employee'
END;

-- Step 3: Drop the is_admin column
ALTER TABLE users
DROP COLUMN is_admin;

-- Step 4: Add any necessary indexes or constraints
CREATE INDEX idx_users_role ON users(role);

-- Optional: If you want to enforce specific roles, you can add a check constraint
ALTER TABLE users
ADD CONSTRAINT check_valid_role
CHECK (role IN ('Admin', 'Employee', 'Manager')); -- Add other roles as needed

-- Step 5: Update the updated_at column
UPDATE users
SET updated_at = CURRENT_TIMESTAMP;