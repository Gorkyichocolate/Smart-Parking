CREATE TABLE IF NOT EXISTS invoices (
                                        id VARCHAR(255) PRIMARY KEY,
                                        payment_id VARCHAR(255) NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
                                        user_id VARCHAR(255) NOT NULL,
                                        amount DECIMAL(10,2) NOT NULL,
                                        pdf_url VARCHAR(500) NOT NULL,
                                        issued_at TIMESTAMP NOT NULL
);

-- Индексы для invoices
CREATE INDEX IF NOT EXISTS idx_invoices_payment_id ON invoices(payment_id);
CREATE INDEX IF NOT EXISTS idx_invoices_user_id ON invoices(user_id);
CREATE INDEX IF NOT EXISTS idx_invoices_issued_at ON invoices(issued_at);