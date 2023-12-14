CREATE TABLE IF NOT EXISTS farms
(
    id           SERIAL PRIMARY KEY,
    name         TEXT NOT NULL,
    city         TEXT,
    address      TEXT,
    latitude     FLOAT8,
    longitude    FLOAT8,
    user_id      INTEGER NOT NULL,
    created_date TIMESTAMP,
    updated_date TIMESTAMP,
    deleted_date TIMESTAMP NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);