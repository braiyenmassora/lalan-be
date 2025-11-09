/*
Mengaktifkan ekstensi UUID untuk menghasilkan UUID.
*/
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

/*
Membuat tabel hosters untuk menyimpan data hoster.
*/
CREATE TABLE hosters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name VARCHAR(255) NOT NULL,
    profile_photo VARCHAR(500),
    store_name VARCHAR(255) NOT NULL,
    description TEXT,
    phone_number VARCHAR(20),
    email VARCHAR(255) UNIQUE NOT NULL,
    address TEXT NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    website VARCHAR(500),
    instagram VARCHAR(255),
    tiktok VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

/*
Membuat index pada kolom email untuk meningkatkan performa pencarian.
*/
CREATE INDEX idx_hosters_email ON hosters(email);

/*
Membuat index pada kolom created_at untuk meningkatkan performa pencarian.
*/
CREATE INDEX idx_hosters_created_at ON hosters(created_at);

/*
Fungsi untuk memperbarui kolom updated_at secara otomatis.
*/
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

/*
Trigger untuk memanggil fungsi update updated_at sebelum update.
*/
CREATE TRIGGER update_hosters_updated_at BEFORE UPDATE ON hosters FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();