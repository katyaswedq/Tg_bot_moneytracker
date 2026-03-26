CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    telegram_id BIGINT UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    user_name VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);