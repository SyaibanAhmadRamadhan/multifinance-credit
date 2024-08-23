-- Create the loan table
CREATE TABLE loans
(
    id              INT AUTO_INCREMENT PRIMARY KEY,
    limit_id        INT          NOT NULL,
    consumer_id     INT          NOT NULL,
    bank_account_id INT          NOT NULL,
    admin_fee DOUBLE NOT NULL,
    contract_number BIGINT          NOT NULL,
    date            TIMESTAMP    NOT NULL,
    amount DOUBLE NOT NULL,
    installment_amount DOUBLE NOT NULL,
    interest_rate DOUBLE NOT NULL,
    tenor           INT          NOT NULL,
    status          VARCHAR(255) NOT NULL,

    FOREIGN KEY (limit_id) REFERENCES limits (id) ON DELETE CASCADE,
    FOREIGN KEY (consumer_id) REFERENCES consumers (id) ON DELETE CASCADE,
    FOREIGN KEY (bank_account_id) REFERENCES bank_accounts (id) ON DELETE CASCADE
);

CREATE INDEX index_loans_limit_id ON loans (limit_id);
CREATE INDEX index_loans_consumer_id ON loans (consumer_id);
CREATE INDEX index_loans_bank_account_id ON loans (bank_account_id);
CREATE INDEX index_loans_contract_number ON loans (contract_number);
CREATE INDEX index_loans_date ON loans (date);