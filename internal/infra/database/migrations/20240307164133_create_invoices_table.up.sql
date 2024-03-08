CREATE TYPE invoice_status AS ENUM
(
    'created',
    'processing',
    'hold',
    'success',
    'failure',
    'reversed',
    'expired'
);

CREATE TABLE invoices
(
    invoice_id     VARCHAR(255) PRIMARY KEY,
    status         invoice_status NOT NULL DEFAULT 'created',
    failure_reason VARCHAR(255),
    err_code       VARCHAR(255),
    amount         NUMERIC(12,2),
    final_amount   NUMERIC(12,2),
    created_date   TIMESTAMP NOT NULL DEFAULT timezone('UTC'::text, now()),
    updated_date   TIMESTAMP NOT NULL DEFAULT timezone('UTC'::text, now())
);