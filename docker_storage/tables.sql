CREATE TABLE accounts
(
    account_id INTEGER CONSTRAINT account_id_pk PRIMARY KEY,
    balance NUMERIC CONSTRAINT positive_balance CHECK (balance>=0),
    CONSTRAINT positive_id CHECK (account_id>0)
);

--INSERT INTO accounts (account_id,balance) VALUES (-1,100),(2,200),(14,300),(0,0);

CREATE TABLE transactions_history
(
    account_id INTEGER REFERENCES accounts ON DELETE RESTRICT,
    delta NUMERIC,
    remaining_balance NUMERIC CONSTRAINT positive_balance CHECK (remaining_balance>=0),
    transaction_message TEXT,
    created_at TIMESTAMP 
);