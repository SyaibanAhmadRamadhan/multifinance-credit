CREATE TABLE installments
(
    id                INT AUTO_INCREMENT PRIMARY KEY,
    limit_id          INT          NOT NULL,
    payment_method_id INT NULL,
    contract_number   BIGINT       NOT NULL,
    amount DOUBLE NOT NULL,
    due_date          TIMESTAMP    NOT NULL,
    payment_date      TIMESTAMP,
    status            VARCHAR(255) NOT NULL,

    FOREIGN KEY (limit_id) REFERENCES limits (id) ON DELETE CASCADE,
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods (id) ON DELETE CASCADE
);

CREATE INDEX index_installments_limit_id ON installments (limit_id);
CREATE INDEX index_installments_payment_method_id ON installments (payment_method_id);
CREATE INDEX index_installments_contract_number ON installments (contract_number);
CREATE INDEX index_installments_due_date ON installments (due_date);
CREATE INDEX index_installments_payment_date ON installments (payment_date);