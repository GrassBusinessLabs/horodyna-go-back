CREATE TABLE invoice_cancellations
(
  invoice_id    VARCHAR(255) PRIMARY KEY,
  status        VARCHAR(255) NOT NULL,
  amount        NUMERIC(12,2),
  approval_code VARCHAR(255),
  rrn           VARCHAR(255),
  created_date  TIMESTAMP NOT NULL DEFAULT timezone('UTC'::text, now()),
  updated_date  TIMESTAMP NOT NULL DEFAULT timezone('UTC'::text, now()),
  CONSTRAINT fk_invoice_id FOREIGN KEY (invoice_id) REFERENCES invoices(invoice_id) ON DELETE CASCADE
);