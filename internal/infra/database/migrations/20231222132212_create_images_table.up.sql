CREATE TABLE IF NOT EXISTS images
(
    id  SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    entity TEXT NOT NULL,
    entity_id INTEGER NOT NULL,
    created_date TIMESTAMP,
    updated_date TIMESTAMP,
    deleted_date TIMESTAMP NULL
);