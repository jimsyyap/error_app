-- Create Users Table
CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP,
    CONSTRAINT chk_email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
);

-- Create Match_Sessions Table
CREATE TABLE match_sessions (
    session_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    start_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMP,
    opponent_name VARCHAR(100),
    location VARCHAR(100),
    score VARCHAR(50),
    notes TEXT,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT chk_end_after_start CHECK (end_time IS NULL OR end_time >= start_time)
);

-- Create Error_Types Table
CREATE TABLE error_types (
    error_type_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

-- Create Error_Logs Table
CREATE TABLE error_logs (
    error_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL,
    error_type_id INTEGER NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES match_sessions(session_id) ON DELETE CASCADE,
    FOREIGN KEY (error_type_id) REFERENCES error_types(error_type_id) ON DELETE RESTRICT
);

-- Create indexes for performance
CREATE INDEX idx_error_logs_session ON error_logs(session_id);
CREATE INDEX idx_match_sessions_user ON match_sessions(user_id);
CREATE INDEX idx_error_logs_timestamp ON error_logs(timestamp);

-- Seed Error_Types table with initial values
INSERT INTO error_types (name) VALUES 
('Forehand'), 
('Backhand'), 
('Serve'), 
('Volley'),
('Return'),
('Drop Shot'),
('Overhead'),
('Footwork');
