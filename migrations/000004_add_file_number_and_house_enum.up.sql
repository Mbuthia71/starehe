-- Add file_number column to users table
ALTER TABLE users ADD COLUMN file_number VARCHAR(10);

-- Create house enum type
CREATE TYPE house_enum AS ENUM (
    'Gikubu',
    'Ngala',
    'Geturo',
    'Shaw',
    'Horsten',
    'Mboya',
    'Shell',
    'Chaka',
    'Njonjo',
    'Kirkley',
    'Muriuki',
    'Kibaki'
);

-- Update house column to use the new enum type
ALTER TABLE profiles ALTER COLUMN house TYPE house_enum USING house::house_enum;

-- Add index for file_number
CREATE INDEX idx_users_file_number ON users(file_number);

-- Add index for house
CREATE INDEX idx_profiles_house ON profiles(house);
