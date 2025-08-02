-- Initial database setup for Task Manager
-- This script will be executed when PostgreSQL container starts

-- Create additional schemas if needed
-- CREATE SCHEMA IF NOT EXISTS task_manager;

-- You can add initial data here
-- Example:
-- INSERT INTO users (id, username, email, password, first_name, last_name, created_at, updated_at) 
-- VALUES (
--     gen_random_uuid(),
--     'admin',
--     'admin@taskmanager.com',
--     '$2a$10$example_hashed_password',
--     'Admin',
--     'User',
--     NOW(),
--     NOW()
-- );

-- Create indexes for better performance
-- These will be created automatically by GORM, but you can add custom ones here

-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable pgcrypto for password hashing functions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
