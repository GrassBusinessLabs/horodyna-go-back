CREATE TABLE IF NOT EXISTS addresses
(
    id           SERIAL PRIMARY KEY,
    user_id      INTEGER NOT NULL,
    title        TEXT NOT NULL,
    city         TEXT NOT NULL,
    country      TEXT NOT NULL,
    address      TEXT NOT NULL,
    lat          TEXT NOT NULL,
    lon          TEXT NOT NULL,
    created_date TIMESTAMP,
    updated_date TIMESTAMP,
    deleted_date TIMESTAMP NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
