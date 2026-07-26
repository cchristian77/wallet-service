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

INSERT INTO transactions (id, transaction_id, status)
VALUES
    (1, 'TOPUP-WALLET-1', 'SUCCESS'),
    (2, 'TOPUP-WALLET-2', 'SUCCESS'),
    (3, 'TOPUP-WALLET-3', 'SUCCESS');

INSERT INTO transaction_ledgers (id, transaction_id, wallet_id, ledger, reference, amount)
VALUES
    (1, 1, 1, 'CREDIT', 'TOP_UP', 10000),
    (2, 2, 2, 'CREDIT', 'TOP_UP', 5000),
    (3, 3, 3, 'CREDIT', 'TOP_UP', 2500);

SELECT setval(pg_get_serial_sequence('users', 'id'), (SELECT COALESCE(MAX(id), 1) FROM users));
SELECT setval(pg_get_serial_sequence('wallets', 'id'), (SELECT COALESCE(MAX(id), 1) FROM wallets));
SELECT setval(pg_get_serial_sequence('transactions', 'id'), (SELECT COALESCE(MAX(id), 1) FROM transactions));
SELECT setval(pg_get_serial_sequence('transaction_ledgers', 'id'), (SELECT COALESCE(MAX(id), 1) FROM transaction_ledgers))
