/*
Membuat tabel admin dengan kolom pribadi, kontak, dan timestamp.
Menyediakan struktur dasar untuk menyimpan data administrator sistem.
*/
    CREATE TABLE admin (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        full_name VARCHAR(255) NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );

/*
Menambahkan index pada kolom email di tabel admin.
Mempercepat pencarian berdasarkan email dan menegakkan keunikan.
*/
    CREATE INDEX idx_admin_email
        ON admin(email);

/*
Menambahkan index pada kolom created_at di tabel admin.
Mempercepat pengurutan berdasarkan waktu pembuatan.
*/
    CREATE INDEX idx_admin_created_at
        ON admin(created_at);

/*
Membuat fungsi untuk memperbarui kolom updated_at secara otomatis.
Digunakan oleh trigger untuk menjaga timestamp pembaruan.
*/
    CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
    BEGIN
        NEW.updated_at = NOW();
        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

/*
Membuat trigger untuk memanggil fungsi update sebelum perubahan pada tabel admin.
Memastikan kolom updated_at selalu diperbarui saat update.
*/
    CREATE TRIGGER update_admin_updated_at
    BEFORE UPDATE ON admin
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();