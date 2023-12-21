CREATE TABLE IF NOT EXISTS offers
(
    id           SERIAL PRIMARY KEY,
    title        TEXT NOT NULL,
    description  TEXT NOT NULL,
    category     TEXT NOT NULL,
    price        FLOAT4 NOT NULL,
    unit         TEXT NOT NULL,
    stock        INTEGER NOT NULL,
    cover        TEXT NOT NULL,
    user_id      INTEGER NOT NULL,
    farm_id      INTEGER NOT NULL,
    status       BOOLEAN NOT NULL,
    created_date TIMESTAMP,
    updated_date TIMESTAMP,
    deleted_date TIMESTAMP NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_farm_id FOREIGN KEY (farm_id) REFERENCES farms(id) ON DELETE CASCADE
);