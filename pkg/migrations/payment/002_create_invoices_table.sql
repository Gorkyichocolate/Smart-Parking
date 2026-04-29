-- Создание таблицы счетов/квитанций
CREATE TABLE IF NOT EXISTS invoices (
    id VARCHAR(50) PRIMARY KEY,
    payment_id VARCHAR(50) NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    user_id VARCHAR(50) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    pdf_url TEXT,
    issued_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для поиска по платежу
/*

CREATE INDEX IF NOT EXISTS idx_invoices_payment_id ON invoices(payment_id);
CREATE INDEX IF NOT EXISTS idx_invoices_user_id ON invoices(user_id);

-- Добавляем тестовый счет (привязывается к тестовому платежу)
-- Сначала создаем тестовый платеж, потом счет
INSERT INTO payments (id, booking_id, user_id, amount, status, payment_method) 
VALUES (
    'pay_test_001',
    'book_test_001',
    'user_test_001',
    1000.00,
    'completed',
    'card'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO invoices (id, payment_id, user_id, amount, pdf_url) 
VALUES (
    'inv_test_001',
    'pay_test_001',
    'user_test_001',
    1000.00,
    '/invoices/inv_test_001.pdf'
) ON CONFLICT (id) DO NOTHING;

*/