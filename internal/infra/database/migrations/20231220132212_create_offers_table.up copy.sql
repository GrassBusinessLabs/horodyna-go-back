CREATE TABLE IF NOT EXISTS images
(
    Id  SERIAL PRIMARY KEY,
    Title TEXT NOT NULL,
    Entity TEXT NOT NULL,
    EntityId TEXT NOT NULL,
    Created_date TIMESTAMP,
    Updated_date TIMESTAMP,
    Deleted_date TIMESTAMP NULL
);