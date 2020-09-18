CREATE TABLE accounts
(
    account_id SERIAL CONSTRAINT account_id_pk PRIMARY KEY,
    balance NUMERIC CONSTRAINT positive_balance CHECK (balance>=0)
);

INSERT INTO accounts (balance) VALUES
(100),(200),(300),(0);

CREATE TABLE logs
(
    account_id INTEGER REFERENCES accounts ON DELETE RESTRICT,
    delta NUMERIC,
    log_user_message TEXT,
    log_internal_message TEXT,
    operation_completed BOOLEAN,
    created_at TIMESTAMP 
);