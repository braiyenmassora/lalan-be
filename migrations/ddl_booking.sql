-- Tabel utama booking dengan referensi ke hoster
CREATE TABLE booking (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hoster_id UUID NOT NULL REFERENCES hoster(id),
    locked_until TIMESTAMP NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    total_days INTEGER NOT NULL,
    delivery_type VARCHAR NOT NULL,
    rental INTEGER NOT NULL,
    deposit INTEGER NOT NULL,
    discount INTEGER NOT NULL,
    total INTEGER NOT NULL,
    outstanding INTEGER NOT NULL,
    user_id UUID NOT NULL REFERENCES customer(id),
    identity_id UUID REFERENCES identity(id),
    status VARCHAR NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Tabel item dalam booking
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

-- Tabel snapshot customer untuk booking
CREATE TABLE booking_customer (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL REFERENCES booking(id) ON DELETE CASCADE,
    name VARCHAR NOT NULL,
    phone VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    address TEXT NOT NULL,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index untuk filter cepat berdasarkan hoster
CREATE INDEX idx_booking_hoster_id
    ON booking(hoster_id);

-- Index lainnya
CREATE INDEX idx_booking_user_id
    ON booking(user_id);

CREATE INDEX idx_booking_start_date
    ON booking(start_date);

CREATE INDEX idx_booking_end_date
    ON booking(end_date);

CREATE INDEX idx_booking_created_at
    ON booking(created_at);

-- Index untuk booking_item
CREATE INDEX idx_booking_item_booking_id
    ON booking_item(booking_id);

CREATE INDEX idx_booking_item_item_id
    ON booking_item(item_id);

-- Index untuk booking_customer
CREATE INDEX idx_booking_customer_booking_id
    ON booking_customer(booking_id);
