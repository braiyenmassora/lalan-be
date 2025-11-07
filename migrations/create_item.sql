-- Membuat ekstensi UUID jika belum ada
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Membuat tabel items
CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), -- ID unik dengan UUID
    name VARCHAR(255) NOT NULL, -- Nama item
    description TEXT, -- Deskripsi item
    photos JSONB, -- Array URL foto dalam format JSON
    stock INTEGER NOT NULL DEFAULT 0, -- Jumlah stok
    pickup_type VARCHAR(50) NOT NULL CHECK (pickup_type IN ('pickup', 'delivery')), -- Tipe pengambilan
    price_per_day INTEGER NOT NULL, -- Harga sewa per hari
    deposit INTEGER NOT NULL DEFAULT 0, -- Jumlah deposit
    discount INTEGER DEFAULT 0, -- Diskon persentase
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(), -- Waktu pembuatan
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(), -- Waktu pembaruan
    category_id UUID NOT NULL, -- ID kategori (foreign key)
    user_id UUID NOT NULL, -- ID pengguna/hoster (foreign key)
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES hosters(id) ON DELETE CASCADE
);

-- Index untuk performa pencarian nama item
CREATE INDEX idx_items_name ON items(name);
-- Index untuk performa pencarian user_id
CREATE INDEX idx_items_user_id ON items(user_id);
-- Index untuk performa pencarian category_id
CREATE INDEX idx_items_category_id ON items(category_id);
-- Index untuk performa pencarian waktu pembuatan
CREATE INDEX idx_items_created_at ON items(created_at);

-- Fungsi untuk update kolom updated_at (menggunakan fungsi yang sama seperti di hosters)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger untuk auto-update updated_at
CREATE TRIGGER update_items_updated_at BEFORE UPDATE ON items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
