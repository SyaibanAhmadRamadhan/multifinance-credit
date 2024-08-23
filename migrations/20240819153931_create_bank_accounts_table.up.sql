CREATE TABLE bank_accounts
(
    id                  INT AUTO_INCREMENT PRIMARY KEY,
    consumer_id         INT          NOT NULL,
    queryParamBindToStruct                VARCHAR(100) NOT NULL,
    account_number      VARCHAR(100) NOT NULL,
    account_holder_name VARCHAR(100) NOT NULL,

    FOREIGN KEY (consumer_id) REFERENCES consumers(id) ON DELETE CASCADE
);

CREATE INDEX index_bank_accounts_consumer_id ON bank_accounts(consumer_id);
