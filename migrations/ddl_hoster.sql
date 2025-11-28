/*
Membuat tabel hosters dengan kolom pribadi, toko, kontak, dan timestamp.
Menyediakan struktur untuk menyimpan data penyedia layanan hosting.
*/
    CREATE TABLE hoster (
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
Menambahkan index pada kolom email di tabel hosters.
Mempercepat pencarian berdasarkan email hoster.
*/
    CREATE INDEX idx_hoster_email
        ON hoster(email);

/*
Menambahkan index pada kolom created_at di tabel hosters.
Mempercepat pengurutan berdasarkan waktu pembuatan hoster.
*/
    CREATE INDEX idx_hoster_created_at
        ON hoster(created_at);

/*
Membuat fungsi untuk memperbarui kolom updated_at secara otomatis.
Digunakan oleh trigger untuk menjaga timestamp pembaruan di tabel hosters.
*/
    CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
    BEGIN
        NEW.updated_at = NOW();
        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

/*
Membuat trigger untuk memanggil fungsi update sebelum perubahan pada tabel hosters.
Memastikan kolom updated_at selalu diperbarui saat update.
*/
    CREATE TRIGGER update_hoster_updated_at
    BEFORE UPDATE ON hoster
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();