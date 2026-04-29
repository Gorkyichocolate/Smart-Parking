-- Создание таблицы бронирований
CREATE TABLE IF NOT EXISTS bookings (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL, -- ID из auth_db
    spot_id VARCHAR(50) NOT NULL REFERENCES parking_spots(id) ON DELETE CASCADE,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- pending, active, completed, cancelled
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Проверка: время начала раньше времени конца
    CHECK (start_time < end_time)
);

-- Индексы для быстрого поиска
/*

CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_spot_id ON bookings(spot_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status);
CREATE INDEX IF NOT EXISTS idx_bookings_time_range ON bookings(start_time, end_time);

-- Индекс для проверки пересечения бронирований
CREATE INDEX IF NOT EXISTS idx_bookings_spot_time ON bookings(spot_id, start_time, end_time);

-- Функция для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Триггер для обновления updated_at
CREATE TRIGGER update_bookings_updated_at
    BEFORE UPDATE ON bookings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Добавляем тестовое бронирование (для демонстрации)
INSERT INTO bookings (id, user_id, spot_id, start_time, end_time, total_price, status) 
VALUES (
    'book_test_001',
    'user_test_001',
    'spot_center_A1',
    NOW() + INTERVAL '1 day',
    NOW() + INTERVAL '1 day' + INTERVAL '2 hours',
    1000.00,
    'pending'
) ON CONFLICT (id) DO NOTHING;

*/