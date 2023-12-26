CREATE TABLE IF NOT EXISTS orders
(
    id             SERIAL PRIMARY KEY,
    comment        TEXT,
    user_id        INTEGER NOT NULL,
    address_id     INTEGER NOT NULL,
    products_price FLOAT8 NOT NULL,
    shipping_price FLOAT8,
    total_price    FLOAT8 NOT NULL,
    status         TEXT NOT NULL,
    created_date   TIMESTAMP,
    updated_date   TIMESTAMP,
    deleted_date   TIMESTAMP NULL,
    CONSTRAINT fk_addresses_id FOREIGN KEY (address_id) REFERENCES addresses(id) ON DELETE CASCADE,    
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);