-- Создание таблицы парковочных мест
CREATE TABLE IF NOT EXISTS parking_spots (
    id VARCHAR(50) PRIMARY KEY,
    zone_id VARCHAR(50) NOT NULL REFERENCES parking_zones(id) ON DELETE CASCADE,
    spot_number VARCHAR(10) NOT NULL,
    status VARCHAR(20) DEFAULT 'free', -- free, occupied, reserved
    price_per_hour DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Уникальность номера места в пределах одной зоны
    UNIQUE(zone_id, spot_number)
);

-- Индексы для быстрого поиска
/*

CREATE INDEX IF NOT EXISTS idx_spots_zone_id ON parking_spots(zone_id);
CREATE INDEX IF NOT EXISTS idx_spots_status ON parking_spots(status);
CREATE INDEX IF NOT EXISTS idx_spots_zone_status ON parking_spots(zone_id, status);

-- Добавляем места для Центральной зоны
INSERT INTO parking_spots (id, zone_id, spot_number, status, price_per_hour) VALUES
('spot_center_A1', 'zone_center', 'A1', 'free', 500.00),
('spot_center_A2', 'zone_center', 'A2', 'free', 500.00),
('spot_center_A3', 'zone_center', 'A3', 'occupied', 500.00),
('spot_center_A4', 'zone_center', 'A4', 'free', 500.00),
('spot_center_B1', 'zone_center', 'B1', 'free', 700.00),
('spot_center_B2', 'zone_center', 'B2', 'reserved', 700.00),
('spot_center_B3', 'zone_center', 'B3', 'free', 700.00),
('spot_center_B4', 'zone_center', 'B4', 'free', 700.00)
ON CONFLICT (id) DO NOTHING;

-- Добавляем места для Северной зоны
INSERT INTO parking_spots (id, zone_id, spot_number, status, price_per_hour) VALUES
('spot_north_1', 'zone_north', '1', 'free', 400.00),
('spot_north_2', 'zone_north', '2', 'free', 400.00),
('spot_north_3', 'zone_north', '3', 'occupied', 400.00),
('spot_north_4', 'zone_north', '4', 'free', 400.00),
('spot_north_5', 'zone_north', '5', 'free', 400.00)
ON CONFLICT (id) DO NOTHING;

-- Добавляем места для Южной зоны
INSERT INTO parking_spots (id, zone_id, spot_number, status, price_per_hour) VALUES
('spot_south_1', 'zone_south', '1', 'free', 300.00),
('spot_south_2', 'zone_south', '2', 'free', 300.00),
('spot_south_3', 'zone_south', '3', 'free', 300.00),
('spot_south_4', 'zone_south', '4', 'reserved', 300.00)
ON CONFLICT (id) DO NOTHING;

*/