/*
Membuat tabel identity dengan kolom untuk verifikasi identitas customer.
Menyediakan struktur untuk menyimpan data KTP, status verifikasi, dan foreign key ke customer.
*/
CREATE TABLE identity (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    ktp_url VARCHAR(255) NOT NULL,  -- URL atau path file KTP
    verified BOOLEAN DEFAULT FALSE,
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'approved', 'rejected')),
    rejected_reason TEXT,
    verified_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES customer(id) ON DELETE CASCADE
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