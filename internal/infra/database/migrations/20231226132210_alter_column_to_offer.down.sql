ALTER TABLE offers
ALTER COLUMN price TYPE FLOAT4 USING price::FLOAT4;