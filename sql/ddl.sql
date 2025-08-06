CREATE TABLE payments (
  id UUID PRIMARY KEY,
  amount NUMERIC(10, 2) NOT NULL,
  processor VARCHAR(20) NOT NULL CHECK (processor IN ('default', 'fallback')),
  requested_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

-- Índices otimizados para performance
CREATE INDEX idx_payments_processor_requested_at ON payments(processor, requested_at);
CREATE INDEX idx_payments_requested_at ON payments(requested_at);
CREATE INDEX idx_payments_processor ON payments(processor);

-- Otimizações do PostgreSQL
ALTER TABLE payments SET (fillfactor = 90);
ALTER TABLE payments ALTER COLUMN amount SET STORAGE MAIN;
ALTER TABLE payments ALTER COLUMN processor SET STORAGE MAIN;