SET TIMEZONE = 'UTC';

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
    auth_id UUID NOT NULL,  -- связь с таблицей auth
    full_name VARCHAR(255) NOT NULL,
    bio TEXT,
    photo_link VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- creating a table of subscriptions
CREATE TABLE IF NOT EXISTS subscriptions (
   id UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
   subscriber_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
   author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
   created_at TIMESTAMP DEFAULT NOW()
);

-- TRIGGERS

CREATE OR REPLACE TRIGGER update_users_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_timestamp();
