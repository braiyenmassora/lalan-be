/*
Membuat tabel customers dengan kolom pribadi, kontak, dan timestamp.
Menyediakan struktur untuk menyimpan data pelanggan.
*/
    CREATE TABLE customers (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        full_name VARCHAR(255) NOT NULL,
        profile_photo VARCHAR(500),
        phone_number VARCHAR(20),
        email VARCHAR(255) UNIQUE NOT NULL,
        address TEXT,
        password_hash VARCHAR(255) NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );

/*
Menambahkan index pada kolom email di tabel customers.
Mempercepat pencarian berdasarkan email pelanggan.
*/
    CREATE INDEX idx_customers_email
        ON customers(email);

/*
Menambahkan index pada kolom created_at di tabel customers.
Mempercepat pengurutan berdasarkan waktu pembuatan pelanggan.
*/
    CREATE INDEX idx_customers_created_at
        ON customers(created_at);

/*
Membuat fungsi untuk memperbarui kolom updated_at secara otomatis.
Digunakan oleh trigger untuk menjaga timestamp pembaruan di tabel customers.
*/
    CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
    BEGIN
        NEW.updated_at = NOW();
        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

/*
Membuat trigger untuk memanggil fungsi update sebelum perubahan pada tabel customers.
Memastikan kolom updated_at selalu diperbarui saat update.
*/
    CREATE TRIGGER update_customers_updated_at
    BEFORE UPDATE ON customers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();