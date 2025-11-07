-- Membuat tabel categories
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- ID unik dengan UUID
    name VARCHAR(255) NOT NULL UNIQUE, -- Nama kategori, unik
    description TEXT, -- Deskripsi kategori
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(), -- Waktu pembuatan
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() -- Waktu pembaruan
);

-- Index untuk performa pencarian nama
CREATE INDEX idx_categories_name ON categories(name);
-- Index untuk performa pencarian waktu pembuatan
CREATE INDEX idx_categories_created_at ON categories(created_at);

-- Fungsi untuk update kolom updated_at
CREATE OR REPLACE FUNCTION update_categories_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk auto-update updated_at
CREATE TRIGGER update_categories_updated_at
BEFORE UPDATE ON categories
FOR EACH ROW
EXECUTE FUNCTION update_categories_updated_at_column();