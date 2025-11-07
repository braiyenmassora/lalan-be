-- Membuat ekstensi UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Membuat tabel hosters
CREATE TABLE hosters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- ID unik dengan UUID
    full_name VARCHAR(255) NOT NULL, -- Nama lengkap hoster
    profile_photo VARCHAR(500), -- URL foto profil (opsional)
    store_name VARCHAR(255) NOT NULL, -- Nama toko
    description TEXT, -- Deskripsi toko
    phone_number VARCHAR(20), -- Nomor telepon
    email VARCHAR(255) UNIQUE NOT NULL, -- Email unik
    address TEXT NOT NULL, -- Alamat
    password_hash VARCHAR(255) NOT NULL, -- Hash password
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(), -- Waktu pembuatan
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() -- Waktu pembaruan
);

-- Index untuk performa pencarian email
CREATE INDEX idx_hosters_email ON hosters(email);
-- Index untuk performa pencarian waktu pembuatan
CREATE INDEX idx_hosters_created_at ON hosters(created_at);

-- Fungsi untuk update kolom updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger untuk auto-update updated_at
CREATE TRIGGER update_hosters_updated_at BEFORE UPDATE ON hosters FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();