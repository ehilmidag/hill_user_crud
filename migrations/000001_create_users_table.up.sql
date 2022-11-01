CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    email text UNIQUE NOT NULL,
    password text NOT NULL,
    created_at date NOT NULL,
    updated_at date NOT NULL
    );