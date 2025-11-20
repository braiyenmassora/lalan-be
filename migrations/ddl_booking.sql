-- Tabel utama booking dengan referensi ke hoster
CREATE TABLE booking (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR NOT NULL,
    hoster_id UUID NOT NULL REFERENCES hosters(id),
    locked_until TIMESTAMP NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    total_days INTEGER NOT NULL,
    delivery_type VARCHAR NOT NULL,
    rental INTEGER NOT NULL,
    deposit INTEGER NOT NULL,
    discount INTEGER NOT NULL,
    total INTEGER NOT NULL,
    outstanding INTEGER NOT NULL,
    user_id UUID NOT NULL REFERENCES customer(id),
    identity_id UUID REFERENCES identity(id),
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
