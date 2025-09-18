CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT DEFAULT 'user',
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE players (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE SET NULL,
    display_name TEXT NOT NULL,
    rating INT NOT NULL DEFAULT 1500,
    created_at TIMESTAMP DEFAULT now()
);
