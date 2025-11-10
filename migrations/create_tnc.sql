/*
Membuat tabel untuk menyimpan syarat dan ketentuan per user.
Menghasilkan struktur tabel dengan kolom user_id, description sebagai JSON, dan timestamp.
*/
CREATE TABLE terms_and_conditions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    description JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES hosters(id) ON DELETE CASCADE,
    UNIQUE (user_id)
);

/*
Membuat index pada kolom user_id.
Meningkatkan performa query filter berdasarkan user.
*/
CREATE INDEX idx_terms_and_conditions_user_id ON terms_and_conditions(user_id);

/*
Membuat index pada kolom created_at.
Meningkatkan performa query pengurutan berdasarkan waktu pembuatan.
*/
CREATE INDEX idx_terms_and_conditions_created_at ON terms_and_conditions(created_at);

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
CREATE TRIGGER update_terms_and_conditions_updated_at
BEFORE UPDATE ON terms_and_conditions
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();