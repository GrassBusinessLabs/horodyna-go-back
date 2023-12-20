CREATE TABLE IF NOT EXISTS order_items
(
    id           SERIAL PRIMARY KEY,
    amount       INTEGER NOT NULL,
    title        TEXT,
    price        FLOAT4 NOT NULL,
    total_price  FLOAT8 NOT NULL,
    order_id     INTEGER NOT NULL,
    offer_id     INTEGER NOT NULL,
    created_date TIMESTAMP,
    updated_date TIMESTAMP,
    deleted_date TIMESTAMP NULL,
    CONSTRAINT fk_order_id FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    CONSTRAINT fk_offer_id FOREIGN KEY (offer_id) REFERENCES offers(id) ON DELETE CASCADE
);