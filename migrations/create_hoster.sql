CREATE TABLE IF NOT EXISTS hoster (
    id VARCHAR(26) PRIMARY KEY,
    owner_name VARCHAR(255) NOT NULL,
    store_name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(50),
    email VARCHAR(255) UNIQUE NOT NULL,
    address TEXT NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
