/*
Membuat tabel untuk menyimpan data kategori.
Menghasilkan struktur tabel dengan kolom id, name, description, dan timestamp.
*/
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

/*
Membuat index pada kolom name.
Meningkatkan performa query pencarian berdasarkan nama.
*/
CREATE INDEX idx_categories_name ON categories(name);

/*
Membuat index pada kolom created_at.
Meningkatkan performa query pengurutan berdasarkan waktu pembuatan.
*/
CREATE INDEX idx_categories_created_at ON categories(created_at);

/*
Membuat fungsi untuk update otomatis kolom updated_at.
Mengembalikan baris yang diperbarui dengan timestamp baru.
*/
CREATE OR REPLACE FUNCTION update_categories_updated_at_column()
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
CREATE TRIGGER update_categories_updated_at
BEFORE UPDATE ON categories
FOR EACH ROW
EXECUTE FUNCTION update_categories_updated_at_column();