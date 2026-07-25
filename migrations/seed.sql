INSERT INTO users (id, full_name, email, password)
VALUES
    (1, 'Justin Bieber', 'justin@example.com', 'password'),
    (2, 'Taylor Swift', 'taylor@example.com', 'password'),
    (3, 'Ed Sheeran', 'ed@example.com', 'password');

INSERT INTO wallets (id, user_id, balance)
VALUES
    (1, 1, 10000),
    (2, 2, 5000),
    (3, 3, 2500);

SELECT setval(pg_get_serial_sequence('users', 'id'), (SELECT COALESCE(MAX(id), 1) FROM users));
SELECT setval(pg_get_serial_sequence('wallets', 'id'), (SELECT COALESCE(MAX(id), 1) FROM wallets));
