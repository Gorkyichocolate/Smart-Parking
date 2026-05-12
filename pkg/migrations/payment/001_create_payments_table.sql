CREATE TABLE IF NOT EXISTS payments (
                                        id VARCHAR(255) PRIMARY KEY,
                                        booking_id VARCHAR(255) NOT NULL,
                                        user_id VARCHAR(255) NOT NULL,
                                        amount DECIMAL(10,2) NOT NULL,
                                        status VARCHAR(50) NOT NULL DEFAULT 'pending',
                                        payment_method VARCHAR(50) NOT NULL,
                                        created_at TIMESTAMP NOT NULL,
                                        updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_payments_booking_id ON payments(booking_id);
CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_status ON payments(status);