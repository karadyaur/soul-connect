SET TIMEZONE = 'UTC';

CREATE TABLE IF NOT EXISTS auth (
    id UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS login_attempts (
    id UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
    username VARCHAR(255) NOT NULL,
    success BOOLEAN NOT NULL,
    attempt_time TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
    user_id UUID REFERENCES auth(id) ON DELETE CASCADE,
    session_token VARCHAR(255) NOT NULL,
    session_expires_at TIMESTAMP NOT NULL
);

-- TRIGGERS

-- Logging successful logs
CREATE OR REPLACE FUNCTION log_login_attempt()
RETURNS TRIGGER AS $$
BEGIN
INSERT INTO login_attempts (username, success)
VALUES (NEW.username, TRUE);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER log_successful_login
    AFTER INSERT ON auth
    FOR EACH ROW
    EXECUTE FUNCTION log_login_attempt();

-- Updating the last time logging
CREATE OR REPLACE FUNCTION update_last_login()
RETURNS TRIGGER AS $$
BEGIN
UPDATE auth SET last_login = NOW() WHERE id = NEW.id;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_login_timestamp
    AFTER INSERT ON auth
    FOR EACH ROW
    EXECUTE FUNCTION update_last_login();

-- Deleting old sessions
CREATE OR REPLACE FUNCTION delete_old_sessions()
RETURNS TRIGGER AS $$
BEGIN
DELETE FROM sessions WHERE user_id = NEW.id AND session_expires_at < NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER clear_old_sessions
    AFTER INSERT OR UPDATE ON auth
                        FOR EACH ROW
                        EXECUTE FUNCTION delete_old_sessions();