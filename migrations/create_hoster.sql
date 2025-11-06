-- Up migration: create hosters table
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE hosters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name VARCHAR(255) NOT NULL,
    profile_photo VARCHAR(500),  -- optional, untuk URL gambar
    store_name VARCHAR(255) NOT NULL,
    description TEXT,
    phone_number VARCHAR(20),
    email VARCHAR(255) UNIQUE NOT NULL,
    address TEXT NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index untuk performa
CREATE INDEX idx_hosters_email ON hosters(email);
CREATE INDEX idx_hosters_created_at ON hosters(created_at);

-- Trigger untuk auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_hosters_updated_at BEFORE UPDATE ON hosters FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();