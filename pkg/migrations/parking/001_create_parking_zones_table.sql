-- Создание таблицы зон парковки
CREATE TABLE IF NOT EXISTS parking_zones (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    total_spots INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для поиска по имени зоны
/*

CREATE INDEX IF NOT EXISTS idx_zones_name ON parking_zones(name);

-- Добавляем тестовые зоны
INSERT INTO parking_zones (id, name, address, total_spots) VALUES
('zone_center', 'Центральная парковка', 'ул. Абая 15', 20),
('zone_north', 'Северная парковка', 'пр. Назарбаева 45', 15),
('zone_south', 'Южная парковка', 'ул. Пушкина 8', 12)
ON CONFLICT (id) DO NOTHING;

*/