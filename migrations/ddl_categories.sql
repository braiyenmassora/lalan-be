/*
Membuat tabel categories dengan kolom id, name, description, dan timestamp.
Menyediakan struktur untuk mengelompokkan item berdasarkan kategori.
*/
    CREATE TABLE categories (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        name VARCHAR(255) NOT NULL UNIQUE,
        description TEXT,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );

/*
Menambahkan index pada kolom name di tabel categories.
Mempercepat pencarian berdasarkan nama kategori.
*/
    CREATE INDEX idx_categories_name
        ON categories(name);

/*
Menambahkan index pada kolom created_at di tabel categories.
Mempercepat pengurutan berdasarkan waktu pembuatan kategori.
*/
    CREATE INDEX idx_categories_created_at
        ON categories(created_at);

/*
Membuat fungsi untuk memperbarui kolom updated_at secara otomatis di tabel categories.
Digunakan oleh trigger untuk menjaga timestamp pembaruan.
*/
    CREATE OR REPLACE FUNCTION update_categories_updated_at_column()
    RETURNS TRIGGER AS $$
    BEGIN
        NEW.updated_at = NOW();
        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

/*
Membuat trigger untuk memanggil fungsi update sebelum perubahan pada tabel categories.
Memastikan kolom updated_at selalu diperbarui saat update.
*/
    CREATE TRIGGER update_categories_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW
    EXECUTE FUNCTION update_categories_updated_at_column();