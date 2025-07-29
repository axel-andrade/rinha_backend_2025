CREATE TABLE payments (
  id UUID PRIMARY KEY,
  amount NUMERIC(10, 2) NOT NULL,
  processor VARCHAR(20) NOT NULL CHECK (processor IN ('default', 'fallback')),
  requested_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);