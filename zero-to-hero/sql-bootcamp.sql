DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id  SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    email TEXT UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    amount INT NOT NULL,
    status TEXT DEFAULT 'new'
);

INSERT INTO users (username, email) VALUES
    ('maxim', 'maksim@gmail.com'),
    ('leonid', 'leonid@gmail.com'),
    ('katya', 'katya@gmail.com');

INSERT INTO orders (user_id, status, amount)
VALUES 
    (1, 'pending', 1000),
    (2, 'pending', 500),
    (2, 'new', 1200),
    (2, 'pending', 300),
    (1, 'new', 123);

SELECT u.username,
    COUNT(o.id) as orders_count,
    COALESCE(SUM(o.amount), 0) as total_spent
FROM users u
LEFT JOIN orders o ON u.id = o.user_id
GROUP BY u.username;