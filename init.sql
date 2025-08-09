-- Initialize expense database
CREATE DATABASE IF NOT EXISTS expense_db;

-- Connect to the database
\c expense_db;

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Note: Tables will be created by GORM migrations
-- This file is mainly for database initialization 