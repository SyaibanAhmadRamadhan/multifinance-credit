-- Create the transactions table
CREATE TABLE transactions
(
    id               INT AUTO_INCREMENT PRIMARY KEY,
    limit_id         INT          NOT NULL,
    consumer_id      INT          NOT NULL,
    contract_number  BIGINT          NOT NULL,
    amount DOUBLE NOT NULL,
    transaction_date TIMESTAMP    NOT NULL,
    status           VARCHAR(255) NOT NULL,
    created_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (limit_id) REFERENCES limits (id) ON DELETE CASCADE,
    FOREIGN KEY (consumer_id) REFERENCES consumers (id) ON DELETE CASCADE
);

CREATE INDEX index_transactions_limit_id ON transactions (limit_id);
CREATE INDEX index_transactions_consumer_id ON transactions (consumer_id);
CREATE INDEX index_transactions_transaction_date ON transactions (transaction_date);
CREATE INDEX index_transactions_contract_number ON transactions (contract_number);
