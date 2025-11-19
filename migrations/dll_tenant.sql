/*
Membuat tabel tenant dengan kolom id, name, hoster_id, created_at, dan updated_at.
Menyediakan struktur untuk menyimpan data tenant yang berelasi one-to-one dengan hoster.
*/
CREATE TABLE tenant (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    hoster_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (hoster_id) REFERENCES hoster(id) ON DELETE CASCADE
);

/*
Menambahkan index pada kolom name di tabel tenant.
Mempercepat pencarian berdasarkan nama tenant.
*/
CREATE INDEX idx_tenant_name
    ON tenant(name);

/*
Menambahkan index pada kolom hoster_id di tabel tenant.
Mempercepat query filter berdasarkan hoster terkait tenant.
*/
CREATE INDEX idx_tenant_hoster_id
    ON tenant(hoster_id);

/*
Menambahkan index pada kolom created_at di tabel tenant.
Mempercepat pengurutan berdasarkan waktu pembuatan tenant.
*/
CREATE INDEX idx_tenant_created_at
    ON tenant(created_at);

/*
Membuat fungsi untuk memperbarui kolom updated_at secara otomatis.
Digunakan oleh trigger untuk menjaga timestamp pembaruan di tabel tenant.
*/
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

/*
Membuat trigger untuk memanggil fungsi update sebelum perubahan pada tabel tenant.
Memastikan kolom updated_at selalu diperbarui saat update.
*/
CREATE TRIGGER update_tenant_updated_at
BEFORE UPDATE ON tenant
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();