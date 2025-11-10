/*
Membuat tabel untuk menyimpan data item dengan relasi ke kategori dan user.
Menghasilkan struktur tabel dengan kolom detail item, harga, stok, dan foreign key.
*/
CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    photos JSONB,
    stock INTEGER NOT NULL DEFAULT 0,
    pickup_type VARCHAR(50) NOT NULL CHECK (pickup_type IN ('pickup', 'delivery')),
    price_per_day INTEGER NOT NULL,
    deposit INTEGER NOT NULL DEFAULT 0,
    discount INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    category_id UUID NOT NULL,
    user_id UUID NOT NULL,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES hosters(id) ON DELETE CASCADE
);

/*
Membuat index pada kolom name.
Meningkatkan performa query pencarian berdasarkan nama item.
*/
CREATE INDEX idx_items_name ON items(name);

/*
Membuat index pada kolom user_id.
Meningkatkan performa query filter berdasarkan user.
*/
CREATE INDEX idx_items_user_id ON items(user_id);

/*
Membuat index pada kolom category_id.
Meningkatkan performa query filter berdasarkan kategori.
*/
CREATE INDEX idx_items_category_id ON items(category_id);

/*
Membuat index pada kolom created_at.
Meningkatkan performa query pengurutan berdasarkan waktu pembuatan.
*/
CREATE INDEX idx_items_created_at ON items(created_at);

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
CREATE TRIGGER update_items_updated_at
BEFORE UPDATE ON items
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
