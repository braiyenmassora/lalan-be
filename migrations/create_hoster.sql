/*
Membuat tabel untuk menyimpan data hoster.
Menghasilkan struktur tabel dengan kolom pribadi, toko, kontak, dan timestamp.
*/
CREATE TABLE hosters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
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
Membuat index pada kolom email.
Meningkatkan performa query pencarian berdasarkan email.
*/
CREATE INDEX idx_hosters_email ON hosters(email);

/*
Membuat index pada kolom created_at.
Meningkatkan performa query pengurutan berdasarkan waktu pembuatan.
*/
CREATE INDEX idx_hosters_created_at ON hosters(created_at);

/*
Membuat fungsi untuk update otomatis kolom updated_at.
Mengembalikan baris yang diperbarui dengan timestamp baru.
*/
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

/*
Membuat trigger untuk memanggil fungsi update sebelum perubahan.
Memastikan kolom updated_at selalu diperbarui saat update.
*/
CREATE TRIGGER update_hosters_updated_at
BEFORE UPDATE ON hosters
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();