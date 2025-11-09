/*
Mengaktifkan ekstensi UUID untuk menghasilkan UUID.
*/
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

/*
Membuat tabel items untuk menyimpan data item dengan foreign key ke categories dan hosters.
*/
CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
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
Membuat index pada kolom name untuk meningkatkan performa pencarian.
*/
CREATE INDEX idx_items_name ON items(name);

/*
Membuat index pada kolom user_id untuk meningkatkan performa pencarian.
*/
CREATE INDEX idx_items_user_id ON items(user_id);

/*
Membuat index pada kolom category_id untuk meningkatkan performa pencarian.
*/
CREATE INDEX idx_items_category_id ON items(category_id);

/*
Membuat index pada kolom created_at untuk meningkatkan performa pencarian.
*/
CREATE INDEX idx_items_created_at ON items(created_at);

/*
Fungsi untuk memperbarui kolom updated_at secara otomatis.
*/
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

/*
Trigger untuk memanggil fungsi update updated_at sebelum update.
*/
CREATE TRIGGER update_items_updated_at BEFORE UPDATE ON items FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
