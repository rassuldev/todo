-- Create databases for all services
CREATE DATABASE IF NOT EXISTS user_db;
CREATE DATABASE IF NOT EXISTS task_db;
CREATE DATABASE IF NOT EXISTS notification_db;

-- Note: PostgreSQL doesn't support CREATE DATABASE IF NOT EXISTS in the same way
-- If databases don't exist, they will be created by the services themselves
-- This file is mainly for reference and can be expanded for more complex initialization

