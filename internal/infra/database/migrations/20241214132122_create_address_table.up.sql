CREATE TABLE IF NOT EXISTS addresses
(
    id           SERIAL PRIMARY KEY,
    user_id      INTEGER NOT NULL,
    street       TEXT NOT NULL,
    city         TEXT NOT NULL,
    state        TEXT NOT NULL,
    zip_code     TEXT NOT NULL,
    created_date TIMESTAMP,
    updated_date TIMESTAMP,
    deleted_date TIMESTAMP NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
