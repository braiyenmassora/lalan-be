/*
Membuat tabel categories untuk menyimpan data kategori.
*/
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

/*
Membuat index pada kolom name untuk meningkatkan performa pencarian.
*/
CREATE INDEX idx_categories_name ON categories(name);

/*
Membuat index pada kolom created_at untuk meningkatkan performa pencarian.
*/
CREATE INDEX idx_categories_created_at ON categories(created_at);

/*
Fungsi untuk memperbarui kolom updated_at secara otomatis.
*/
CREATE OR REPLACE FUNCTION update_categories_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

/*
Trigger untuk memanggil fungsi update updated_at sebelum update.
*/
CREATE TRIGGER update_categories_updated_at
BEFORE UPDATE ON categories
FOR EACH ROW
EXECUTE FUNCTION update_categories_updated_at_column();