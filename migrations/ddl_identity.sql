/*
Membuat tabel identity dengan kolom untuk verifikasi identitas customer.
Menyediakan struktur untuk menyimpan data KTP, status verifikasi, dan foreign key ke customer.
*/
CREATE TABLE identity (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES customer(id) ON DELETE CASCADE,
    ktp_url VARCHAR NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    status VARCHAR NOT NULL DEFAULT 'pending',
    reason TEXT,
    verified_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

/*
Menambahkan index pada kolom user_id di tabel identity.
Mempercepat query filter berdasarkan customer pemilik identitas.
*/
CREATE INDEX idx_identity_user_id
    ON identity(user_id);

/*
Menambahkan index pada kolom status di tabel identity.
Mempercepat query filter berdasarkan status verifikasi.
*/
CREATE INDEX idx_identity_status
    ON identity(status);

/*
Menambahkan index pada kolom created_at di tabel identity.
Mempercepat pengurutan berdasarkan waktu pembuatan identitas.
*/
CREATE INDEX idx_identity_created_at
    ON identity(created_at);

/*
Membuat fungsi untuk memperbarui kolom updated_at secara otomatis.
Digunakan oleh trigger untuk menjaga timestamp pembaruan di tabel identity.
*/
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

/*
Membuat trigger untuk memanggil fungsi update sebelum perubahan pada tabel identity.
Memastikan kolom updated_at selalu diperbarui saat update.
*/
CREATE TRIGGER update_identity_updated_at
    BEFORE UPDATE ON identity
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();