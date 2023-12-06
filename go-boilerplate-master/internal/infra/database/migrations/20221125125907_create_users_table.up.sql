CREATE TABLE IF NOT EXISTS users
(
    id           SERIAL PRIMARY KEY,
    name         TEXT,
    email        TEXT,
    password     TEXT,
    created_date TIMESTAMP,
    updated_date TIMESTAMP,
    deleted_date TIMESTAMP NULL
);
