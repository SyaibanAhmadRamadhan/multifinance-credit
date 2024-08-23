CREATE TABLE transaction_items
(
    id             SERIAL PRIMARY KEY,
    transaction_id INTEGER      NOT NULL,
    merchant_id    INTEGER      NOT NULL,
    queryParamBindToStruct           VARCHAR(255) NOT NULL,
    image          varchar(255),
    qty            INTEGER,
    unit_price     DOUBLE PRECISION,
    amount         INTEGER,

    CONSTRAINT fk_transaction_items_transaction
        FOREIGN KEY (transaction_id)
            REFERENCES transactions (id)
            ON DELETE CASCADE,

    CONSTRAINT fk_transaction_items_merchant
        FOREIGN KEY (merchant_id)
            REFERENCES merchants (id)
            ON DELETE CASCADE
);

CREATE INDEX idx_transaction_items_transaction_id ON transaction_items (transaction_id);
CREATE INDEX idx_transaction_items_merchant_id ON transaction_items (merchant_id);