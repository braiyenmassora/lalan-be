/*
Membuat tabel tnc dengan kolom user_id, description sebagai JSON, dan timestamp.
Menyediakan struktur untuk menyimpan syarat dan ketentuan per hoster.
*/
    CREATE TABLE tnc (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        user_id UUID NOT NULL,
        description JSONB NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
        FOREIGN KEY (user_id) REFERENCES hosters(id) ON DELETE CASCADE,
        UNIQUE (user_id)
    );

/*
Menambahkan index pada kolom user_id di tabel tnc.
Mempercepat query filter berdasarkan hoster pemilik syarat dan ketentuan.
*/
    CREATE INDEX idx_tnc_user_id
        ON tnc(user_id);

/*
Menambahkan index pada kolom created_at di tabel tnc.
Mempercepat pengurutan berdasarkan waktu pembuatan syarat dan ketentuan.
*/
    CREATE INDEX idx_tnc_created_at
        ON tnc(created_at);

/*
Membuat fungsi untuk memperbarui kolom updated_at secara otomatis.
Digunakan oleh trigger untuk menjaga timestamp pembaruan di tabel tnc.
*/
    CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
    BEGIN
        NEW.updated_at = NOW();
        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;

/*
Membuat trigger untuk memanggil fungsi update sebelum perubahan pada tabel tnc.
Memastikan kolom updated_at selalu diperbarui saat update.
*/
    CREATE TRIGGER update_tnc_updated_at
    BEFORE UPDATE ON tnc
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
