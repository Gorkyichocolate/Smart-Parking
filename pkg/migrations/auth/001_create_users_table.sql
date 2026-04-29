-- Создание таблицы users
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(50) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'driver',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для быстрого поиска по email
/*
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Индекс для поиска по роли
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Добавляем администратора по умолчанию (пароль: admin123)
INSERT INTO users (id, email, password_hash, full_name, role) 
VALUES (
    'admin_001',
    'admin@smartparking.com',
    '8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918',
    'System Admin',
    'admin'
) ON CONFLICT (email) DO NOTHING;

-- Добавляем тестового пользователя (пароль: driver123)
INSERT INTO users (id, email, password_hash, full_name, role) 
VALUES (
    'user_test_001',
    'driver@example.com',
    '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92',
    'Test Driver',
    'driver'
) ON CONFLICT (email) DO NOTHING;
*/