CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS users (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(200) NOT NULL UNIQUE,
    name VARCHAR(80) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    create_at TIMESTAMP WITH TIME ZONE NOT NULL,
    update_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX email_users_idx ON users(email);

CREATE OR REPLACE FUNCTION auto_update_time()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
--$$ language 'plpgsql';

CREATE TRIGGER auto_update_users_trig
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION auto_update_time();

