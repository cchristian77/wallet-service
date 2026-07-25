-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS users
(
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    full_name  VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL,
    password   TEXT         NOT NULL
);

-- Partial unique / non-unique indexes cannot be declared inside CREATE TABLE in Postgres.
CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique_idx ON users (email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS users_deleted_at_idx ON users (deleted_at);


CREATE TABLE IF NOT EXISTS wallets
(
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    user_id    BIGINT    NOT NULL REFERENCES users (id) UNIQUE,
    balance    BIGINT    NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX IF NOT EXISTS wallets_user_id_unique_idx ON wallets (user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS wallets_deleted_at_idx ON wallets (deleted_at);


CREATE TABLE IF NOT EXISTS transactions
(
    id             BIGSERIAL PRIMARY KEY,
    created_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,

    transaction_id VARCHAR(255) UNIQUE,
    status         VARCHAR(16)  NOT NULL CHECK (status IN ('PENDING', 'SUCCESS', 'FAILED'))
);

CREATE INDEX IF NOT EXISTS transactions_status_idx ON transactions (status);


CREATE TABLE IF NOT EXISTS transaction_ledgers
(
    id             BIGSERIAL PRIMARY KEY,
    created_at     TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,

    transaction_id BIGINT      NOT NULL REFERENCES transactions (id),
    wallet_id      BIGINT      NOT NULL REFERENCES wallets (id),
    ledger         VARCHAR(16) NOT NULL CHECK (ledger IN ('DEBIT', 'CREDIT')),
    reference      VARCHAR(32) NOT NULL CHECK (reference IN ('TOP_UP', 'TRANSFER', 'WITHDRAWAL')),
    amount         BIGINT      NOT NULL CHECK (amount > 0)
);

CREATE INDEX IF NOT EXISTS transaction_ledgers_transaction_id_idx ON transaction_ledgers (transaction_id);
CREATE INDEX IF NOT EXISTS transaction_ledgers_wallet_id_idx ON transaction_ledgers (wallet_id);
CREATE INDEX IF NOT EXISTS transaction_ledgers_reference_idx ON transaction_ledgers (reference);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP INDEX IF EXISTS transaction_ledgers_reference_idx;
DROP INDEX IF EXISTS transaction_ledgers_wallet_id_idx;
DROP INDEX IF EXISTS transaction_ledgers_transaction_id_idx;
DROP TABLE IF EXISTS transaction_ledgers;

DROP INDEX IF EXISTS transactions_status_idx;
DROP TABLE IF EXISTS transactions;

DROP INDEX IF EXISTS wallets_deleted_at_idx;
DROP INDEX IF EXISTS wallets_user_id_unique_idx;
DROP TABLE IF EXISTS wallets;

DROP INDEX IF EXISTS users_deleted_at_idx;
DROP INDEX IF EXISTS users_email_unique_idx;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
