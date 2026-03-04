-- Удаляем таблицы если существуют
DROP TABLE IF EXISTS workouts;
DROP TABLE IF EXISTS memberships;
DROP TABLE IF EXISTS users;

-- Создаем таблицу пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создаем таблицу абонементов
CREATE TABLE memberships (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('standard', 'premium', 'unlimited')),
    price DECIMAL(10,2) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создаем таблицу тренировок
CREATE TABLE workouts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    duration_minutes INTEGER NOT NULL,
    calories_burned INTEGER,
    notes TEXT,
    workout_date DATE DEFAULT CURRENT_DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Начальные данные
INSERT INTO users (name, email, phone) VALUES
    ('Алексей Петров', 'alex@energy.ru', '+79161234567'),
    ('Мария Сидорова', 'maria@energy.ru', '+79167654321'),
    ('Иван Иванов', 'ivan@energy.ru', '+79161112233');

INSERT INTO memberships (user_id, type, price, start_date, end_date, is_active) VALUES
    (1, 'premium', 2999.99, '2024-01-01', '2024-12-31', true),
    (2, 'standard', 1999.99, '2024-02-01', '2024-11-30', true),
    (3, 'unlimited', 4999.99, '2024-03-01', '2025-02-28', true);

INSERT INTO workouts (user_id, type, duration_minutes, calories_burned, notes) VALUES
    (1, 'Йога', 60, 300, 'Утренняя практика'),
    (1, 'Кардио', 45, 500, 'Беговая дорожка'),
    (2, 'Силовая', 90, 600, 'Тренировка с весом'),
    (3, 'Плавание', 60, 400, 'Бассеин 50м');