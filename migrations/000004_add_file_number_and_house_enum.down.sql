-- Drop indexes
DROP INDEX IF EXISTS idx_users_file_number;
DROP INDEX IF EXISTS idx_profiles_house;

-- Revert house column to VARCHAR
ALTER TABLE profiles ALTER COLUMN house TYPE VARCHAR(100) USING house::VARCHAR(100);

-- Drop house enum type
DROP TYPE IF EXISTS house_enum;

-- Remove file_number column from users table
ALTER TABLE users DROP COLUMN IF EXISTS file_number;
