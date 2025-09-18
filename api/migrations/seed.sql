INSERT INTO users(email, password_hash, role) VALUES
('test@example.com', 'hashed_pw', 'user');

INSERT INTO players(user_id, display_name, rating) VALUES
(1, 'Alice', 1600),
(1, 'Bob', 1500);
