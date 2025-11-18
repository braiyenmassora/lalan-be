-- Tabel utama booking
CREATE TABLE booking (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR NOT NULL,
    status VARCHAR NOT NULL,
    locked_until TIMESTAMP NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    total_days INTEGER NOT NULL,
    delivery_type VARCHAR NOT NULL,
    rental INTEGER NOT NULL,
    deposit INTEGER NOT NULL,
    delivery INTEGER NOT NULL,
    discount INTEGER NOT NULL,
    total INTEGER NOT NULL,
    outstanding INTEGER NOT NULL,
    user_id UUID NOT NULL REFERENCES customer(id),
    identity_id UUID REFERENCES identity(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

/*
Menambahkan index pada kolom user_id di tabel booking.
Mempercepat query filter berdasarkan customer pemilik booking.
*/
CREATE INDEX idx_booking_user_id
    ON booking(user_id);

/*
Menambahkan index pada kolom status di tabel booking.
Mempercepat query filter berdasarkan status booking.
*/
CREATE INDEX idx_booking_status
    ON booking(status);

/*
Menambahkan index pada kolom start_date di tabel booking.
Mempercepat query range berdasarkan tanggal mulai rental.
*/
CREATE INDEX idx_booking_start_date
    ON booking(start_date);

/*
Menambahkan index pada kolom end_date di tabel booking.
Mempercepat query range berdasarkan tanggal selesai rental.
*/
CREATE INDEX idx_booking_end_date
    ON booking(end_date);

/*
Menambahkan index pada kolom created_at di tabel booking.
Mempercepat pengurutan berdasarkan waktu pembuatan booking.
*/
CREATE INDEX idx_booking_created_at
    ON booking(created_at);

/*
Membuat tabel booking_item untuk menyimpan detail item dalam booking.
Relasi one-to-many dengan booking, menyediakan struktur untuk subtotal harga.
*/
CREATE TABLE booking_item (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES booking(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES item(id),
    name VARCHAR NOT NULL,
    quantity INTEGER NOT NULL,
    price_per_day INTEGER NOT NULL,
    deposit_per_unit INTEGER NOT NULL,
    subtotal_rental INTEGER NOT NULL,
    subtotal_deposit INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

/*
Menambahkan index pada kolom booking_id di tabel booking_item.
Mempercepat query join dengan tabel booking.
*/
CREATE INDEX idx_booking_item_booking_id
    ON booking_item(booking_id);

/*
Menambahkan index pada kolom item_id di tabel booking_item.
Mempercepat query join dengan tabel item.
*/
CREATE INDEX idx_booking_item_item_id
    ON booking_item(item_id);

/*
Membuat tabel booking_customer sebagai snapshot data customer saat booking.
Relasi one-to-one dengan booking, menyimpan info pengiriman.
*/
CREATE TABLE booking_customer (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES booking(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    phone VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    delivery_address VARCHAR NOT NULL,
    notes VARCHAR,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

/*
Menambahkan index pada kolom booking_id di tabel booking_customer.
Mempercepat query join dengan tabel booking.
*/
CREATE INDEX idx_booking_customer_booking_id
    ON booking_customer(booking_id);

/*
Membuat tabel booking_identity sebagai snapshot status identity saat booking.
Relasi one-to-one dengan booking, menyimpan info verifikasi KTP.
*/
CREATE TABLE booking_identity (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES booking(id) ON DELETE CASCADE,
    uploaded BOOLEAN NOT NULL,
    status VARCHAR NOT NULL,
    rejection_reason VARCHAR,
    reupload_allowed BOOLEAN NOT NULL,
    estimated_time VARCHAR NOT NULL,
    status_check_url VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

/*
Menambahkan index pada kolom booking_id di tabel booking_identity.
Mempercepat query join dengan tabel booking.
*/
CREATE INDEX idx_booking_identity_booking_id
    ON booking_identity(booking_id);

/*
Membuat fungsi untuk memperbarui kolom updated_at secara otomatis.
Digunakan oleh trigger untuk menjaga timestamp pembaruan di tabel booking.
*/
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

/*
Membuat trigger untuk memanggil fungsi update sebelum perubahan pada tabel booking.
Memastikan kolom updated_at selalu diperbarui saat update.
*/
CREATE TRIGGER update_booking_updated_at
    BEFORE UPDATE ON booking
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

/*
Membuat trigger untuk memanggil fungsi update sebelum perubahan pada tabel booking_item.
Memastikan kolom updated_at selalu diperbarui saat update.
*/
CREATE TRIGGER update_booking_item_updated_at
    BEFORE UPDATE ON booking_item
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

/*
Membuat trigger untuk memanggil fungsi update sebelum perubahan pada tabel booking_customer.
Memastikan kolom updated_at selalu diperbarui saat update.
*/
CREATE TRIGGER update_booking_customer_updated_at
    BEFORE UPDATE ON booking_customer
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

/*
Membuat trigger untuk memanggil fungsi update sebelum perubahan pada tabel booking_identity.
Memastikan kolom updated_at selalu diperbarui saat update.
*/
CREATE TRIGGER update_booking_identity_updated_at
    BEFORE UPDATE ON booking_identity
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();